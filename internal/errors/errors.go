package errors

import (
	"errors"
	"fmt"
)

var ErrUnsupported = errors.ErrUnsupported

func New(text string) error {
	return errors.New(text)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Swrap(err error, msg string) error {
	return Swrapf(err, "%s", msg)
}

func Swrapf(err error, format string, a ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), err)
}

func Wrap(target error, wrapper error) error {
	return fmt.Errorf("%w: %w", wrapper, target)
}
