package cnf

import (
	"testing"
)

func TestNewClause(t *testing.T) {
	is1 := []int{-3, 2, -5, 1}
	c1 := NewClause(is1)
	if !sameStats(c1.ClauseStats, ClauseStats{false, false, false, false}) {
		t.Log("c1's stats are incorrect\n")
		t.Fail()
	}

	is2 := []int{-3, -5, 1}
	c2 := NewClause(is2)
	if !sameStats(c2.ClauseStats, ClauseStats{false, true, true, true}) {
		t.Log("c2's stats are incorrect\n")
		t.Fail()
	}

	is3 := []int{-3, -86, -5, 1}
	c3 := NewClause(is3)
	if !sameStats(c3.ClauseStats, ClauseStats{false, false, true, true}) {
		t.Log("c3's stats are incorrect\n")
		t.Fail()
	}

	is4 := []int{-3, -5}
	c4 := NewClause(is4)
	if !sameStats(c4.ClauseStats, ClauseStats{true, false, true, false}) {
		t.Log("c4's stats are incorrect\n")
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	c1 := NewClause([]int{-3, 2, -5, 1, 10})

	in1, pol1 := c1.Contains(2)
	if !in1 {
		t.Log("c1 reports that it does not contain 2\n")
		t.Fail()
	}
	if pol1 != Pos {
		t.Log("c1 reports having 2 in the wrong polarity\n")
	}
	in2, pol2 := c1.Contains(3)
	if !in2 {
		t.Log("c1 reports that it does not contain 3\n")
		t.Fail()
	}
	if pol2 != Neg {
		t.Log("c1 reports having 3 in the wrong polarity\n")
	}
	in3, pol3 := c1.Contains(4)
	if in3 {
		t.Log("c1 reports that it contains 4\n")
		t.Fail()
	}
	if pol3 != 0 {
		t.Log("c1 reports a polarity for 4\n")
	}
}

func sameStats(a, b ClauseStats) bool {
	if a.binary != b.binary {
		return false
	}
	if a.ternary != b.ternary {
		return false
	}
	if a.horn != b.horn {
		return false
	}
	if a.definite != b.definite {
		return false
	}
	return true
}
