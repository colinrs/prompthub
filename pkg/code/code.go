package code

import "net/http"

var (
	OKErr                      = &Err{HTTPCode: http.StatusOK, Code: 0, Msg: "success"}
	ErrParam                   = &Err{HTTPCode: http.StatusOK, Code: 10003, Msg: "参数有误"}
	UnknownErr                 = &Err{HTTPCode: http.StatusOK, Code: 10004, Msg: "unknown error"}
	ErrValidation              = &Err{HTTPCode: http.StatusOK, Code: 20001, Msg: "Validation failed."}
	ErrDatabase                = &Err{HTTPCode: http.StatusOK, Code: 20002, Msg: "Database error."}
	BizTagAlreadyExist         = &Err{HTTPCode: http.StatusOK, Code: 20003, Msg: "业务已存在"}
	BizTagNotExist             = &Err{HTTPCode: http.StatusOK, Code: 20004, Msg: "业务不存在"}
	EtcdKeyNotExist            = &Err{HTTPCode: http.StatusOK, Code: 20005, Msg: "Etcd key 不存在"}
	HTTPClientErr              = &Err{HTTPCode: http.StatusOK, Code: 20006, Msg: "http client error"}
	ErrUserExist               = &Err{HTTPCode: http.StatusOK, Code: 20007, Msg: "用户已经存在"}
	ErrUserNotExist            = &Err{HTTPCode: http.StatusOK, Code: 20008, Msg: "用户不存在"}
	ErrPasswordIncorrect       = &Err{HTTPCode: http.StatusOK, Code: 20009, Msg: "用户密码错误"}
	ErrUserNoLogin             = &Err{HTTPCode: http.StatusOK, Code: 20010, Msg: "用户未登录"}
	PromptNotFound             = &Err{HTTPCode: http.StatusOK, Code: 20011, Msg: "prompt not found"}
	ErrCategoryAlreadyExists   = &Err{HTTPCode: http.StatusOK, Code: 20012, Msg: "分类已存在"}
	ErrVerificationCodeInvalid = &Err{HTTPCode: http.StatusOK, Code: 20013, Msg: "验证码无效"}
	ErrEmailVerification       = &Err{HTTPCode: http.StatusOK, Code: 20014, Msg: "用户邮件待验证"}
	ErrVerificationLimitExceed = &Err{HTTPCode: http.StatusOK, Code: 20015, Msg: "验证次数超过限制"}
	ErrPasswordSameInvalid     = &Err{HTTPCode: http.StatusOK, Code: 20016, Msg: "密码和之前密码相同"}
	ErrSensitiveWord           = &Err{HTTPCode: http.StatusOK, Code: 20017, Msg: "包含敏感词"}
)
