package helper

import (
	"book_service/domain/dto"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

func GetCurrentUserFromGinCtx(ctx *gin.Context) (*dto.CurrentUser, error) {
	rawCurrentUser, ok := ctx.Get("currentUser")
	if !ok {
		return nil, fmt.Errorf("invalid currentUser: %v", rawCurrentUser)
	}
	currentUser, ok := rawCurrentUser.(dto.CurrentUser)
	if !ok {
		return nil, fmt.Errorf("invalid currentUser: %v", rawCurrentUser)
	}
	return &currentUser, nil
}
