package handler

import (
	"context"
	"errors"

	bderrors "github.com/tcarzverey/course-go-python/homeworks/hw1/handler/errors"
)

type Handler struct {
	db UsersDB
}

func NewHandler(db UsersDB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) UpdateUserBalance(ctx context.Context, userID, balance int64) error {
	err := h.db.UpdateBalance(ctx, userID, balance)

	if err == nil {
		return nil
	}

	var retryableError *bderrors.RetryableError

	if errors.As(err, &retryableError) {
		iterations := retryableError.RetryCount()

		for ; iterations != 0; iterations-- {
			err = h.db.UpdateBalance(ctx, userID, balance)
			if err == nil {
				return nil
			}
		}

		return err
	}

	var notFoundError *bderrors.NotFoundError

	if errors.As(err, &notFoundError) {
		nfErrorWrapped := bderrors.NewAdditionalMessageError(notFoundError, "new error")
		return nfErrorWrapped
	}

	var additionalMessageError *bderrors.AdditionalMessageError

	if errors.As(err, &additionalMessageError) {
		return err
	}

	panic(err)
}
