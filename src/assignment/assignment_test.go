package assignment

import (
	"testing"
)

func TestNew(t *testing.T) {
	var (
		a Assignment
		b *Assignment
	)
	if a.Initialized() {
		t.Logf("Assignment reports it is initialized before call to New()\n")
		t.Fail()
	}
	b = NewAssignment(10)
	if !b.Initialized() {
		t.Logf("Assignment reports it is uninitialized after call to New()\n")
		t.Fail()
	}
	for i := 0; i < 10; i++ {
		if b.top.vars[i] != Unassigned {
			t.Logf("Assignment reports assigned variables after initialization\n")
			t.Fail()
		}
	}
}

func TestAssign(t *testing.T) {
	b := NewAssignment(10)

	if p, e := b.Get(2); p != Unassigned || !e {
		t.Logf("Assign nothing, get 2, it returns %s\n", p)
		t.Fail()
	}

	if e := b.Assign(2, Pos); !e {
		t.Log("Assigning 2 failed\n")
		t.Fail()
	}
	if p, e := b.Get(2); p != Pos || !e {
		t.Logf("Assign 2, get, it returns %s\n", p)
		t.Fail()
	}

	if e := b.Assign(2, Neg); !e {
		t.Log("Assigning -2 failed\n")
		t.Fail()
	}
	if p, e := b.Get(2); p != Neg || !e {
		t.Logf("Assign -2, get, it returns %s\n", p)
		t.Fail()
	}

	if e := b.Assign(2, Unassigned); !e {
		t.Log("Assigning <2> failed\n")
		t.Fail()
	}
	if p, e := b.Get(2); p != Unassigned || !e {
		t.Logf("Assign <2>, get, it returns %s\n", p)
		t.Fail()
	}

	if e := b.Assign(11, Pos); e {
		t.Log("Assigning 11 succeeded\n")
		t.Fail()
	}
	if _, e := b.Get(11); e {
		t.Log("Getting 11 succeeded\n")
		t.Fail()
	}
	if e := b.Assign(0, Pos); e {
		t.Log("Assigning 0 succeeded\n")
		t.Fail()
	}
	if _, e := b.Get(0); e {
		t.Log("Getting 0 succeeded\n")
		t.Fail()
	}
}

func TestPushPopAssign(t *testing.T) {
	a := NewAssignment(10)

	if a.Depth() != 0 {
		t.Log("Depth initializes improperly\n")
		t.Fail()
	}
	if p, e := a.Get(2); p != Unassigned || !e {
		t.Logf("Initial value is %s\n", p)
		t.Fail()
	}

	a.PushAssign(2, Pos)
	if a.Depth() != 1 {
		t.Log("Depth updates improperly\n")
		t.Fail()
	}
	if p, e := a.Get(2); p != Pos || !e {
		t.Logf("PushAssigned 2, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 2, p)
		t.Fail()
	}
	a.PushAssign(3, Neg)
	if a.Depth() != 2 {
		t.Log("Depth updates improperly\n")
		t.Fail()
	}
	if p, e := a.Get(3); p != Neg || !e {
		t.Logf("PushAssigned 3, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 3, p)
		t.Fail()
	}
	if p, e := a.Get(2); p != Pos || !e {
		t.Logf("PushAssigned 2, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 2, p)
		t.Fail()
	}

	a.PopAssign()
	if a.Depth() != 1 {
		t.Log("Depth PopAssigns improperly\n")
		t.Fail()
	}
	if p, e := a.Get(3); p != Unassigned || !e {
		t.Logf("PopAssigned value is %s\n", p)
		t.Fail()
	}
	if p, e := a.Get(2); p != Pos || !e {
		t.Logf("PushAssigned 2, value is %s\n", p)
		t.Logf("val: %d, pol: %d\n", 2, p)
		t.Fail()
	}

	a.PopAssign()
	if a.Depth() != 0 {
		t.Log("Depth PopAssigns improperly\n")
		t.Fail()
	}
	if p, e := a.Get(2); p != Unassigned || !e {
		t.Logf("PopAssigned value is %s\n", p)
		t.Fail()
	}
}