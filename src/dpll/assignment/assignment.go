package assignment

import (
	"dpll/assignment/guess"
	"errors"
)

type Assignment struct {
	top   *assignmentNode
	depth uint
}

// A node in the assignment stack.
type assignmentNode struct {
	prev     *assignmentNode
	g        *guess.Guess
}

// Allocates and returns the empty assignment
func NewAssignment(nVars int) (a *Assignment) {
	a = new(Assignment)
	a.top = new(assignmentNode)
	a.top.g = guess.NewGuess(nVars)
	return
}

// Returns the number of items pushed onto the stack
func (a *Assignment) Depth() uint {
	return a.depth
}

func (a *Assignment) Get(i uint) (byte, error) {
	return a.top.g.Get(i)
}

func (a *Assignment) Len() uint {
	return uint(a.top.g.Len())
}

func (a *Assignment) PushAssign(v uint, pol byte) error {

	if check, err := a.top.g.Get(v); check != guess.Unassigned || err != nil {
		return err
	}
	newNode := &assignmentNode{nil, nil}
	newNode.prev = a.top

	newNode.g = a.top.g.Copy()

	a.top = newNode
	a.depth++
	a.top.g.Set(v, pol)
	return nil
}

func (a *Assignment) Guess() *guess.Guess {
	return a.top.g
}

func (a *Assignment) PopAssign() error {
	if a.top == nil {
		return errors.New("Tried to Assignment.PopAssign() empty assignment")
	}
	a.top = a.top.prev
	a.depth--
	return nil
}
