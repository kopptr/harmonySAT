package db

import (
   "guess"
   "cnf"
)

func (db *DB) Bcp(g *guess.Guess, lit cnf.Lit) {

   lq := newLitQ() // queue of literals to be assigned

   lq.PushBack(lit)

   // For each new literal in the queue
   for l,ok := lq.PopFront(); ok; l,ok = lq.PopFront() {
      g.Set(l.Val, l.Pol)
      // For each clause watching that literal
//      wl := db.GetWatchList(l)
//      for w := wl.First(); w != nil; w = wl.Next() {

 //     }
   }
}

type litQ struct {
   First *litQNode
   Last *litQNode
}

type litQNode struct {
   Next *litQNode
   L cnf.Lit
}

func newLitQNode(l cnf.Lit) (lqn *litQNode) {
   lqn = new(litQNode)
   lqn.L.Val = l.Val
   lqn.L.Pol = l.Pol
   lqn.Next = nil
   return
}

func newLitQ() (lq *litQ) {
   lq.First = nil
   lq.Last = nil
   return
}

func (lq litQ) PushBack(l cnf.Lit) {
   lqn := newLitQNode(l)
   if lq.First == nil { // empty
      lq.First = lqn
      lq.Last = lqn
   } else {
      lq.Last.Next = lqn
      lq.Last = lqn
   }
}

func (lq litQ) PopFront() (l cnf.Lit, ok bool) {
   if lq.First == nil { // empty
      l = cnf.Lit{0,0}
      ok = false
   } else {
      l = lq.First.L
      ok = true
      lq.First = lq.First.Next
   }
   return
}
