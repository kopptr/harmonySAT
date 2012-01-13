package sat

import (
   "testing"
)

func TestNew(t *testing.T) {
   var (
      a Assignment
      b *Assignment
   )
   if a.Initialized() {
      t.Logf("Assignment reports it is initialized before call to New()\n")
      t.Fail()
   }
   b = NewAssignment(10)
   if !b.Initialized() {
      t.Logf("Assignment reports it is uninitialized after call to New()\n")
      t.Fail()
   }
   for i := 0; i < 10; i++ {
      if b.top.vars[i] != Unassigned {
         t.Logf("Assignment reports assigned variables after initialization\n")
         t.Fail()
      }
   }
}

func TestAssign(t *testing.T) {
   b := NewAssignment(10)

   if l,e := b.Get(2); l.Pol != Unassigned || l.Val != 2 || !e {
      t.Logf("Assign nothing, get 2, it returns %s\n", l)
      t.Fail()
   }

   if e := b.Assign(Lit{2,Pos}); !e {
      t.Log("Assigning 2 failed\n")
      t.Fail()
   }
   if l,e := b.Get(2); l.Pol != Pos || l.Val != 2 || !e {
      t.Logf("Assign 2, get, it returns %s\n", l)
      t.Fail()
   }

   if e := b.Assign(Lit{2,Neg}); !e {
      t.Log("Assigning -2 failed\n")
      t.Fail()
   }
   if l,e := b.Get(2); l.Pol != Neg || l.Val != 2 || !e {
      t.Logf("Assign -2, get, it returns %s\n", l)
      t.Fail()
   }

   if e := b.Assign(Lit{2,Unassigned}); !e {
      t.Log("Assigning <2> failed\n")
      t.Fail()
   }
   if l,e := b.Get(2); l.Pol != Unassigned || l.Val != 2 || !e {
      t.Logf("Assign <2>, get, it returns %s\n", l)
      t.Fail()
   }

   if e := b.Assign(Lit{11,Pos}); e {
      t.Log("Assigning 11 succeeded\n")
      t.Fail()
   }
   if _,e := b.Get(11); e {
      t.Log("Getting 11 succeeded\n")
      t.Fail()
   }
   if e := b.Assign(Lit{0,Pos}); e {
      t.Log("Assigning 0 succeeded\n")
      t.Fail()
   }
   if _,e := b.Get(0); e {
      t.Log("Getting 0 succeeded\n")
      t.Fail()
   }
}



