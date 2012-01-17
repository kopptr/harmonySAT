package assignment

import (
	"log"
   "cnf"
)

const (
   Unassigned byte = 0
   Pos byte = cnf.Pos
   Neg byte = cnf.Neg
}

/* A stack of boolean assignments.
 * Each item on the stack represents a snapshot of the assignment the programmer
 * may want to return to. One can PushAssign, which makes a snapshot and assigns
 * the given literal to the next item, or PopAssign and return to the previous
 * state.
 */
type Assignment struct {
	// top of the stack of assignmentNodes
	top *assignmentNode
	// Number of nodes on the stack, not counting the empty assignment.
	// Also |PushAssign| - |PopAssign|
	depth uint
}

// A node in the assignment stack.
type assignmentNode struct {
	prev     *assignmentNode
	vars     []byte
	assigned int
}

// Allocates and returns the empty assignment */
func NewAssignment(nVars int) (a *Assignment) {
	a = &Assignment{nil, 0}
	a.top = &assignmentNode{nil, nil, 0}
	a.top.vars = make([]byte, nVars)
	return
}

// Returns the number of items pushed onto the stack
func (a *Assignment) Depth() uint {
	return a.depth
}

// Assigns a literal to the top assignment in the stack.
func (a *Assignment) Assign(l Lit) bool {
	if a.top == nil {
		log.Print("Assign() called on uninitialized Assignment")
		return false
	}
	if l.Val < 1 || l.Val > uint(len(a.top.vars)) {
		log.Printf("Attempted to assign %d (#vars=%d)", l.Val, len(a.top.vars))
		return false
	}
	if a.top.vars[l.Val-1] == Unassigned {
		a.top.assigned++
	}
	a.top.vars[l.Val-1] = l.Pol
	return true
}

func (a *Assignment) Get(i uint) (byte, bool) {
	if a.top == nil {
		log.Print("Get() called on uninitialized Assignment")
		return Unassigned, false
	}
	if i < 1 || i > uint(len(a.top.vars)) {
		log.Printf("Attempted to get %d (#vars=%d)", i, len(a.top.vars))
		return Unassigned, false
	}
	return a.top.vars[i-1], true
}

func (a *Assignment) PushAssign(l Lit) bool {
	if l.Pol == Unassigned {
		log.Print("Attempted to PushAssign an unassigned literal\n")
		return false
	}
	if r, e := a.Get(l.Val); r.Pol != Unassigned && e {
		log.Printf("Attempted to PushAssign a previously assigned %s\n", l)
		return false
	}

	newNode := &assignmentNode{nil, nil, 0}
	newNode.prev = a.top

	newNode.assigned = a.top.assigned
	newNode.vars = make([]byte, len(a.top.vars))
	copy(newNode.vars, a.top.vars[:])

	a.top = newNode
	a.depth++
	return a.Assign(l)
}

func (a *Assignment) PopAssign() bool {
	if a.top.prev == nil {
		log.Printf("Attempted to PopAssign the empty assignment\n")
		return false
	}
	a.top = a.top.prev
	a.depth--

	return true
}

func (a *Assignment) Initialized() bool {
	if a.top == nil {
		return false
	}
	return true
}
