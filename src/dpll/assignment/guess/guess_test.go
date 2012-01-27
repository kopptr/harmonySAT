package guess

import (
   "testing"
)

func TestGuess(t *testing.T) {
   g := NewGuess(10)

   if err := g.Set(4,Neg); err != nil {
      t.Logf("Setting -4 failed: %s\n", err.Error())
      t.FailNow()
   }

   if p, err := g.Get(4); p != Neg || err != nil {
      if err != nil {
         t.Logf("error: %s\n", err.Error())
      }
      t.Logf("Get returns %d, should return %d\n", p, Neg)
      t.Fail()
   }
}
