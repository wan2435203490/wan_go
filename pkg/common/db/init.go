package db

import (
	"fmt"
	"github.com/dtm-labs/rockscache"
	go_redis "github.com/go-redis/redis/v8"
	"gopkg.in/mgo.v2"
	"gorm.io/gorm"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/api"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
	//go_redis "github.com/go-redis/redis/v8"
)

var DB DataBases

type DataBases struct {
	MysqlDB     *gorm.DB
	mgoSession  *mgo.Session
	mongoClient *mongo.Client
	RDB         go_redis.UniversalClient
	Rc          *rockscache.Client
	WeakRc      *rockscache.Client
}

func init() {

	fmt.Println("initiating mysql redis mongo ")

	blog.InitMysql(&DB)

	initMongo()

	initRedis()

	fmt.Println("mysql redis mongo initiated")
}

func Mysql() *gorm.DB {
	return DB.MysqlDB
}

// Page todo 怎么改成传BaseRequestVO[T]？
func Page(pagination *api.Pagination) (db *gorm.DB) {

	db = Mysql().
		Order(pagination.Order()).
		Limit(pagination.Size).
		Offset(pagination.Current)

	return
}