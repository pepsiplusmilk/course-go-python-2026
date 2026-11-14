package errors

import "fmt"

// AdditionalMessageError ошибка, которая
type AdditionalMessageError struct {
	wrappedError error
	newMessage   string
}

func (a *AdditionalMessageError) Error() string {
	if a == nil {
		return ""
	}

	newPart := a.newMessage
	oldPart := ""

	if a.wrappedError != nil {
		oldPart = a.wrappedError.Error()
	}

	if oldPart != "" && newPart != "" {
		newPart += ": "
	}

	return fmt.Sprintf("%s%s", newPart, oldPart)
}

// NewAdditionalMessageError создает ошибку
func NewAdditionalMessageError(err error, format string, args ...any) error {
	newMessage := ""

	if format != "" {
		newMessage = fmt.Sprintf(format, args...)
	}

	return &AdditionalMessageError{
		wrappedError: err,
		newMessage:   newMessage,
	}
}

// TODO: добавьте сюда необходимые методы

func (a *AdditionalMessageError) Unwrap() error {
	if a == nil {
		return nil
	}

	return a.wrappedError
}
