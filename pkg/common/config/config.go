package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "../../..")
)

//const ConfigName = "landlordConf"

var Config config

type config struct {
	Version string `yaml:"version"`

	User struct {
		CaptchaFormat string `yaml:"captchaFormat"`
	} `yaml:"user"`

	Mail struct {
		Host            string `yaml:"host"`
		Username        string `yaml:"userName"`
		Password        string `yaml:"password"`
		DefaultEncoding string `yaml:"defaultEncoding"`
		Protocol        string `yaml:"protocol"`
		Port            string `yaml:"port"`
		//Properties      struct {
		//	smtp struct{
		//		socketFactory struct{
		//
		//		}`yaml:"socketFactory"`
		//		ssl struct{
		//			enable bool `yaml:"enable"`
		//		}`yaml:"ssl"`
		//	}`yaml:"smtp"`
		//} `yaml:"properties"`
	} `yaml:"mail"`

	Server struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	}

	Api struct {
		Port     []int  `yaml:"port"`
		ListenIP string `yaml:"listenIP"`
	}
	CmsApi struct {
		Port     []int  `yaml:"port"`
		ListenIP string `yaml:"listenIP"`
	}
	Sdk struct {
		WsPort  []int    `yaml:"wsPort"`
		DataDir []string `yaml:"dataDir"`
	}
	Mysql struct {
		Address       []string `yaml:"address"`
		UserName      string   `yaml:"userName"`
		Password      string   `yaml:"password"`
		DatabaseName  string   `yaml:"databaseName"`
		TableName     string   `yaml:"tableName"`
		MsgTableNum   int      `yaml:"msgTableNum"`
		MaxOpenConns  int      `yaml:"maxOpenConns"`
		MaxIdleConns  int      `yaml:"maxIdleConns"`
		MaxLifeTime   int      `yaml:"maxLifeTime"`
		LogLevel      int      `yaml:"logLevel"`
		SlowThreshold int      `yaml:"slowThreshold"`
	}
	Mongo struct {
		DBUri                string   `yaml:"dbUri"`
		DBAddress            []string `yaml:"dbAddress"`
		DBDirect             bool     `yaml:"dbDirect"`
		DBTimeout            int      `yaml:"dbTimeout"`
		DBDatabase           string   `yaml:"dbDatabase"`
		DBSource             string   `yaml:"dbSource"`
		DBUserName           string   `yaml:"dbUserName"`
		DBPassword           string   `yaml:"dbPassword"`
		DBMaxPoolSize        int      `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords  int      `yaml:"dbRetainChatRecords"`
		ChatRecordsClearTime string   `yaml:"chatRecordsClearTime"`
	}
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBUserName    string   `yaml:"dbUserName"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}
	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
		UserName   string   `yaml:"userName"`
		Password   string   `yaml:"password"`
		Secret     string   `yaml:"secret"`
	}

	Kafka struct {
		SASLUserName string `yaml:"SASLUserName"`
		SASLPassword string `yaml:"SASLPassword"`
		Ws2mschat    struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		//Ws2mschatOffline struct {
		//	Addr  []string `yaml:"addr"`
		//	Topic string   `yaml:"topic"`
		//}
		MsgToMongo struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		MsgToModify struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		ConsumerGroupID struct {
			MsgToRedis  string `yaml:"msgToTransfer"`
			MsgToMongo  string `yaml:"msgToMongo"`
			MsgToMySql  string `yaml:"msgToMySql"`
			MsgToPush   string `yaml:"msgToPush"`
			MsgToModify string `yaml:"msgToModify"`
		}
	}
	Landlords struct {
		MaxSecondsForEveryRound int64 `yaml:"max_seconds_for_every_round"`
	}

	Session struct {
		UserSessionKey string `yaml:"user_session_key"`
		Secret         string `yaml:"secret"`
		Name           string `yaml:"name"`
	}

	TLS struct {
		Addr string `yaml:"addr"`
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	}

	TokenPolicy struct {
		JwtSecret string `yaml:"jwt_secret"`
		JwtExpire int64  `yaml:"jwt_expire"`
	}

	Websocket struct {
		Port             []string `yaml:"port"`
		MaxConnNum       int      `yaml:"max_conn_num"`
		MaxMsgLen        int      `yaml:"max_msg_len"`
		HandshakeTimeOut int      `yaml:"handshake_time_out"`
		OnlineTimeOut    int      `yaml:"online_time_out"`
	}

	Log struct {
		StorageLocation       string   `yaml:"storageLocation"`
		RotationTime          int      `yaml:"rotationTime"`
		RemainRotationCount   uint     `yaml:"remainRotationCount"`
		RemainLogLevel        uint     `yaml:"remainLogLevel"`
		ElasticSearchSwitch   bool     `yaml:"elasticSearchSwitch"`
		ElasticSearchAddr     []string `yaml:"elasticSearchAddr"`
		ElasticSearchUser     string   `yaml:"elasticSearchUser"`
		ElasticSearchPassword string   `yaml:"elasticSearchPassword"`
	}

	Blog struct {
		Name         string `yaml:"name"`
		Port         []int  `yaml:"port"`
		ListenIP     string `yaml:"listenIP"`
		AliSMSVerify struct {
			AccessKeyID                  string `yaml:"accessKeyId"`
			AccessKeySecret              string `yaml:"accessKeySecret"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
			Enable                       bool   `yaml:"enable"`
		}
		TencentSMS struct {
			AppID                        string `yaml:"appID"`
			Region                       string `yaml:"region"`
			SecretID                     string `yaml:"secretID"`
			SecretKey                    string `yaml:"secretKey"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
			Enable                       bool   `yaml:"enable"`
		}
		SuperCode    string `yaml:"superCode"`
		CodeTTL      int    `yaml:"codeTTL"`
		UseSuperCode bool   `yaml:"useSuperCode"`
		Mail         struct {
			Title                   string `yaml:"title"`
			SenderMail              string `yaml:"senderMail"`
			SenderAuthorizationCode string `yaml:"senderAuthorizationCode"`
			SmtpAddr                string `yaml:"smtpAddr"`
			SmtpPort                int    `yaml:"smtpPort"`
		}
	}

	Prometheus struct {
		Enable                        bool  `yaml:"enable"`
		UserPrometheusPort            []int `yaml:"userPrometheusPort"`
		FriendPrometheusPort          []int `yaml:"friendPrometheusPort"`
		MessagePrometheusPort         []int `yaml:"messagePrometheusPort"`
		MessageGatewayPrometheusPort  []int `yaml:"messageGatewayPrometheusPort"`
		GroupPrometheusPort           []int `yaml:"groupPrometheusPort"`
		AuthPrometheusPort            []int `yaml:"authPrometheusPort"`
		PushPrometheusPort            []int `yaml:"pushPrometheusPort"`
		AdminCmsPrometheusPort        []int `yaml:"adminCmsPrometheusPort"`
		OfficePrometheusPort          []int `yaml:"officePrometheusPort"`
		OrganizationPrometheusPort    []int `yaml:"organizationPrometheusPort"`
		ConversationPrometheusPort    []int `yaml:"conversationPrometheusPort"`
		CachePrometheusPort           []int `yaml:"cachePrometheusPort"`
		RealTimeCommPrometheusPort    []int `yaml:"realTimeCommPrometheusPort"`
		MessageTransferPrometheusPort []int `yaml:"messageTransferPrometheusPort"`
	} `yaml:"prometheus"`
}

func init() {
	unmarshalConfig(&Config, "config.yaml")
}

func unmarshalConfig(config interface{}, configName string) {

	env := "CONFIG_NAME"

	cfgName := os.Getenv(env)

	if len(cfgName) != 0 {
		bytes, err := os.ReadFile(filepath.Join(cfgName, "config", configName))
		if err != nil {
			bytes, err = os.ReadFile(filepath.Join(Root, "config", configName))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", configName))
			}
		} else {
			Root = cfgName
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	} else {
		dir, _ := os.Getwd()
		bytes, err := os.ReadFile(fmt.Sprintf("%s/config/%s", dir, configName))
		if err != nil {
			panic(err.Error() + configName)
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	}
}
