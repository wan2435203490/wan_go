package apis

import (
	json2 "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"sync"
	"wan_go/internal/blog/vo/music"
	"wan_go/internal/blog/vo/music/netease"
	"wan_go/internal/blog/vo/music/qq"
	"wan_go/pkg/common/api"
	blogRedis "wan_go/pkg/common/db/redis/blog"
	"wan_go/pkg/utils"
)

type Hash func(data []byte) uint32

func New(hash Hash) *MusicApi {
	return &MusicApi{}
	//return &MusicApi{hash: hash}
}

type MusicApi struct {
	api.Api
	//hash Hash
}

const maxCacheCount = 100
const keyNeteaseRand = "neteaseRand"
const keyQQRand = "qqRand"

const neteaseHot = "https://api.uomg.com/api/rand.music?sort=热歌榜&format=json"

// todo 集成 https://rain120.github.io/qq-music-api

const qqRandom = "https://apis.jxcxin.cn/api/qqrandommusic?type=json"
const qqDetail = "https://apis.jxcxin.cn/api/qqmusic?url=https://y.qq.com/n/ryqq/songDetail/"
const qqSplit = ".m4a?"

//todo 追书api
//const zhuishu = "https://www.cnblogs.com/Stars-are-shining/p/13345856.html"

//todo 快看漫画api

func Get(url string, result any) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json2.Unmarshal(body, result)
	fmt.Println(url, utils.GetCurrentTimestampByMill())
	return err
}

func GetQQRandom() (*music.Song, error) {

	randomResult := qq.RandomResult{}
	if err := Get(qqRandom, &randomResult); err != nil || randomResult.Data == nil {
		return nil, err
	}

	randomSong := randomResult.Data
	url := randomSong.Url
	index := strings.Index(url, qqSplit)
	mid := url[index-14 : index]

	detail := qq.MusicDetail{}
	if err := Get(qqDetail+mid, &detail); err != nil {
		return nil, err
	}

	song := music.Song{
		Name:        randomSong.SongName,
		Url:         randomSong.Url,
		PicUrl:      detail.Pic,
		ArtistsName: randomSong.Author,
	}

	return &song, nil
}

func GetNeteaseHot() (*music.Song, error) {

	song := netease.HotSong{}
	if err := Get(neteaseHot, &song); err != nil {
		return nil, err
	}

	return song.Data, nil
}

func (a MusicApi) RandomMusic(c *gin.Context) {
	a.MakeContext(c)

	count := 0
	if a.IntFailed(&count, "count") {
		return
	}

	//校验redis里存了多少个 当有 maxCacheCount 个时则不去api获取，直接去redis拿
	countNetease := blogRedis.CountSongInfo(keyNeteaseRand)
	countQQ := blogRedis.CountSongInfo(keyQQRand)

	neteaseSongs := music.RandomSong{TypeName: "网易云热歌"}
	qqSongs := music.RandomSong{TypeName: "QQ音乐随机"}

	//不去校验那么仔细 大约存 maxCacheCount 个数据就好
	if countNetease > maxCacheCount || countQQ > maxCacheCount {
		netEaseResult := blogRedis.GetSongInfo(keyNeteaseRand, count)
		qqResult := blogRedis.GetSongInfo(keyQQRand, count)

		neteaseSongs.Songs = netEaseResult
		qqSongs.Songs = qqResult
	} else {
		//是否能优化？
		wg := sync.WaitGroup{}
		var neteaseMu sync.Mutex
		//var qqMu sync.Mutex
		netEaseResult := make([]*music.Song, 0)
		qqResult := make([]*music.Song, 0)
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func() {
				song, _ := GetNeteaseHot()
				json, _ := json2.Marshal(song)
				_ = blogRedis.CacheSongInfo(keyNeteaseRand, song.Url, json)
				neteaseMu.Lock()
				netEaseResult = append(netEaseResult, song)
				neteaseMu.Unlock()
				wg.Done()
			}()
			//go func() {
			//	song, _ := GetQQRandom()
			//	json, _ := json2.Marshal(song)
			//	if song == nil {
			//		return
			//	}
			//	_ = blogRedis.CacheSongInfo(keyQQRand, song.Url, json)
			//	qqMu.Lock()
			//	qqResult = append(qqResult, song)
			//	qqMu.Unlock()
			//	wg.Done()
			//}()
		}
		wg.Wait()

		neteaseSongs.Songs = &netEaseResult
		qqSongs.Songs = &qqResult
	}

	result := []*music.RandomSong{&neteaseSongs, &qqSongs}

	a.OK(result)
}

func (a MusicApi) SongUrl(c *gin.Context) {
	a.MakeContext(c)

	id := a.Param("id")
	//tr := http.Request{}
	if id == "" {
		id = "33894312"
	}
	response, err := http.Get("http://localhost:3000/song/url?id=" + id)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	songUrl := music.SongUrl{}
	err = json2.Unmarshal(body, &songUrl)
	if err != nil {
		return
	}
}
