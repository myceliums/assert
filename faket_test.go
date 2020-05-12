package assert

type fakeT struct {
	gotError bool
}

func newFakeT() (*fakeT, Assert) {
	var ft fakeT
	return &ft, New(&ft)
}

func (t *fakeT) Error(_ ...interface{}) {
	t.gotError = true
}

func (_ *fakeT) Helper()          {}
func (_ *fakeT) Cleanup(_ func()) {}

func (t *fakeT) GotError() bool {
	r := t.gotError
	t.gotError = false

	return r
}
