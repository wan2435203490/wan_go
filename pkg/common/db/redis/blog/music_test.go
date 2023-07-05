package blog

import (
	"github.com/go-redis/redis/v8"
	"reflect"
	"testing"
)

func TestCacheSongInfo(t *testing.T) {
	type args struct {
		url  string
		info any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CacheSongInfo(tt.args.url, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("CacheSongInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCountSongInfo(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountSongInfo(); got != tt.want {
				t.Errorf("CountSongInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rdb(t *testing.T) {
	tests := []struct {
		name string
		want redis.UniversalClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rdb(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rdb() = %v, want %v", got, tt.want)
			}
		})
	}
}
