package sqlm

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	full  string
	wrong string
	desc  string
}

func (e *Error) Error() string {
	start := strings.Index(e.full, e.wrong)
	end := start + len(e.wrong)
	var sql string
	if start > 0 {
		if end > len(e.full) {
			end = len(e.full)
		}
		sql = e.full[:start]
		sql += mark(e.full[start:end])
		if end < len(e.full) {
			sql += e.full[end:]
		}
	}
	if sql == "" {
		sql = e.full
	}
	return fmt.Sprintf("%s:\n%s", red(e.desc), sql)
}

type MethodError struct {
	Method string
	Err    error
}

func (e *MethodError) Error() string {
	return fmt.Sprintf("parse method %s failed: %v",
		cyan(e.Method), e.Err)
}

var colorEnable = func() bool {
	return runtime.GOOS == "darwin" || runtime.GOOS == "linux"
}()

func mark(s string) string {
	if !colorEnable {
		// mark donot support in windows os.
		return s
	}
	return fmt.Sprintf("\033[31;4m%s\033[0m", s)
}

func red(s string) string {
	if !colorEnable {
		return s
	}
	return fmt.Sprintf("\033[31;1m%s\033[0m", s)
}

func cyan(s string) string {
	if !colorEnable {
		return s
	}
	return fmt.Sprintf("\033[36;1m%s\033[0m", s)
}
