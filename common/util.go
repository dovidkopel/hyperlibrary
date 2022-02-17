package common

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

type History struct {
	Time    time.Time              `json:"time"`
	Data    map[string]interface{} `json:"data"`
	Deleted bool                   `json:"deleted"`
}

func GetApproxTime(ts *timestamp.Timestamp) time.Time {
	return time.Unix(ts.Seconds, int64(ts.Nanos)).Round(time.Minute).UTC()
}

//func GetKeys[T, S](m map[T]S) []T {
//	keys := make([]T, len(m))
//	for k, _ := range m {
//		keys = append(keys, k)
//	}
//	return keys
//}

func Contains(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//func GetValues(m map[interface{}]interface{}) []interface{} {
//	values := make([]interface{}, len(m))
//	for _, v := range m {
//		values = append(values, v)
//	}
//	return values
//}
