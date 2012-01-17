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
		s = strconv.Uitoa(l.Val)
	} else if l.Pol == Neg {
		s = "-" + strconv.Uitoa(l.Val)
	}
	return
}
