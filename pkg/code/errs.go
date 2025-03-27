package code

import "net/http"

type Err struct {
	HTTPCode int      `json:"http_code"`
	Code     int      `json:"code"`
	Msg      string   `json:"msg"`
	Errors   []*Error `json:"errors,omitempty"`
}

func (e *Err) Error() string {
	return e.Msg
}

func (e *Err) GetCode() int {
	return e.Code
}

func (e *Err) GetHTTPCode() int {
	return e.HTTPCode
}

func (e *Err) GetMsg() string {
	return e.Msg
}

func (e *Err) WithErrors(items []*Error) {
	e.Errors = append(e.Errors, items...)
}

func (e *Err) GetErrors() []*Error {
	return e.Errors
}

func (e *Err) Copy() *Err {
	tempE := *e
	return &tempE
}

type Error struct {
	Attr   string `json:"attr,omitempty"`
	Code   int    `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (e *Error) Error() string {
	return e.Detail
}

func (e *Error) GetCode() int {
	return e.Code
}

type Option func(e *Err)

func WithHTTPCode(httpCode int) Option {
	return func(e *Err) {
		e.HTTPCode = httpCode
	}
}

func WithCode(code int) Option {
	return func(e *Err) {
		e.Code = code
	}
}

func WithMsg(msg string) Option {
	return func(e *Err) {
		e.Msg = msg
	}
}

func WithErrors(errs ...*Error) Option {
	return func(e *Err) {
		e.Errors = append(e.Errors, errs...)
	}
}

func NewErr(options ...Option) *Err {
	e := &Err{
		HTTPCode: http.StatusOK,
		Code:     0,
		Msg:      "",
		Errors:   nil,
	}
	for _, option := range options {
		option(e)
	}
	return e
}
