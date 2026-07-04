// Package errors is a drop-in replacement for the unpublished internal
// chaitin.cn/dev/go/errors module. It implements only the functions this
// codebase actually calls (New, Wrap, Wrapf, Annotate, Annotatef), each
// wrapping the underlying error with fmt.Errorf's %w so errors.Is/As still
// work against the original error.
package errors

import (
	"errors"
	"fmt"
)

// New returns an error with the given message.
func New(message string) error {
	return errors.New(message)
}

// Wrap annotates err with a message. Returns nil if err is nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf annotates err with a formatted message. Returns nil if err is nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Annotate annotates err with a message. Returns nil if err is nil.
func Annotate(err error, message string) error {
	return Wrap(err, message)
}

// Annotatef annotates err with a formatted message. Returns nil if err is nil.
func Annotatef(err error, format string, args ...interface{}) error {
	return Wrapf(err, format, args...)
}
