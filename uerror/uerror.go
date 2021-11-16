package uerror

import (
	"fmt"
	"runtime"
)

type Error struct {
	Code int
	Msg  string
	Desc string

	// 调用栈
	st []uintptr
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

func New(code int, msg string) error {
	err := &Error{
		Code: code,
		Msg:  msg,
	}
	err.st = callers()
	return err
}

func callers() []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(4, pcs[:])
	st := pcs[0:n]
	return st
}

func Code(e error) int {
	if e == nil {
		return 0
	}
	err, ok := e.(*Error)
	if !ok {
		return 0
	}
	if err == (*Error)(nil) {
		return 0
	}
	return int(err.Code)
}

func Msg(e error) string {
	if e == nil {
		return "success"
	}
	err, ok := e.(*Error)
	if !ok {
		return e.Error()
	}
	if err == (*Error)(nil) {
		return "success"
	}
	return err.Msg
}

func Err(e error) error {
	return New(500, e.Error())
}
