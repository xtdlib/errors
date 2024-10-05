package throw

import (
	"encoding/json"
	"errors"
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

func SlogAttr(err error) slog.Attr {
	return slog.Any("throw", Wrap(err))
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var terr ThrowError

	// do not re-wrap
	if errors.As(err, &terr) {
		terr.err = err
		return terr
	}

	return ThrowError{err: err, stacktrace: getStackTrace()}
}

func getStackTrace() []string {
	stackBuffer := make([]uintptr, MaxDepth)
	length := runtime.Callers(3, stackBuffer[:])
	stack := stackBuffer[:length]

	throwList := make([]string, 0, MaxDepth)
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

		if strings.HasSuffix(frame.File, "/try/try.go") {
			continue
		}

		// TODO: add lib to skip throw
		throwList = append(throwList, fmt.Sprintf("%s:%s:%d", frame.Function, frame.File, frame.Line))
	}
	return throwList
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
