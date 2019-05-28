package web

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"github.com/zhin/go-codex/cerror"
)

var Default *gin.Engine

var showErrValue = false

var codeField = "code"
var msgField = "msg"
var dataField = "data"
var errIDField = "errid"

func init() {

	Default = gin.New()

}

type JSONResult struct {
	Code  int
	Msg   string
	ErrID string
	Data  interface{}

	codeField  string
	msgField   string
	dataField  string
	errIDField string
}

func (u *JSONResult) MarshalJSON() ([]byte, error) {

	val := map[string]interface{}{}

	if u.codeField != "" {
		val[u.codeField] = u.Code
	} else {
		val[codeField] = u.Code
	}

	if u.msgField != "" {
		val[u.msgField] = u.Msg
	} else {
		val[msgField] = u.Msg
	}

	if u.Data != nil {
		if u.dataField != "" {
			val[u.dataField] = u.Data
		} else {
			val[dataField] = u.Data
		}
	}

	if u.ErrID != "" {
		if u.ErrID != "" {
			val[u.errIDField] = u.ErrID
		} else {
			val[errIDField] = u.ErrID
		}
	}
	return json.Marshal(&val)
}

type DBErrorHandle func(errID string, err error)

var errorHandles = []DBErrorHandle{}

func SetDBErrorHook(handle DBErrorHandle) {
	errorHandles = append(errorHandles, handle)
}

func triggerErrorHandles(errID string, err error) {
	if errorHandles != nil {
		for _, handle := range errorHandles {
			if handle != nil {
				handle(errID, err)
			}
		}
	}
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
	var err error
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
			} else if er, ok := msg[0].(error); ok {
				r.Code = 1
				err = er
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
	if err != nil {
		guid := uuid.NewV4()
		if showErrValue {
			r.Msg = err.Error()
		} else {
			r.Msg = "ERR OCCURRED"
		}
		r.ErrID = guid.String()
		triggerErrorHandles(guid.String(), err)
	}
	return r
}

func (r *JSONResult) SetData(data interface{}) *JSONResult {
	r.Data = data
	return r
}

func (r *JSONResult) SetFiled(codeField string, msgField string, dataField string, errIDField string) *JSONResult {

	if codeField != "" {
		r.codeField = codeField
	}
	if msgField != "" {
		r.msgField = msgField
	}
	if dataField != "" {
		r.dataField = dataField
	}
	if errIDField != "" {
		r.errIDField = errIDField
	}
	return r
}

func ShowErrorDetail(value bool) {
	showErrValue = value
}

func SetJSONResultField(_codeField string, _msgField string, _dataField string, _errIDField string) {
	if _codeField != "" {
		codeField = _codeField
	}
	if _msgField != "" {
		msgField = _msgField
	}
	if _dataField != "" {
		dataField = _dataField
	}
	if _errIDField != "" {
		errIDField = _errIDField
	}
}
