package gomads

import (
	"reflect"
)

type _SliceBoxed struct {
	T reflect.Type
	V reflect.Value
}

func newSliceBoxed(v interface{}) Boxed {
	slice, err := baseType(reflect.TypeOf(v), reflect.Slice)
	if err != nil {
		panic("gomads: sliceBox got no slice (" + reflect.TypeOf(v).String() + ")")
	}
	b := _SliceBoxed{
		T: slice.Elem(),
		V: reflect.ValueOf(v),
	}
	return &b
}

func (b *_SliceBoxed) Map(fmap interface{}) Boxed {
	ft := reflect.TypeOf(fmap)
	panicCondition(ft.Kind() != reflect.Func, "gomads: Map (not function)")
	panicCondition(ft.NumIn() != 1, "gomads: Map (need one input param)")
	panicCondition(ft.NumOut() != 1, "gomads: Map (need one output param)")

	fv := reflect.ValueOf(fmap)
	outs := reflect.MakeSlice(reflect.SliceOf(ft.Out(0)), 0, b.V.Len())
	for i := 0; i < b.V.Len(); i++ {
		call := fv.Call([]reflect.Value{b.V.Index(i)})
		outs = reflect.Append(outs, call[0])
	}
	return newSliceBoxed(outs.Interface())
}

func (b *_SliceBoxed) Unbox(dest interface{}) {
	value := reflect.ValueOf(dest)
	panicCondition(value.Kind() != reflect.Ptr, "gomads: Unbox needs a pointer")
	panicCondition(value.IsNil(), "gomads: Unbox (nil)")

	slice, err := baseType(value.Type(), reflect.Slice)
	panicCondition(err != nil, "gomads: Unbox (not slice destination)")
	panicCondition(slice.Elem() != b.T, "gomads: Unbox (not same element type)")

	direct := reflect.Indirect(value)
	for i := 0; i < b.V.Len(); i++ {
		direct.Set(reflect.Append(direct, b.V.Index(i)))
	}
}