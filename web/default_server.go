package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zhin/go-codex/cerror"
)

var Default *gin.Engine

func init() {

	Default = gin.New()

}

type JSONResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewJSONResult() *JSONResult {
	return &JSONResult{}
}

func (r *JSONResult) ParameterError(msg ...string) *JSONResult {
	r.Code = 22
	if msg == nil {
		msg = []string{"参数错误"}
	}
	r.Msg = msg[0]
	return r
}

func (r *JSONResult) Success(msg ...string) *JSONResult {
	r.Code = 0
	if msg == nil {
		msg = []string{"OK"}
	}
	r.Msg = msg[0]
	return r
}

func (r *JSONResult) Error(msg ...interface{}) *JSONResult {

	if msg != nil {
		if len(msg) == 1 {
			if val, ok := msg[0].(int); ok {
				r.Code = val
			} else if val, ok := msg[0].(cerror.CodeError); ok {
				r.Code = val.Code
				r.Msg = val.Error()
			} else if val, ok := msg[0].(string); ok {
				r.Code = 1
				r.Msg = val
			}
		} else if len(msg) == 2 {
			r.Code = msg[0].(int)
			r.Msg = msg[1].(string)
		}
	} else {

	}
	if r.Msg == "" {
		r.Msg = "ERR"
	}

	return r
}

func (r *JSONResult) SetData(data interface{}) *JSONResult {
	r.Data = data
	return r
}
