package assert

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// Assert is a helper for tests
type Assert struct {
	t T
}

func (a Assert) f(ok bool, msg []interface{}, format string, extra ...interface{}) {
	if !ok {
		if format != `` {
			msg = prepMsg(msg, format, extra...)
		}
		if msg == nil {
			msg = []interface{}{`Assertion failed`}
		}
		msg = append(append([]interface{}{shell(1) + shell(97) + shell(41) + `FAIL!` + shell(0) + shell(1)}, msg...), shell(0), "\n")
		a.t.Helper()
		a.t.Error(msg...)
	}
}

// New returns a new Assert
func New(t T) Assert {
	a := Assert{t}

	return a
}

// TODO: Panics, Len

///// Boolean /////

// True asserts the given value is true
func (a Assert) True(actual bool, msg ...interface{}) {
	a.t.Helper()
	a.f(actual, msg, `Should be true, but it isn't`)
}

// False sserts the given value is false
func (a Assert) False(actual bool, msg ...interface{}) {
	a.t.Helper()
	a.f(!actual, msg, `Should be false, but it isn't`)
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
	a.t.Helper()
	a.f(isNil(actual), msg, `Should be nil, but got %#v`, actual)
}

// NotNil sserts the given value is not nil
func (a Assert) NotNil(actual interface{}, msg ...interface{}) {
	a.t.Helper()
	a.f(!isNil(actual), msg, `Should not be nil, but it is`)
}

///// Errors /////

// Error asserts the given error is not nil
func (a Assert) Error(actual error, msg ...interface{}) {
	a.t.Helper()
	a.f(actual != nil, msg, `Expected an error, but got nil`)
}

// NoError asserts the given error is not nil
func (a Assert) NoError(actual error, msg ...interface{}) {
	a.t.Helper()
	a.f(actual == nil, msg, `Expected no error, but got %#v`, actual)
}

///// Comparisons /////

// Eq asserts the given values match
func (a Assert) Eq(expected, actual interface{}, msg ...interface{}) {
	a.t.Helper()
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		a.f(false, msg, `Expected %T(%#v), but got %T(%#v)`, expected, expected, actual, actual)
		return
	}

	a.f(expected == actual, msg, `Expected %#v, but got %#v`, expected, actual)
}

// Ne asserts the given values don't match
func (a Assert) Ne(expected, actual interface{}, msg ...interface{}) {
	a.t.Helper()
	a.f(expected != actual, msg, `Should not be %#v, but it is`, expected)
}

///// Lists /////

// Contains asserts the expected value is in the given list
func (a Assert) Contains(expected, list interface{}, msg ...interface{}) {
	a.t.Helper()

	rlist := reflect.ValueOf(list)
	a.f(rlist.Kind() == reflect.Slice || rlist.Kind() == reflect.Array, nil, `Can only call assert.Contains on a slice or array`)
	for i := 0; i < rlist.Len(); i++ {
		if rlist.Index(i).Interface() == expected {
			return
		}
	}

	a.f(false, msg, `Expected %#v to be in %#v, but it isn't`, expected, list)
}

// SameElements asserts the values have the same elements. It ignores the order of the elements
func (a Assert) SameElements(expected, actual interface{}, msg ...interface{}) {
	a.t.Helper()

	rexpected, ractual := reflect.ValueOf(expected), reflect.ValueOf(actual)
	a.f(rexpected.Kind() == reflect.Slice || rexpected.Kind() == reflect.Array, nil, `Can only call assert.SameElements on a slice or array`)
	a.f(ractual.Kind() == reflect.Slice || ractual.Kind() == reflect.Array, nil, `Can only call assert.SameElements on a slice or array`)

	if rexpected.Len() != ractual.Len() {
		a.f(false, msg, `Expected elements of %#v to match %#v, but they don't`, expected, actual)
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

	a.f(false, msg, ``)
}

// Cmp assert wrapper for go-cmp
func (a Assert) Cmp(expected, actual interface{}, opts ...cmp.Option) {
	a.t.Helper()
	diff := cmp.Diff(expected, actual, opts...)
	if diff == `` {
		return
	}

	a.f(false, nil, "\n"+diff)
}

// NCmp assert wrapper for go-cmp but fails when !Equal
func (a Assert) NCmp(expected, actual interface{}, opts ...cmp.Option) {
	a.t.Helper()

	ok := cmp.Equal(expected, actual, opts...)
	a.f(!ok, nil, `Should not be %#v, but it is`, expected)
}
