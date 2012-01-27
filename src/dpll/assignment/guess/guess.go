package guess

import (
	"bytes"
	"errors"
	"fmt"
)

// TODO make this a real type
const (
	Unassigned byte = 0
	Pos        byte = 1
	Neg        byte = 2
)

// Guess is a struct because we want to create a thin abstraction over the fact
// that it's just a simple array. This way, the user doesn't need to worry about
// subtracting 1. Silly, but I'm okay with it.
type Guess struct {
	vars []byte
}

func NewGuess(nVars int) (g *Guess) {
	g = new(Guess)
	g.vars = make([]byte, nVars)
	return
}

func (g *Guess) Copy() (g1 *Guess) {
	g1 = NewGuess(g.Len())
	copy(g1.vars, g.vars[:])
	return
}

// Sets the variable n. v must be \in {Unassigned,Pos,Neg}
func (g *Guess) Set(n uint, v byte) error {
	if n < 1 || n > uint(g.Len()) || (v != Pos && v != Neg && v != Unassigned) {
		return errors.New(
			fmt.Sprintf("Guess.Set() given invalid index %d or assignment %d", n, v))
	}
	g.vars[n-1] = v
	return nil
}

// Returns the number of variables that can be assigned
func (g *Guess) Len() int {
	return len(g.vars)
}

// Returns what is assigned to the nth variable.
func (g *Guess) Get(n uint) (byte, error) {
	if n < 1 || n > uint(g.Len()) {
		return Unassigned, errors.New("Guess.Get() index out of bounds")
	}
	return g.vars[n-1], nil
}

// Returns an array of ints representing the assigned variables.
func (g *Guess) Vars(flipped bool) (v []int) {
	v = []int{}
	for i, n := range g.vars {
		if n == Pos {
			if !flipped {
				v = append(v, i+1)
			} else {
				v = append(v, -1*(i+1))
			}
		} else if n == Neg {
			if !flipped {
				v = append(v, (i+1)*-1)
			} else {
				v = append(v, i+1)
			}
		}
	}
	return
}

// String representation
func (g Guess) String() string {
	buffer := bytes.NewBufferString("")
	for i, l := range g.vars {
		if l == Pos {
			fmt.Fprintf(buffer, "%d ", i+1)
		} else if l == Neg {
			fmt.Fprintf(buffer, "%d ", -1*(i+1))
		}
	}
	return string(buffer.Bytes())
}
