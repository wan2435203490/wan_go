package ws

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/token"
	"wan_go/pkg/utils"
)

type UserConn struct {
	*websocket.Conn
	mu         sync.Mutex
	UserId     int32
	IsCompress bool
	IsOnline   bool
}

var WS wsServer

type wsServer struct {
	Port       string //443
	MaxConnNum int
	UpGrader   *websocket.Upgrader
	//UserConnMap map[string]*UserConn
	UserConnMap sync.Map
}

func Start() {
	WS.onInit()
	WS.run()
}

func (ws *wsServer) onInit() {
	conf := config.Config.Landlords.Websocket
	ws.Port = conf.Port[0]
	ws.MaxConnNum = conf.MaxConnNum
	ws.UpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(conf.HandshakeTimeOut) * time.Second,
		ReadBufferSize:   conf.MaxMsgLen,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
}

func (ws *wsServer) run() {
	http.HandleFunc("/ws", ws.Handler)
	err := http.ListenAndServe(":"+ws.Port, nil)
	if err != nil {
		panic("ws listening err:" + err.Error())
	}
}

func (ws *wsServer) Handler(w http.ResponseWriter, r *http.Request) {
	//query := r.URL.Query()

	if isPass, userId := ws.headerCheck(w, r); isPass {
		conn, err := ws.UpGrader.Upgrade(w, r, nil)
		if err == nil {
			uc := &UserConn{UserId: userId, Conn: conn, IsOnline: true}
			ws.UserConnMap.Store(userId, uc)
			go ws.readMsg(uc)

			//ticker := time.NewTicker(time.Duration(config.Config.Websocket.OnlineTimeOut) * time.Second)
			//go uc.Ping(ticker)
		}
	}
}

func (uc *UserConn) Ping(ticker *time.Ticker) {
	for {
		<-ticker.C
		uc.SetWriteDeadline(time.Now().Add(time.Duration(config.Config.Websocket.OnlineTimeOut) * time.Second))

		uc.mu.Lock()
		if err := uc.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("ping error: %s\n", err.Error())
			uc.IsOnline = false
			//todo
			Push("ExitRoom", &blog.User{ID: uc.UserId})
			//svc.ExitRoom(&db.User{ID: uc.UserId})
			uc.mu.Unlock()
			break
		} else {
			uc.mu.Unlock()
		}

	}
}

var registers sync.Map

func Push(method string, args any) {

	//todo
	registers.Range(func(key, value any) bool {
		//if key ...
		//binder := value.(Binder)
		//binder.Func()
		return true
	})
}

func (ws *wsServer) headerCheck(w http.ResponseWriter, r *http.Request) (isPass bool, userId int32) {
	status := http.StatusUnauthorized
	var err error
	query := r.URL.Query()
	if len(query["token"]) != 0 {
		if userId, err = token.GetUserIdFromToken(query["token"][0]); err != nil {
			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("ws_err_msg", fmt.Sprintf("decode token failure:%s", err.Error()))
			http.Error(w, "error token", status)
			return false, -1
		} else {
			return true, userId
		}
	} else {
		status = int(constant.ErrArgs.ErrCode)
		w.Header().Set("Sec-Websocket-Version", "13")
		errMsg := "args err, need token"
		w.Header().Set("ws_err_msg", errMsg)
		http.Error(w, errMsg, status)
		return false, -1
	}
}

func (ws *wsServer) readMsg(conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()

		if err != nil {
			//ws.UserConnMap.Delete(conn.UserId)
			return
		}

		switch messageType {
		case 0:

		}

		if conn.IsCompress {
			buff := bytes.NewBuffer(msg)
			reader, err := gzip.NewReader(buff)
			if err != nil {

				continue
			}
			msg, err = io.ReadAll(reader)
			if err != nil {

				continue
			}
			err = reader.Close()
			if err != nil {

			}
		}

		println(msg)
		//ws.msgParse(conn, msg)
	}
}

func (ws *wsServer) writeMsg(conn *UserConn, msgType int, msg []byte) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.IsCompress {
		var buffer bytes.Buffer
		gz := gzip.NewWriter(&buffer)
		if _, err := gz.Write(msg); err != nil {
			return utils.WrapMsg(err, "")
		}
		if err := gz.Close(); err != nil {
			return utils.WrapMsg(err, "")
		}
		msg = buffer.Bytes()
	}
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	return conn.WriteMessage(msgType, msg)
}

func (ws *wsServer) IsOnline(userId int32) bool {
	if value, ok := ws.UserConnMap.Load(userId); !ok {
		return false
	} else {
		conn := value.(*UserConn)
		conn.mu.Lock()
		defer conn.mu.Unlock()
		return true
		//if conn.IsOnline {
		//	return true
		//} else {
		//	conn = nil
		//	ws.UserConnMap.Delete(userId)
		//	return false
		//}
	}
	return true
}

func (ws *wsServer) Send2Users(userIds []int32, content string) error {

	//res := strings.Builder{}

	var errs error

	for _, id := range userIds {
		//res.WriteString(ws.Send2User(id, content))
		if err := ws.Send2User(id, content); err != nil {
			if errs == nil {
				errs = utils.Wrap(err)
			} else {
				errs = utils.WrapMsg(err, errs.Error())
			}
		}
	}

	return errs
}

func (ws *wsServer) Send2User(userId int32, content string) error {
	if value, ok := ws.UserConnMap.Load(userId); ok {
		conn := value.(*UserConn)
		if !conn.IsOnline {
			return errors.New(fmt.Sprintf("玩家不在线userId:%d\n", userId))
		} else {
			log.Printf("websocket:%d %s\n", userId, content)
			err := conn.WriteMessage(websocket.TextMessage, []byte(content))
			if err != nil {
				return errors.New(fmt.Sprintf("消息推送异常userId:%d,err:%s\n", userId, err.Error()))
			} else {
				return nil
			}
		}
	}

	return errors.New("用户连接不存在，请重试")
}

func (ws *wsServer) Send2AllUser(content string) error {
	//var sb strings.Builder

	var errs error

	ws.UserConnMap.Range(func(key, value any) bool {
		conn := value.(*UserConn)
		if !conn.IsOnline {
			//sb.WriteString("玩家不在线")
			if errs == nil {
				errs = utils.Wrap(errors.New(fmt.Sprintf("玩家%d不在线", conn.UserId)))
			} else {
				errs = utils.WrapMsg(errors.New(fmt.Sprintf("玩家%d不在线", conn.UserId)), errs.Error())
			}
		} else {
			err := conn.WriteMessage(websocket.TextMessage, []byte(content))
			if err != nil {
				msg := fmt.Sprintf("消息推送异常userId:%s,err:%s\n", key, err.Error())
				//sb.WriteString(msg)
				if errs == nil {
					errs = utils.Wrap(errors.New(msg))
				} else {
					errs = utils.WrapMsg(errors.New(msg), errs.Error())
				}
			}
		}
		return true
	})

	return errs
}
