package helper

import (
	"encoding/json"
	"fmt"
	"time"
)

func ArrayContains(arr interface{}, item interface{}) bool {
	newArr, ok := arr.([]interface{})
	if !ok {
		return false
	}

	for _, v := range newArr {
		if v == item {
			return true
		}
	}
	return false
}

func PrettyJson(data interface{}) string {
	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("<failed to parse json: %v>", err.Error())
	}
	return string(res)
}

func TimeNowUTC() time.Time {
	return time.Now().UTC()
}
