package cnf

import (
	"testing"
)

func TestComposite(t *testing.T) {
	var l1 Lit = Lit{2, Pos}
	if l1.Val != 2 || l1.Pol != Pos {
		t.Logf("l1 expected: 2\n")
		t.Logf("l1 actual:   %s\n", l1)
		t.Fail()
	}
	var l2 Lit = Lit{500, Neg}
	if l2.Val != 500 || l2.Pol != Neg {
		t.Logf("l2 expected: -500\n")
		t.Logf("l2 actual:   %s\n", l2)
		t.Fail()
	}
}

func TestNewLit(t *testing.T) {
	var l1 Lit = NewLit(2)
	if l1.Val != 2 || l1.Pol != Pos {
		t.Logf("l1 expected: 2\n")
		t.Logf("l1 actual:   %s\n", l1)
		t.Fail()
	}
	var l2 Lit = NewLit(-500)
	if l2.Val != 500 || l2.Pol != Neg {
		t.Logf("l2 expected: -500\n")
		t.Logf("l2 actual:   %s\n", l2)
		t.Fail()
	}
	var l3 Lit = NewLit(0)
	if l3.Val != 0 || l3.Pol != 0 {
		t.Logf("l3 expected: ")
		t.Logf("l3 actual:   %s\n", l3)
		t.Fail()
	}
}

func TestString(t *testing.T) {
	var l1 Lit = Lit{2, Pos}
	if l1.String() != "2" {
		t.Logf("l1 expected: 2\n")
		t.Logf("l1 actual:   %s\n", l1)
		t.Fail()
	}
	var l2 Lit = Lit{500, Neg}
	if l2.String() != "-500" {
		t.Logf("l2 expected: -500\n")
		t.Logf("l2 actual:   %s\n", l2)
		t.Fail()
	}
	var l3 Lit
	if l3.String() != "" {
		t.Logf("l3 expected: \n")
		t.Logf("l3 actual:   %s\n", l3)
		t.Fail()
	}

}

func TestInt(t *testing.T) {
	var l1 Lit = Lit{2, Pos}
	if l1.Int() != 2 {
		t.Logf("l1 expected: 2\n")
		t.Logf("l1 actual:   %s\n", l1)
		t.Fail()
	}
	var l2 Lit = Lit{500, Neg}
	if l2.Int() != -500 {
		t.Logf("l2 expected: -500\n")
		t.Logf("l2 actual:   %s\n", l2)
		t.Fail()
	}
	var l3 Lit
	if l3.Int() != 0 {
		t.Logf("l3 expected: \n")
		t.Logf("l3 actual:   %s\n", l3)
		t.Fail()
	}

}
