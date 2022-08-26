package mock

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/asynkron/protoactor-go/actor"
)

func MakeFuture(result interface{}, err error) *actor.Future {
	future := new(actor.Future)
	setUnexportedField(
		reflect.ValueOf(future).
			Elem().
			FieldByName("result"),
		result,
	)
	setUnexportedField(
		reflect.ValueOf(future).
			Elem().
			FieldByName("err"),
		err,
	)
	setUnexportedField(
		reflect.ValueOf(future).
			Elem().
			FieldByName("done"),
		true,
	)
	setUnexportedField(
		reflect.ValueOf(future).
			Elem().
			FieldByName("cond"),
		sync.NewCond(&sync.Mutex{}),
	)

	return future
}

func setUnexportedField(field reflect.Value, value interface{}) {
	if value != nil {
		reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
			Elem().
			Set(reflect.ValueOf(value))
	}
}
