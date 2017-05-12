package libs

import (
	"cn/com/hengwei/commons"
	"fmt"
	"net/http"
)

type ErrNotFound struct {
	id interface{}
}

func (err *ErrNotFound) Code() int {
	return http.StatusNotFound
}

func (err *ErrNotFound) Error() string {
	if nil == err.id {
		return "not found"
	}
	return "record with id is '" + fmt.Sprint(err.id) + "' isn't found"
}

func NotFound(id interface{}) commons.RuntimeError {
	return &ErrNotFound{id: id}
}
