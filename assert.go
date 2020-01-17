package assert

import (
	"reflect"
	"runtime"
)

// Assert is a helper for tests
type Assert func(bool, ...interface{})

// New returns a new Assert
func New(t T) Assert {
	a := func(ok bool, msg ...interface{}) {
		if !ok {
			if msg == nil {
				msg = []interface{}{`Assertion failed`}
			}
			msg = append(append([]interface{}{shell(1) + shell(97) + shell(41) + `FAIL!` + shell(0) + shell(1)}, msg...), shell(0), "\n")
			t.Helper()
			t.Error(msg...)
		}
	}

	f := runtime.FuncForPC(reflect.ValueOf(a).Pointer())
	ts[f] = t

	return a
}

// TODO: Panics, Len

///// Boolean /////

// True asserts the given value is true
func (a Assert) True(actual bool, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Should be true, but it isn't`)
	a(actual, msg...)
}

// False sserts the given value is false
func (a Assert) False(actual bool, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Should be false, but it isn't`)
	a(!actual, msg...)
}

///// Nil /////

func isNil(val interface{}) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}

// Nil asserts the given value is nil
func (a Assert) Nil(actual interface{}, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Should be nil, but got %#v`, actual)

	a(isNil(actual), msg...)
}

// NotNil sserts the given value is not nil
func (a Assert) NotNil(actual interface{}, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Should not be nil, but it is`)

	a(!isNil(actual), msg...)
}

///// Errors /////

// Error asserts the given error is not nil
func (a Assert) Error(actual error, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Expected an error, but got nil`)
	a(actual != nil, msg...)
}

// NoError asserts the given error is not nil
func (a Assert) NoError(actual error, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Expected no error, but got %#v`, actual)
	a(actual == nil, msg...)
}

///// Comparisons /////

// Eq asserts the given values match
func (a Assert) Eq(expected, actual interface{}, msg ...interface{}) {
	t(a).Helper()
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		msg = prepMsg(msg, `Expected %T(%#v), but got %T(%#v)`, expected, expected, actual, actual)
		a(false, msg...)
		return
	}

	msg = prepMsg(msg, `Expected %#v, but got %#v`, expected, actual)
	a(expected == actual, msg...)
}

// Ne asserts the given values don't match
func (a Assert) Ne(expected, actual interface{}, msg ...interface{}) {
	t(a).Helper()
	msg = prepMsg(msg, `Should not be %#v, but it is`, expected)
	a(expected != actual, msg...)
}

///// Lists /////

// Contains asserts the expected value is in the given list
func (a Assert) Contains(expected, list interface{}, msg ...interface{}) {
	t(a).Helper()

	rlist := reflect.ValueOf(list)
	a(rlist.Kind() == reflect.Slice || rlist.Kind() == reflect.Array, `Can only call assert.Contains on a slice or array`)
	for i := 0; i < rlist.Len(); i++ {
		if rlist.Index(i).Interface() == expected {
			return
		}
	}

	msg = prepMsg(msg, `Expected %#v to be in %#v, but it isn't`, expected, list)
	a(false, msg...)
}

// SameElements asserts the values have the same elements. It ignores the order of the elements
func (a Assert) SameElements(expected, actual interface{}, msg ...interface{}) {
	t(a).Helper()

	rexpected, ractual := reflect.ValueOf(expected), reflect.ValueOf(actual)
	a(rexpected.Kind() == reflect.Slice || rexpected.Kind() == reflect.Array, `Can only call assert.SameElements on a slice or array`)
	a(ractual.Kind() == reflect.Slice || ractual.Kind() == reflect.Array, `Can only call assert.SameElements on a slice or array`)

	msg = prepMsg(msg, `Expected elements of %#v to match %#v, but they don't`, expected, actual)
	if rexpected.Len() != ractual.Len() {
		a(false, msg...)
		return
	}

	var same int
	for i := 0; i < rexpected.Len(); i++ {
		for j := 0; j < ractual.Len(); j++ {
			if rexpected.Index(i).Interface() == ractual.Index(j).Interface() {
				same++
				break
			}
		}
	}

	if same == rexpected.Len() {
		return
	}

	a(false, msg...)
}
