package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"wan_go/pkg/common/constant"
)

func IsEmpty(strs ...string) bool {
	if strs == nil {
		return true
	}
	for _, s := range strs {
		if strings.TrimSpace(s) != "" {
			return false
		}
	}
	return true
}

func IsNotEmpty(strs string) bool {
	return !IsEmpty(strs)
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}
func StringToInt32(i string) int32 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return int32(j)
}
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Uint32ToString(i uint32) string {
	return strconv.FormatInt(int64(i), 10)
}

// judge a string whether in the  string list
func IsContain(target string, List []string) bool {
	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false
}
func IsContainInt32(target int32, List []int32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
func IsContainInt(target int, List []int) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
func InterfaceArrayToStringArray(data []interface{}) (i []string) {
	for _, param := range data {
		i = append(i, param.(string))
	}
	return i
}
func StructToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func StructToJsonBytes(param interface{}) []byte {
	dataType, _ := json.Marshal(param)
	return dataType
}

// The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	err := json.Unmarshal([]byte(s), args)
	return err
}

func GetMsgID(sendID string) string {
	t := int64ToString(GetCurrentTimestampByNano())
	return Md5(t + sendID + int64ToString(rand.Int63n(GetCurrentTimestampByNano())))
}
func GetConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case constant.SingleChatType:
		return "single_" + sourceID
	case constant.GroupChatType:
		return "group_" + sourceID
	case constant.SuperGroupChatType:
		return "super_group_" + sourceID
	case constant.NotificationChatType:
		return "notification_" + sourceID
	}
	return ""
}
func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func RemoveDuplicateElement(idList []string) []string {
	result := make([]string, 0, len(idList))
	temp := map[string]struct{}{}
	for _, item := range idList {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func UUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func Intn(n int) int {
	r := rand.New(rand.NewSource(time.Millisecond.Milliseconds()))
	return r.Intn(n)
}

func CreateCaptcha(num int) string {
	str := "1"
	for i := 0; i < num; i++ {
		str += strconv.Itoa(0)
	}
	str10 := str
	int10, err := strconv.ParseInt(str10, 10, 32)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		j := int32(int10)
		return fmt.Sprintf("%0"+strconv.Itoa(num)+"v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(j))
	}
}
