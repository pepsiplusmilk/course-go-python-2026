package errors

// AdditionalMessageError ошибка, которая
type AdditionalMessageError struct {
	// TODO: добавьте сюда свои поля
}

func (a AdditionalMessageError) Error() string {
	//TODO implement me
	panic("implement me")
}

// NewAdditionalMessageError создает ошибку
func NewAdditionalMessageError(err error, format string, args ...any) error {
	// TODO: добавьте сюда создание ошибки
	panic("unimplemented")
}

// TODO: добавьте сюда необходимые методы
