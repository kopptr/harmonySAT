package cnf

import (
	"strconv"
)

// Represents a propositional literal. Val is a unique id for the
// variable, Pol is it's polarity, Pos(itive) or Neg(ative)
type Lit struct {
	Val uint
	Pol byte
}

// Possible values for Lit.Pol
const (
	Pos byte = 1
	Neg byte = 2
)

// Returns a Lit given the int representation
func NewLit(x int) (l Lit) {
	if x > 0 {
		l.Val = uint(x)
		l.Pol = Pos
	} else if x < 0 {
		l.Val = uint(x * -1)
		l.Pol = Neg
	} // else Lit{0,0} is implicitly set
	return
}

// Returns the int representation of a Lit
func (l *Lit) Int() (x int) {
	if l.Pol == Pos {
		x = int(l.Val)
	} else {
		x = -1 * int(l.Val)
	}
	return
}

// Returns the string representation of a Lit
func (l Lit) String() (s string) {
	if l.Pol == Pos {
		s = strconv.FormatUint(uint64(l.Val), 10)
	} else if l.Pol == Neg {
		s = "-" + strconv.FormatUint(uint64(l.Val), 10)
	}
	return
}

func (l *Lit) Flip() {
	if l.Pol == Pos {
		l.Pol = Neg
	} else if l.Pol == Neg {
		l.Pol = Pos
	}
}

func (l *Lit) Eq(l1 *Lit) bool {
	return l.Val == l1.Val && l.Pol == l1.Pol
}
