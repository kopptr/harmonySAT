package assignment

import (
	"testing"
   "dpll/assignment/guess"
)

func TestPushPopAssign(t *testing.T) {
	a := NewAssignment(10)

	if a.Depth() != 0 {
		t.Log("Depth initializes improperly\n")
		t.Fail()
	}

	a.PushAssign(2, guess.Pos)
	if a.Depth() != 1 {
		t.Log("Depth updates improperly\n")
		t.Fail()
	}
	if p := a.Get(2); p != guess.Pos {
		t.Logf("PushAssigned 2, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 2, p)
		t.Fail()
	}
	a.PushAssign(3, guess.Neg)
	if a.Depth() != 2 {
		t.Log("Depth updates improperly\n")
		t.Fail()
	}
	if p := a.Get(3); p != guess.Neg {
		t.Logf("PushAssigned 3, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 3, p)
		t.Fail()
	}
	if p := a.Get(2); p != guess.Pos {
		t.Logf("PushAssigned 2, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 2, p)
		t.Fail()
	}

	a.PopAssign()
	if a.Depth() != 1 {
		t.Log("Depth PopAssigns improperly\n")
		t.Fail()
	}
	if p := a.Get(3); p != guess.Unassigned {
		t.Logf("PopAssigned value is %s\n", p)
		t.Fail()
	}
	if p := a.Get(2); p != guess.Pos {
		t.Logf("PushAssigned 2, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 2, p)
		t.Fail()
	}

	a.PopAssign()
	if a.Depth() != 0 {
		t.Log("Depth PopAssigns improperly\n")
		t.Fail()
	}
	if p := a.Get(2); p != guess.Unassigned {
		t.Logf("PopAssigned value is %s\n", p)
		t.Fail()
	}
}
