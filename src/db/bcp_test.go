package db

import (
   "testing"
   "guess"
   "cnf"
)

func TestLitQ(t *testing.T) {
   lq := newLitQ()
   lq.PushBack(cnf.Lit{1,cnf.Pos})

   l, ok := lq.PopFront();
   if !ok {
      t.Logf("PopFront ok is Broken\n")
      t.Fail()
   }
   if l.Val != 1 || l.Pol != cnf.Pos {
      t.Logf("PopFront Lit is %s\n", l)
      t.Fail()
   }
}

func TestBcp(t *testing.T) {
   db := NewDB(10)
   db.AddEntry([]int{1,2,3,4,5})
   db.AddEntry([]int{-1,2})

   g := guess.NewGuess(10)

   // BCP on a 1. It should assign 2.
   result := db.Bcp(g, cnf.Lit{1, guess.Pos})
   if !result {
      t.Logf("Bcp returned false\n")
      t.Fail()
   }
   if g.Get(1) != guess.Pos {
      t.Logf("Bcp does not assign 1 as it should\n")
      t.Fail()
   }
   if g.Get(2) != guess.Pos {
      t.Logf("Bcp does not assign 2 as it should\n")
      t.Fail()
   }
}
