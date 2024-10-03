package app

import (
	"context"
	"reflect"

	"studentgit.kata.academy/quant/torque/internal/domain"

	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/errors"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

func process[T any](
	logger logster.Logger,
	name string,
	ch <-chan T,
	fn func(context.Context, T) error,
) func() error {
	return func() error {
		for value := range ch {
			ctx := logster.Span(context.Background(), name)
			errX := fn(ctx, value)

			errX = errors.Wrap(name, errX).
				WithField(f.Message, value).
				WithField(f.MessageType, messageType(value)).E()

			for _, err := range errors.Errors(errX) {
				logger.
					WithContext(ctx).
					WithFields(errors.FieldsFromError(err)).
					WithError(err).
					Errorf("%s error", name)
			}
		}
		logger.Infof("%s stopped", name)
		return nil
	}
}

func processAckable[T any](
	logger logster.Logger,
	name string,
	ch <-chan domain.Ackable[T],
	fn func(context.Context, T) error,
) func() error {
	return func() error {
		for value := range ch {
			ctx := logster.Span(context.Background(), name)

			errX := fn(ctx, value.Value)

			for _, err := range errors.Errors(errX) {
				logger.
					WithContext(ctx).
					WithFields(errors.FieldsFromError(err)).
					WithError(err).
					Errorf("%s error", name)
			}

			value.Ack()
		}
		return nil
	}
}

func messageType[T any](v T) string {
	x := reflect.TypeOf(&v)
	for x.Kind() == reflect.Pointer {
		x = x.Elem()
	}
	return x.String()
}
