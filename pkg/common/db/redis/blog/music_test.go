package blog

import (
	"testing"
)

func TestCacheSongInfo(t *testing.T) {
	CacheSongInfo("neteaseRand", "url", "qwq")
}

func TestGetCacheInfo(t *testing.T) {
	GetSongInfo("a", 3)
}
