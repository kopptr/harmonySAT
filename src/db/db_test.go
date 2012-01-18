package db

import (
   "testing"
   "fmt"
)

func TestNewDB(t *testing.T) {
   db := NewDB(10)
   if db.Binary != 0 || db.Ternary != 0 || db.Horn != 0 || db.Definite != 0 {
      t.Logf("db.Stats not initialized properly\n")
      t.Fail()
   }
   if db.NLearned() != 0 || db.NGiven() != 0 {
      t.Logf("db.NLearned/Given() not initialized properly\n")
      t.Fail()
   }
   if db.Learned != nil || db.Given != nil {
      t.Logf("db.Learned/Given not initialized properly\n")
      t.Fail()
   }
   for i := range db.WatchRings {
      if db.WatchRings[i].Watching.Pol != 0 || db.WatchRings[i].Watching.Val != 0 {
         t.Logf("watch.Pol/Val initialized incorrectly\n")
         t.Fail()
      }
      if !db.WatchRings[i].IsDummy() {
         t.Logf("watch.IsDummy incorrect\n")
         t.Fail()
      }
      if db.WatchRings[i].Next != db.WatchRings[i] {
         t.Logf("w.Next not initialized to w\n")
         t.Fail()
      }
      if db.WatchRings[i].Prev != db.WatchRings[i] {
         t.Logf("w.Prev not initialized to w\n")
         t.Fail()
      }
   }
}

func TestAddEntry(t *testing.T) {
   db := NewDB(10)
   db.AddEntry([]int{-1,3,-5})
   if db.NLearned() != 0 || db.NGiven() != 1 {
      t.Logf("entry totals update incorrectly\n")
      t.Fail()
   }
   fmt.Printf("%s\n", db.String() )
   db.AddEntry([]int{-10,4,-5,6,1})
   fmt.Printf("%s\n", db.String() )
   db.AddEntry([]int{-9,4,-5,6,1})
   fmt.Printf("%s\n", db.String() )

   db.StartLearning()
   db.AddEntry([]int{-9,-4,-5,6,1})
   fmt.Printf("%s\n", db.String() )
   db.AddEntry([]int{-9,4,-3,6,1})
   fmt.Printf("%s\n", db.String() )
   db.AddEntry([]int{-9,4,-5,2,1})
   fmt.Printf("%s\n", db.String() )
}
