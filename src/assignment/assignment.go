package assignment

import (
	"log"
)

const (
   Unassigned byte = 0
   Pos byte = 1
   Neg byte = 2
)

/* A stack of boolean assignments.
 * Each item on the stack represents a snapshot of the assignment the programmer
 * may want to return to. One can PushAssign, which makes a snapshot and assigns
 * the given literal to the next item, or PopAssign and return to the previous
 * state.
 */
type Assignment struct {
	top *assignmentNode // top of the stack of assignmentNodes
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
func (a *Assignment) Assign(v uint, pol byte) bool {
	if a.top == nil {
		log.Print("Assign() called on uninitialized Assignment")
		return false
	}
	if v < 1 || v > uint(len(a.top.vars)) {
		log.Printf("Attempted to assign %d (#vars=%d)", v, len(a.top.vars))
		return false
	}
	if a.top.vars[v-1] == Unassigned {
		a.top.assigned++
	}
	a.top.vars[v-1] = pol
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

func (a *Assignment) PushAssign(v uint, pol byte) bool {
	if pol == Unassigned {
		log.Print("Attempted to PushAssign an unassigned literal\n")
		return false
	}
	if p, e := a.Get(v); p != Unassigned && e {
		log.Printf("Attempted to PushAssign a previously assigned %s\n", v)
		return false
	}

	newNode := &assignmentNode{nil, nil, 0}
	newNode.prev = a.top

	newNode.assigned = a.top.assigned
	newNode.vars = make([]byte, len(a.top.vars))
	copy(newNode.vars, a.top.vars[:])

	a.top = newNode
	a.depth++
	return a.Assign(v, pol)
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
