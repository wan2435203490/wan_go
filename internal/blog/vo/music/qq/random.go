package qq

type RandomResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data *RandomSong `json:"data"`
}

type RandomSong struct {
	SongName  string `json:"songname"`
	AlbumName string `json:"albumname"`
	Author    string `json:"author"`
	Url       string `json:"url"`
}
