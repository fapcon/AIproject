//nolint:errname // ok
package errors

import (
	"fmt"
)

type Fields = map[string]any

func Errors(err error) []error {
	errs := wrapper{err: err}.Errors()
	res := make([]error, 0, len(errs))
	for _, e := range errs {
		res = append(res, e)
	}
	return res
}

func FieldsFromError(err error) Fields {
	errs := wrapper{err: err}.Errors()
	if len(errs) == 0 || errs[0].fields == nil {
		// It may be handy to update map returned by FieldsFromError.
		// Need to return empty map instead of nil.
		return Fields{}
	}
	return clone(errs[0].fields)
}

type ErrorBuilder struct {
	err treeNode
}

func Err(err error) *ErrorBuilder {
	// Type switch instead of errors.As() because we don't want to extract wrapped error to not miss wrapper.
	switch e := err.(type) { //nolint:errorlint // see comment above
	case nil:
		return nil
	case treeNode:
		return &ErrorBuilder{
			err: e,
		}
	default:
		return &ErrorBuilder{
			err: wrapper{
				err: err,
			},
		}
	}
}

func (e *ErrorBuilder) E() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *ErrorBuilder) Wrap(prefix string) *ErrorBuilder {
	if e == nil {
		return nil
	}
	e.err = withPrefix{
		err:    e.err,
		prefix: prefix,
	}
	return e
}

func (e *ErrorBuilder) WithFields(fields Fields) *ErrorBuilder {
	if e == nil {
		return nil
	}
	e.err = withFields{
		err:    e.err,
		fields: clone(fields),
	}
	return e
}

func (e *ErrorBuilder) WithField(key string, value any) *ErrorBuilder {
	return e.WithFields(Fields{
		key: value,
	})
}

func Wrap(prefix string, err error) *ErrorBuilder {
	return Err(err).Wrap(prefix)
}

func joinFields(outer Fields, inner Fields) Fields {
	switch {
	case len(outer) == 0:
		return inner
	case len(inner) == 0:
		return outer
	}
	res := make(Fields, len(outer)+len(inner))
	for k, v := range outer {
		res[k] = v
	}
	for k, v := range inner {
		res[k] = v // inner map has higher priority for duplicated fields
	}
	return res
}

type errorWithFields struct {
	err    error
	fields Fields
}

func (e errorWithFields) Error() string {
	return e.err.Error()
}

func (e errorWithFields) Unwrap() error {
	return e.err
}

type treeNode interface {
	isMyError()
	error
	Errors() []errorWithFields
	// Optional methods:
	//   Unwrap() error
	//   Unwrap() []error
}

//nolint:exhaustruct // false positive
var (
	_ treeNode = wrapper{}
	_ treeNode = withPrefix{}
	_ treeNode = withFields{}
)

func (e wrapper) isMyError()    {}
func (e withPrefix) isMyError() {}
func (e withFields) isMyError() {}

type wrapper struct {
	err error
}

func (e wrapper) Errors() []errorWithFields {
	// Type switch instead of errors.As() because we don't want to extract wrapped error to not miss wrapper.
	switch err := e.err.(type) { //nolint:errorlint // see comment above
	case nil:
		return nil
	case treeNode:
		return err.Errors()
	case errorWithFields:
		return []errorWithFields{err}
	default:
		return []errorWithFields{{
			err:    err,
			fields: nil,
		}}
	}
}

func (e wrapper) Error() string {
	return e.err.Error()
}

func (e wrapper) Unwrap() error {
	return e.err
}

type withPrefix struct {
	err    treeNode
	prefix string
}

func (e withPrefix) Errors() []errorWithFields {
	errs := e.err.Errors()
	res := make([]errorWithFields, 0, len(errs))
	for _, err := range errs {
		res = append(res, errorWithFields{
			err:    fmt.Errorf("%s: %w", e.prefix, err.err),
			fields: err.fields,
		})
	}
	return res
}

func (e withPrefix) Error() string {
	return fmt.Sprintf("%s: %s", e.prefix, e.err.Error())
}

func (e withPrefix) Unwrap() error {
	return e.err
}

type withFields struct {
	err    treeNode
	fields Fields
}

func (e withFields) Errors() []errorWithFields {
	errs := e.err.Errors()
	res := make([]errorWithFields, 0, len(errs))
	for _, err := range errs {
		res = append(res, errorWithFields{
			err:    err.err,
			fields: joinFields(e.fields, err.fields),
		})
	}
	return res
}

func (e withFields) Error() string {
	return e.err.Error()
}

func (e withFields) Unwrap() error {
	return e.err
}

func clone(fields Fields) Fields {
	res := make(Fields, len(fields))
	for k, v := range fields {
		res[k] = v
	}
	return res
}
