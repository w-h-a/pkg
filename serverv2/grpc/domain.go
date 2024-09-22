package grpc

import (
	"reflect"
)

type Handler struct {
	Name     string
	Receiver reflect.Value
	Methods  map[string]*Method
}

type Method struct {
	Name    string
	Value   reflect.Value
	CtxType reflect.Type
	ReqType reflect.Type
	RspType reflect.Type
	Stream  bool
}

func NewHandler(handler interface{}) *Handler {
	methods := map[string]*Method{}

	// used to get method data
	typeOfHandler := reflect.TypeOf(handler)

	for i := 0; i < typeOfHandler.NumMethod(); i++ {
		m := typeOfHandler.Method(i)

		method := &Method{
			Name:  m.Name,
			Value: m.Func,
		}

		// TODO: make this better/safer
		switch m.Type.NumIn() {
		case 3:
			method.CtxType = m.Type.In(1)
			method.RspType = m.Type.In(2)
			method.Stream = true
		case 4:
			method.CtxType = m.Type.In(1)
			method.ReqType = m.Type.In(2)
			method.RspType = m.Type.In(3)
		}

		methods[m.Name] = method
	}

	// keep the value to use as a receiver in function invocation
	valueOfHandler := reflect.ValueOf(handler)

	// get the name of the handler struct
	nameOfHandler := reflect.Indirect(valueOfHandler).Type().Name()

	h := &Handler{
		Name:     nameOfHandler,
		Receiver: valueOfHandler,
		Methods:  methods,
	}

	return h
}
