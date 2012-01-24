package dpll

import (
   "testing"
   "dpll/assignment"
   "dpll/assignment/guess"
   "dpll/db"
   "dpll/db/cnf"
)

func TestDecide(t *testing.T) {
   a := assignment.NewAssignment(10)
   l := decide(nil, a)
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
}

func TestDpll(t *testing.T) {
   db := db.NewDB(10)
   db.AddEntry([]int{1,2,3})
   db.AddEntry([]int{4,5,6})
   db.AddEntry([]int{-1,-2,3})
   db.AddEntry([]int{3,2,7})



}
