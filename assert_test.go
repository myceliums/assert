package assert

import (
	"testing"
	"time"
)

// -----------------------------------------
//  Slices
// -----------------------------------------

func TestSameElementsSlices(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	a, b := []int{1, 2}, []int{1, 2}
	fa.SameElements(a, a)
	assert.False(ft.GotError())

	fa.SameElements(a, b)
	assert.False(ft.GotError())

	b = []int{2, 1}
	fa.SameElements(b, a) // SameElements ignores order
	assert.False(ft.GotError())

	b = []int{1, 2, 3}
	fa.SameElements(b, a)
	assert.True(ft.GotError())

	b = []int{2, 3}
	fa.SameElements(b, a)
	assert.True(ft.GotError())
}

func TestCmpSlices(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	a, b := []int{1, 2}, []int{1, 2}
	fa.Cmp(a, a)
	assert.False(ft.GotError())

	fa.Cmp(a, b)
	assert.False(ft.GotError())

	b = []int{2, 1}
	fa.Cmp(a, b) // Cmp does not accept different order
	assert.True(ft.GotError())

	b = []int{1, 2, 3}
	fa.Cmp(a, b)
	assert.True(ft.GotError())

	b = []int{2, 3}
	fa.Cmp(a, b)
	assert.True(ft.GotError())
}

// -----------------------------------------
//  Maps
// -----------------------------------------

func TestCmpMaps(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	a, b := map[int]int{1: 2, 3: 4}, map[int]int{1: 2, 3: 4}
	fa.Cmp(a, a)
	assert.False(ft.GotError())

	fa.Cmp(a, b)
	assert.False(ft.GotError())

	b = map[int]int{1: 2, 3: 5}
	fa.Cmp(a, b)
	assert.True(ft.GotError())

	b = map[int]int{1: 2, 3: 4, 5: 6}
	fa.Cmp(a, b)
	assert.True(ft.GotError())

	b = map[int]int{1: 2}
	fa.Cmp(a, b)
	assert.True(ft.GotError())
}

// -----------------------------------------
//  Pointers
// -----------------------------------------

func TestEqPointers(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	var a, b int
	fa.Eq(&a, &a)
	assert.False(ft.GotError())

	fa.Eq(&a, &b)
	assert.True(ft.GotError())

	fa.Eq(&b, &a)
	assert.True(ft.GotError())
}

func TestCmpPointers(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	type B struct {
		Value int
	}
	type A struct {
		B *B
	}

	fa.Cmp(A{B: &B{1}}, A{&B{1}})
	assert.False(ft.GotError())

	fa.Cmp(A{B: &B{1}}, A{&B{2}})
	assert.True(ft.GotError())

	fa.Cmp(A{B: &B{1}}, A{nil})
	assert.True(ft.GotError())
}

// -----------------------------------------
//  Timezones
// -----------------------------------------

func mustLoadLocation(zone string) *time.Location {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		panic(err)
	}

	return loc
}

var locAmsterdam = mustLoadLocation(`Europe/Amsterdam`)
var locTokyo = mustLoadLocation(`Asia/Tokyo`)

func TestEqTimezones(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	ti := time.Now()
	ti, ti2 := ti.In(locAmsterdam), ti.In(locTokyo)

	fa.Eq(ti, ti2)
	assert.True(ft.GotError())
}

func TestCmpTimezones(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	type A struct {
		Time time.Time
	}

	ti := time.Now()
	ti2 := ti.In(time.UTC)
	fa.Cmp(&A{ti}, &A{ti2})
	assert.False(ft.GotError())

	ti2 = ti.In(locAmsterdam)
	fa.Cmp(&A{ti}, &A{ti2})
	assert.False(ft.GotError())

	ti = ti.In(locTokyo)
	fa.Cmp(&A{ti}, &A{ti2})
	assert.False(ft.GotError())

	ti2 = ti2.Add(time.Second)
	fa.Cmp(&A{ti}, &A{ti2})
	assert.True(ft.GotError())
}

func TestNCmp(t *testing.T) {
	assert := New(t)
	ft, fa := newFakeT()

	type A struct {
		S string
	}

	fa.NCmp(A{"not"}, A{"equal"})
	assert.False(ft.GotError())

	fa.NCmp(A{"equal"}, A{"equal"})
	assert.True(ft.GotError())
}
