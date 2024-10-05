package errors

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
)

const MaxDepth = 24

type ThrowError struct {
	err        error
	stacktrace []string
}

var Is func(err, target error) bool = stderrors.Is
var As func(err error, target any) bool = stderrors.As
var Unwrap func(err error) error = stderrors.Unwrap
var Join func(errs ...error) error = stderrors.Join

func New(text string) error {
	return Wrap(stderrors.New(text))
}

func As2[T error](err error) (T, bool) {
	var t T
	ok := As(err, &t)
	return t, ok
}

func (m ThrowError) MarshalJSON() ([]byte, error) {
	v := struct {
		Error      string   `json:"error"`
		Stacktrace []string `json:"stack"`
	}{
		Error:      m.err.Error(),
		Stacktrace: m.stacktrace,
	}
	return json.Marshal(v)
}

func (m ThrowError) Error() string {
	return m.err.Error()
}

func (m ThrowError) Unwrap() error {
	return m.err
}

func Errorf(format string, args ...any) error {
	return Wrap(fmt.Errorf(format, args...))
}

func Attr(err error) slog.Attr {
	return slog.Any("errors", Wrap(err))
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var terr ThrowError

	// do not re-wrap
	if stderrors.As(err, &terr) {
		terr.err = err
		return terr
	}

	return ThrowError{err: err, stacktrace: getStackTrace()}
}

func getStackTrace() []string {
	stackBuffer := make([]uintptr, MaxDepth)
	length := runtime.Callers(3, stackBuffer[:])
	stack := stackBuffer[:length]

	errorsList := make([]string, 0, MaxDepth)
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		if goroot != "" && strings.Contains(frame.File, goroot) {
			continue
		}

		if strings.Contains(frame.File, packageName) {
			continue
		}

		// if strings.HasSuffix(frame.File, "/try/try.go") {
		// 	continue
		// }

		filename := strings.TrimPrefix(frame.File, goroot)

		// TODO: add lib to skip errors
		errorsList = append(errorsList, fmt.Sprintf("%s:%s:%d", frame.Function, filename, frame.Line))
	}
	return errorsList
}

type fake struct{}

var (
	goroot      = runtime.GOROOT()
	packageName = reflect.TypeOf(fake{}).PkgPath()
)

// Not sure this is useful, just experiment
func Wrap1[T1 any](v T1, err error) error {
	return Wrap(err)
}

func Wrap2[T1, T2 any](_ T1, _ T2, err error) error {
	return Wrap(err)
}

func Wrap3[T1, T2, T3 any](_ T1, _ T2, _ T3, err error) error {
	return Wrap(err)
}

func Wrap4[T1, T2, T3, T4 any](_ T1, _ T2, _ T3, _ T4, err error) error {
	return Wrap(err)
}

func Wrap5[T1, T2, T3, T4, T5 any](_ T1, _ T2, _ T3, _ T4, _ T5, err error) error {
	return Wrap(err)
}
