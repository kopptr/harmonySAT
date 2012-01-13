package sat

import (
   "log"
)

type Assignment struct {
   // top of the stack of assignmentNodes
   top *assignmentNode
   // Number of nodes on the stack, not counting the empty assignment.
   // Also |PushAssign| - |PopAssign|
   depth uint
}

type assignmentNode struct {
   prev *assignmentNode
   vars []byte
   assigned int
}

func NewAssignment(nVars int) (a *Assignment) {
   a = &Assignment{nil,0}
   a.top = &assignmentNode{nil,nil,0}
   a.top.vars = make([]byte, nVars)
   return
}

func (a *Assignment) Depth() uint {
   return a.depth
}

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

func (a *Assignment) Get(i uint) (Lit, bool) {
   if a.top == nil {
      log.Print("Get() called on uninitialized Assignment")
      return Lit{}, false
   }
   if i < 1 || i > uint(len(a.top.vars)) {
      log.Printf("Attempted to get %d (#vars=%d)", i, len(a.top.vars))
      return Lit{}, false
   }
   return Lit{i, a.top.vars[i-1]}, true
}

func (a *Assignment) PushAssign(l Lit) bool {
   if l.Pol == Unassigned {
      log.Print("Attempted to PushAssign an unassigned literal\n")
      return false
   }
   if r,e := a.Get(l.Val); r.Pol != Unassigned && e {
      log.Printf("Attempted to PushAssign a previously assigned %s\n", l)
      return false
   }

   newNode := &assignmentNode{nil,nil,0}
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
