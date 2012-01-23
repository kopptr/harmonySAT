package guess

import (
   "testing"
)

func TestGuess(t *testing.T) {
   g := NewGuess(10)

   g.Set(4,Neg)

   if g.Get(4) != Neg {
      t.Logf("It are broken\n")
      t.Fail()
   }
}
