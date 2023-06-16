package log

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"wan_go/pkg/utils"
)

type fileHook struct{}

func newFileHook() *fileHook {
	return &fileHook{}
}

func (f *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//func (f *fileHook) Fire(entry *logrus.Entry) error {
//	entry.Data["FilePath"] = findCaller(6)
//	utils.GetSelfFuncName()
//	return nil
//}

func (f *fileHook) Fire(entry *logrus.Entry) error {
	var s string
	_, file, line, _ := runtime.Caller(8)
	i := strings.SplitAfter(file, "/")
	if len(i) > 3 {
		s = i[len(i)-3] + i[len(i)-2] + i[len(i)-1] + ":" + utils.IntToString(line)
	}
	entry.Data["FilePath"] = s
	return nil
}
