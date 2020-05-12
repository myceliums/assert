package assert

// T is an interface of what we use from testing.T
type T interface {
	Error(...interface{})
	Helper()
	Cleanup(func())
}
