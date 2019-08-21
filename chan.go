package gomads

import (
	"reflect"
)

type _ChanBoxed struct {
	T reflect.Type
	V reflect.Value
}

func newChanBoxed(v interface{}) Boxed {
	slice, err := baseType(reflect.TypeOf(v), reflect.Chan)
	panicCondition(err != nil, "gomads: chanBox got no channel ("+reflect.TypeOf(v).String()+")")
	b := _ChanBoxed{
		T: slice.Elem(),
		V: reflect.ValueOf(v),
	}
	return &b
}

func (b *_ChanBoxed) Map(fmap interface{}) Boxed {
	ft := reflect.TypeOf(fmap)
	panicCondition(ft.Kind() != reflect.Func, "gomads: Map (not function)")
	panicCondition(ft.NumIn() != 1, "gomads: Map (need one input param)")
	panicCondition(ft.NumOut() != 1, "gomads: Map (need one output param)")
	panicCondition(ft.In(0) != b.T, "gomads: Map (not same input type for fmap)")

	fv := reflect.ValueOf(fmap)
	outChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, ft.Out(0)), b.V.Cap())
	go func() {
		for v, ok := b.V.Recv(); ok; v, ok = b.V.Recv() {
			results := fv.Call([]reflect.Value{v})
			outChan.Send(results[0])
		}
		outChan.Close()
	}()
	return newChanBoxed(outChan.Interface())
}

func (b *_ChanBoxed) FlatMap(fmap interface{}) Boxed {
	ft := reflect.TypeOf(fmap)
	panicCondition(ft.Kind() != reflect.Func, "gomads: FlatMap (not function)")
	panicCondition(ft.NumIn() != 1, "gomads: FlatMap (need one input param)")
	panicCondition(ft.NumOut() != 1, "gomads: FlatMap (need one output param)")
	panicCondition(ft.In(0) != b.T, "gomads: FlatMap (not same input type for fmap)")

	return b
}

func (b *_ChanBoxed) Unbox(dest interface{}) {
	value := reflect.ValueOf(dest)
	panicCondition(value.Kind() != reflect.Ptr, "gomads: Unbox needs a pointer")
	panicCondition(value.IsNil(), "gomads: Unbox (nil)")
	ch, err := baseType(value.Type(), reflect.Chan)
	panicCondition(err != nil, "gomads: Unbox (not channel destination)")
	panicCondition(ch.Elem() != b.T, "gomads: Unbox (not same element type)")
	reflect.Indirect(value).Set(b.V)
}
