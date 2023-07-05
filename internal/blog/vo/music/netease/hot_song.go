package netease

import "wan_go/internal/blog/vo/music"

type HotSong struct {
	Code int         `json:"code"`
	Data *music.Song `json:"data"`
}
