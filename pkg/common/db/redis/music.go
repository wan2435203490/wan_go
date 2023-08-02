package redis

import (
	json2 "encoding/json"
	go_redis "github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"time"
	"wan_go/internal/blog/vo/music"
	"wan_go/pkg/common/db"
)

const KeySongUrl = "SongUrl"

func rdb() go_redis.UniversalClient {
	return db.DB.RDB
}

// 这里用 hash 是想将url区分出来 以后可能有别的处理
// CacheSongInfo
// typ 区分网易云 QQ音乐
func CacheSongInfo(typ, url string, info any) error {
	rdb().Expire(context.Background(), typ, time.Hour*2)
	return rdb().HSet(context.Background(), typ, url, info).Err()
}

func CountSongInfo(typ string) int {
	val := rdb().HLen(context.Background(), typ).Val()
	return int(val)
}

func GetSongInfo(typ string, count int) *[]*music.Song {
	val := rdb().HRandField(context.Background(), typ, count, true).Val()
	var ret []*music.Song

	for i, v := range val {
		if i%2 == 0 {
			//暂不处理
			continue
		}
		var song music.Song
		_ = json2.Unmarshal([]byte(v), &song)
		ret = append(ret, &song)
	}

	return &ret
}
