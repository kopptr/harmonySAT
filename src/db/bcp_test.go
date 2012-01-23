package db

import (
   "testing"
   "guess"
   "cnf"
//   "fmt"
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

func TestBcpUnits(t *testing.T) {
   db := NewDB(10)
   db.AddEntry([]int{1,2,3,4,5})
   db.AddEntry([]int{-1,2})
   db.AddEntry([]int{-2,-3})
   db.AddEntry([]int{3,4})
   db.AddEntry([]int{-4,5})
   db.StartLearning()

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
   if g.Get(3) != guess.Neg {
      t.Logf("Bcp does not assign -3 as it should\n")
      t.Fail()
   }
   if g.Get(4) != guess.Pos {
      t.Logf("Bcp does not assign 4 as it should\n")
      t.Fail()
   }
   if g.Get(5) != guess.Pos {
      t.Logf("Bcp does not assign 4 as it should\n")
      t.Fail()
   }
}


func TestBcpConflict(t *testing.T) {
   db := NewDB(10)
   db.AddEntry([]int{1,2,3,4,5})
   db.AddEntry([]int{-1,2})
   db.AddEntry([]int{-2,-3})
   db.AddEntry([]int{3,4})
   db.AddEntry([]int{-4,3})
   db.StartLearning()

   g := guess.NewGuess(10)

   result := db.Bcp(g, cnf.Lit{1, guess.Pos})
   if result {
      t.Logf("Bcp returned true\n")
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
   if g.Get(3) != guess.Neg {
      t.Logf("Bcp does not assign -3 as it should\n")
      t.Fail()
   }
   if g.Get(4) != guess.Pos {
      t.Logf("Bcp does not assign 4 as it should\n")
      t.Fail()
   }
   if db.NGiven() != 5 {
      t.Logf("NGiven is after bcp w/ conflict incorrect\n")
      t.Fail()
   }
   if db.NLearned() != 1 {
      t.Logf("NLearned is after bcp w/ conflict incorrect\n")
      t.Fail()
   }
   //fmt.Printf("%s\n", db.Learned)

}

func TestBcpBoring(t *testing.T) {
   db := NewDB(10)
   db.AddEntry([]int{1,2,3,4,5})
   db.AddEntry([]int{1,2})
   db.AddEntry([]int{-2,-3})
   db.AddEntry([]int{3,4})
   db.AddEntry([]int{-4,3})
   db.StartLearning()

   g := guess.NewGuess(10)

   result := db.Bcp(g, cnf.Lit{1, guess.Pos})
   if !result {
      t.Logf("Bcp returned false\n")
      t.Fail()
   }
   if g.Get(1) != guess.Pos {
      t.Logf("Bcp does not assign 1 as it should\n")
      t.Fail()
   }
   if g.Get(2) != guess.Unassigned {
      t.Logf("Bcp does not assign 2 as it should\n")
      t.Fail()
   }
   if g.Get(3) != guess.Unassigned {
      t.Logf("Bcp does not assign -3 as it should\n")
      t.Fail()
   }
   if g.Get(4) != guess.Unassigned {
      t.Logf("Bcp does not assign 4 as it should\n")
      t.Fail()
   }
   if g.Get(5) != guess.Unassigned {
      t.Logf("Bcp does not assign 4 as it should\n")
      t.Fail()
   }
}
