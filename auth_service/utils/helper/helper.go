package helper

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

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
