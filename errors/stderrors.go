package errors

import stderrors "errors"

func New(text string) error { return stderrors.New(text) }

func Is(err, target error) bool { return stderrors.Is(err, target) }

func As(err error, target interface{}) bool { return stderrors.As(err, target) }

func Unwrap(err error) error { return stderrors.Unwrap(err) }
