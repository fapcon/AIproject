package mocks

import "github.com/stretchr/testify/mock"

func InOrder(calls ...*mock.Call) {
	for i, call := range calls[1:] {
		call.NotBefore(calls[i])
	}
}
