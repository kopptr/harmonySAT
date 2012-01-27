package dpll

import (
   "testing"
   "dpll/assignment"
   "dpll/assignment/guess"
   "dpll/db"
   "dpll/db/cnf"
   "fmt"
)

func TestDecide(t *testing.T) {
   fmt.Printf("foz\n")
   a := assignment.NewAssignment(10)
   l := decide(nil, a)
   fmt.Printf("faz\n")
   if ! l.Eq(&cnf.Lit{1,cnf.Pos}) {
      t.Logf("Decide didn't pick 1, it picked \n", l)
      t.Fail()
   }
   a.PushAssign(2, guess.Pos)
   a.PushAssign(l.Val, l.Pol)
   l = decide(nil, a)
   if !l.Eq(&cnf.Lit{3,cnf.Pos}) {
      t.Logf("Decide didn't pick 3, it picked \n", l)
      t.Fail()
   }
   fmt.Printf("foo\n")
}

func TestDpll(t *testing.T) {
   db := db.NewDB(10)
   db.AddEntry([]int{1,2,3,-9})
   db.AddEntry([]int{4,5,6,-9})
   db.AddEntry([]int{-1,-2,3,-9})
   db.AddEntry([]int{3,2,7,-9})
   a := assignment.NewAssignment(10)

   g := Dpll(db, a)
   if g == nil {
      t.Logf("Dpll returned nil\n")
      t.Fail()
   }
   if !db.Verify(g) {
      t.Logf("Dpll returned incorrect solution\n")
      t.Fail()
   }
   fmt.Printf("bar\n")
}

func TestDpll2(t *testing.T) {
   db := db.NewDB(5)
   db.AddEntry([]int{4, -2, -5})
   db.AddEntry([]int{3, 1, -5})
   db.AddEntry([]int{-5, -4, 3})
   db.AddEntry([]int{-3, 4, -2})
   db.AddEntry([]int{1, 4, -3})
   db.AddEntry([]int{-2, 1, 5})
   db.AddEntry([]int{-1, 5, 3})
   db.AddEntry([]int{-4, 2, 5})
   db.AddEntry([]int{-2, -3, -4})
   db.AddEntry([]int{4, -5, 3})
   a := assignment.NewAssignment(5)

   g := Dpll(db, a)
   if g == nil {
      t.Logf("Dpll returned nil\n")
      t.Fail()
   }
   if !db.Verify(g) {
      t.Logf("Dpll returned incorrect solution\n")
      t.Fail()
   }
   fmt.Printf("baz\n")
}

