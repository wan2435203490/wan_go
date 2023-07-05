package netease

import "wan_go/internal/vo/blog/music"

type HotSong struct {
	Code int         `json:"code"`
	Data *music.Song `json:"data"`
}
