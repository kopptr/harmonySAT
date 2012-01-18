package db

import (
   "testing"
)


func TestNewWatch(t *testing.T) {
   w := NewWatch()
   if w.Next != w {
      t.Logf("Watch.Next is broken\n")
      t.Fail()
   }
   if w.Prev != w {
      t.Logf("Watch.Prev is broken\n")
      t.Fail()
   }
   if w.Watching.Val != 0 {
      t.Logf("Watch.Watching.Val is broken\n")
      t.Fail()
   }
   if w.Watching.Pol != 0 {
      t.Logf("Watch.Watching.Pol is broken\n")
      t.Fail()
   }
}
