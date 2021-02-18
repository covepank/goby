package errors

import (
	"fmt"
	"io"
)

type wrapError struct {
	msg   string
	cause error
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &wrapError{
		cause: err,
		msg:   msg,
	}
}

func (w *wrapError) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *wrapError) Unwrap() error { return w.cause }
func (w *wrapError) Cause() error  { return w.cause }

func (w *wrapError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
