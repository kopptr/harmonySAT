package sat

import (
   "testing"
)

func TestComposite(t *testing.T){
   var l1 Lit = Lit{2,Pos}
   if l1.Val != 2 || l1.Pol != Pos {
      t.Logf("l1 expected: 2\n")
      t.Logf("l1 actual:   %s\n", l1)
      t.Fail()
   }
   var l2 Lit = Lit{500,Neg}
   if l2.Val != 500 || l2.Pol != Neg {
      t.Logf("l1 expected: -500\n")
      t.Logf("l1 actual:   %s\n", l2)
      t.Fail()
   }
}

func TestString(t *testing.T){
   var l1 Lit = Lit{2,Pos}
   if l1.String()!= "2" {
      t.Logf("l1 expected: 2\n")
      t.Logf("l1 actual:   %s\n", l1)
      t.Fail()
   }
   var l2 Lit = Lit{500,Neg}
   if l2.String() != "-500" {
      t.Logf("l2 expected: -500\n")
      t.Logf("l2 actual:   %s\n", l2)
      t.Fail()
   }
   var l3 Lit
   if l3.String() != "" {
      t.Logf("l3 expected: \n")
      t.Logf("l3 actual:   %s\n", l3)
      t.Fail()
   }

}

func TestIsSet(t *testing.T) {
   var l0 Lit
   var l1 Lit = Lit{6,Unassigned}
   var l2 Lit = Lit{6,Neg}
   var l3 Lit = Lit{6,Pos}

   if l0.IsSet() {
      t.Log("l0.IsSet() returns true for uninitialized Lit\n")
      t.Fail()
   }
   if l1.IsSet() {
      t.Log("l1.IsSet() returns true for unassigned Lit\n")
      t.Fail()
   }
   if ! l2.IsSet() {
      t.Log("l2.IsSet() returns false for unassigned Lit\n")
      t.Fail()
   }
   if ! l3.IsSet() {
      t.Log("l3.IsSet() returns false for unassigned Lit\n")
      t.Fail()
   }
}

func TestFlip(t *testing.T) {
   var l0 Lit
   var l1 Lit = Lit{6,Unassigned}
   var l2 Lit = Lit{6,Neg}
   var l3 Lit = Lit{6,Pos}
   l0.Flip()
   if l0.Pol != Unassigned {
      t.Logf("l0 expected: \n")
      t.Logf("l0 actual:   %s\n", l0)
      t.Fail()
   }
   l1.Flip()
   if l1.Pol != Unassigned {
      t.Logf("l1 expected: \n")
      t.Logf("l1 actual:   %s\n", l1)
      t.Fail()
   }
   l2.Flip()
   if l2.Pol != Pos {
      t.Logf("l2 expected: 6\n")
      t.Logf("l2 actual:   %s\n", l2)
      t.Fail()
   }
   l3.Flip()
   if l3.Pol != Neg {
      t.Logf("l3 expected: -6\n")
      t.Logf("l3 actual:   %s\n", l3)
      t.Fail()
   }
}



