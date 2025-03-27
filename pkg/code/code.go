package code

import "net/http"

var (
	OKErr              = &Err{HTTPCode: http.StatusOK, Code: 0, Msg: "success"}
	ErrParam           = &Err{HTTPCode: http.StatusOK, Code: 10003, Msg: "参数有误"}
	UnknownErr         = &Err{HTTPCode: http.StatusOK, Code: 10004, Msg: "unknown error"}
	ErrValidation      = &Err{HTTPCode: http.StatusOK, Code: 20001, Msg: "Validation failed."}
	ErrDatabase        = &Err{HTTPCode: http.StatusOK, Code: 20002, Msg: "Database error."}
	BizTagAlreadyExist = &Err{HTTPCode: http.StatusOK, Code: 20003, Msg: "业务已存在"}
	BizTagNotExist     = &Err{HTTPCode: http.StatusOK, Code: 20004, Msg: "业务不存在"}
	EtcdKeyNotExist    = &Err{HTTPCode: http.StatusOK, Code: 20005, Msg: "Etcd key 不存在"}
	HTTPClientErr      = &Err{HTTPCode: http.StatusOK, Code: 20006, Msg: "http client error"}
)
