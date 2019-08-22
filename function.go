package gomads

import (
	"reflect"
)

type _FunctionBoxed struct {
	T reflect.Type
	V reflect.Value
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func newFunctionBoxed(v interface{}) Boxed {
	fn, err := baseType(reflect.TypeOf(v), reflect.Func)
	panicCondition(err != nil, "gomads: functionBox got no function ("+reflect.TypeOf(v).String()+")")
	b := _FunctionBoxed{
		T: fn,
		V: reflect.ValueOf(v),
	}
	return &b
}

func (b *_FunctionBoxed) Map(fmap interface{}) Boxed {
	ft := reflect.TypeOf(fmap)
	panicCondition(ft.Kind() != reflect.Func, "gomads: Map (not function)")
	panicCondition(hasLastErrorResult(ft), "gomads: Map (cannot take function that returns error as the last output)")
	numOut := b.T.NumOut()
	hasLastError := hasLastErrorResult(b.T)
	if hasLastError {
		numOut--
	}

	panicCondition(numOut != ft.NumIn(), "gomads: Map (function has no type tranformation)")
	for i := 0; i < numOut; i++ {
		panicCondition(b.T.Out(i) != ft.In(i), "gomads: Map (function has no type transformation)")
	}

	inTypes := make([]reflect.Type, 0, b.T.NumIn())
	for i := 0; i < b.T.NumIn(); i++ {
		inTypes = append(inTypes, b.T.In(i))
	}
	outTypes := make([]reflect.Type, 0, ft.NumOut())
	outZeros := make([]reflect.Value, 0, ft.NumOut())
	for i := 0; i < ft.NumOut(); i++ {
		outTypes = append(outTypes, ft.Out(i))
		outZeros = append(outZeros, reflect.Zero(ft.Out(i)))
	}
	if hasLastError {
		outTypes = append(outTypes, errorInterface)
	}

	newFuncType := reflect.FuncOf(inTypes, outTypes, false)
	fv := reflect.ValueOf(fmap)
	fn := reflect.MakeFunc(newFuncType, func(in []reflect.Value) []reflect.Value {
		results := b.V.Call(in)
		if !hasLastError {
			return reflect.ValueOf(fmap).Call(results)
		}
		results, lastError := results[:len(results)-1], results[len(results)-1]
		if err, ok := lastError.Interface().(error); ok && err != nil {
			results = append(outZeros, lastError)
		} else {
			results = fv.Call(results)
			results = append(results, lastError)
		}

		return results
	})
	return newFunctionBoxed(fn.Interface())
}

func (b *_FunctionBoxed) FlatMap(fmap interface{}) Boxed {
	ft := reflect.TypeOf(fmap)
	panicCondition(ft.Kind() != reflect.Func, "gomads: FlatMap (not function)")
	panicCondition(!hasLastErrorResult(ft), "gomads: FlatMap (should take function that returns error as the last output)")
	numOut := b.T.NumOut()
	hasLastError := hasLastErrorResult(b.T)
	if hasLastError {
		numOut--
	}
	panicCondition(numOut != ft.NumIn(), "gomads: Map (function has no type tranformation)")
	for i := 0; i < numOut; i++ {
		panicCondition(b.T.Out(i) != ft.In(i), "gomads: Map (function has no type transformation)")
	}

	inTypes := make([]reflect.Type, 0, b.T.NumIn())
	for i := 0; i < b.T.NumIn(); i++ {
		inTypes = append(inTypes, b.T.In(i))
	}
	outTypes := make([]reflect.Type, 0, ft.NumOut())
	outZeros := make([]reflect.Value, 0, ft.NumOut())
	for i := 0; i < ft.NumOut(); i++ {
		outTypes = append(outTypes, ft.Out(i))
		outZeros = append(outZeros, reflect.Zero(ft.Out(i)))
	}
	newFuncType := reflect.FuncOf(inTypes, outTypes, false)
	fv := reflect.ValueOf(fmap)
	fn := reflect.MakeFunc(newFuncType, func(in []reflect.Value) []reflect.Value {
		results := b.V.Call(in)
		results, lastError := results[:len(results)-1], results[len(results)-1]
		if err, ok := lastError.Interface().(error); ok && err != nil {
			outZeros[len(outZeros)-1] = lastError
			results = outZeros
		} else {
			results = fv.Call(results)
		}
		return results
	})
	return newFunctionBoxed(fn.Interface())
}

func (b *_FunctionBoxed) Unbox(dest interface{}) {
	value := reflect.ValueOf(dest)
	panicCondition(value.Kind() != reflect.Ptr, "gomads: Unbox needs a pointer")
	panicCondition(value.IsNil(), "gomads: Unbox (nil)")

	fn, err := baseType(value.Type(), reflect.Func)
	panicCondition(err != nil, "gomads: Unbox (not function destination)")
	panicCondition(fn.NumIn() != b.T.NumIn(), "gomads: Unbox (not equal input params)")
	panicCondition(fn.NumOut() != b.T.NumOut(), "gomads: Unbox (not equal input params)")
	for i := 0; i < fn.NumIn(); i++ {
		panicCondition(fn.In(i) != b.T.In(i), "gomads: Unbox (not same function signature)")
	}
	for i := 0; i < fn.NumOut(); i++ {
		panicCondition(fn.Out(i) != b.T.Out(i), "gomads: Unbox (not same function signature)")
	}

	reflect.Indirect(value).Set(b.V)
}

func hasLastErrorResult(t reflect.Type) bool {
	panicCondition(t.Kind() != reflect.Func, "gomads: "+t.String()+" not a function")
	if t.NumOut() > 0 && t.Out(t.NumOut()-1).Implements(errorInterface) {
		return true
	}
	return false
}

// ComposeErrors return a composed function chain with monadic error handling
func ComposeErrors(fns ...interface{}) Boxed {
	panicCondition(len(fns) == 0, "gomads: ComposeErrros (should take at least one functions)")
	fnts := make([]reflect.Type, 0, len(fns))
	for _, fn := range fns {
		ft := reflect.TypeOf(fn)
		panicCondition(ft.Kind() != reflect.Func, "gomads: ComposeErrors (take only functions as input params)")
		fnts = append(fnts, ft)
	}

	boxed := Box(fns[0])
	for i := 1; i < len(fns); i++ {
		ft := fnts[i]
		if hasLastErrorResult(ft) {
			boxed = boxed.FlatMap(fns[i])
		} else {
			boxed = boxed.Map(fns[i])
		}
	}
	return boxed
}
