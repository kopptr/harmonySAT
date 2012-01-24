package assignment

import (
   "dpll/assignment/guess"
)

type Assignment struct {
	top *assignmentNode
	depth uint
}

// A node in the assignment stack.
type assignmentNode struct {
	prev     *assignmentNode
	g        *guess.Guess
	assigned int
}

// Allocates and returns the empty assignment */
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

func (a *Assignment) Get(i uint) byte {
	return a.top.g.Get(i)
}

func (a *Assignment) Len() uint {
   return uint(a.top.g.Len())
}

func (a *Assignment) PushAssign(v uint, pol byte) {
	newNode := &assignmentNode{nil, nil, 0}
	newNode.prev = a.top

	newNode.assigned = a.top.assigned
	newNode.g = a.top.g.Copy()

	a.top = newNode
	a.depth++
	a.top.g.Set(v, pol)
}

func (a *Assignment) Guess() *guess.Guess {
   return a.top.g
}

func (a *Assignment) PopAssign() {
	a.top = a.top.prev
	a.depth--
}
