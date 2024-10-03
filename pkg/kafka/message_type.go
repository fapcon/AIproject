package kafka

import "reflect"

func messageType[T any]() string {
	// Following code returns "<nil>" if T is an interface:
	//   return fmt.Sprintf("%T", zero)
	// So need to hack it.
	x := reflect.TypeOf(new(T))
	for x.Kind() == reflect.Pointer {
		x = x.Elem()
	}
	return x.String()
	// An alternatives:
	//   return fmt.Sprintf("%T", []T(nil))[2:]
	//   return fmt.Sprintf("%T", new(*T))[2:]
	//   return fmt.Sprintf("%T", *new(*T))[1:]
}
