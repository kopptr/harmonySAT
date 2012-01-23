package db

import (
   "testing"
)

func TestNewEntry(t *testing.T) {
   e := NewEntry([]int{-5,3,4})
   if e.Next != nil || e.Prev != nil {
      t.Logf("e.Next/Prev is messed up\n")
      t.Fail()
   }
}


