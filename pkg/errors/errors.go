package errors

import (
	"github.com/cockroachdb/errors"
)

// New is an alias for Error
func New(msg string) error {
	return errors.NewWithDepth(1, msg)
}

// Newf is an alias for Errorf
func Newf(format string, args ...any) error {
	return errors.NewWithDepthf(1, format, args...)
}

// Error creates an error with a simple error message.
// A stack trace is retained.
func Error(msg string) error {
	return errors.NewWithDepth(1, msg)
}

// Errorf creates an error with a formatted error message.
// A stack trace is retained.
func Errorf(format string, args ...any) error {
	return errors.NewWithDepthf(1, format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
func WithStack(err error) error {
	return errors.WithStackDepth(err, 1)
}

// WithMessage annotates err with a new message.
func WithMessage(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

// WithMessagef annotates err with the format specifier.
func WithMessagef(err error, format string, args ...any) error {
	return errors.WithMessagef(err, format, args...)
}

// Wrap wraps an error with a message prefix.
// A stack trace is retained.
func Wrap(err error, msg string) error {
	return errors.WrapWithDepth(1, err, msg)
}

// Wrapf wraps an error with a formatted message prefix.
// A stack trace is also retained.
func Wrapf(err error, format string, args ...any) error {
	return errors.WrapWithDepthf(1, err, format, args...)
}

// Unwarp accesses the direct cause of the error if any, otherwise
// returns nil.
func Unwarp(err error) error {
	return errors.Unwrap(err)
}

// UnwrapAll accesses the root cause object of the error.
func UnwrapAll(err error) error {
	return errors.UnwrapAll(err)
}

// As finds the first error in errs chain that matches the type to which target
// points, and if so, sets the target to its value and returns true.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
func Join(errs ...error) error {
	return errors.JoinWithDepth(1, errs...)
}

// Is determines whether one of the causes of the given error or any
// of its causes is equivalent to some reference error.
func Is(err, reference error) bool {
	return errors.Is(err, reference)
}

// IsAny is like Is except that multiple references are compared.
func IsAny(err error, reference ...error) bool {
	return errors.IsAny(err, reference...)
}
