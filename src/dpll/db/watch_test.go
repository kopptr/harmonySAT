package db

import (
   "testing"
   //"fmt"
)

func TestNewWatch(t *testing.T) {
   w := NewWatch()
   if w.Next != nil {
      t.Logf("Watch.Next is broken\n")
      t.Fail()
   }
   if w.Prev != nil {
      t.Logf("Watch.Prev is broken\n")
      t.Fail()
   }
   if w.E != nil {
      t.Logf("Watch.E is broken\n")
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

func TestNewWatchListAdd(t *testing.T) {
   wr := NewWatchList()
   w1 := NewWatch()
   e1 := NewEntry([]int{1,2,3})
   w1.E = e1
   wr.Add(w1)
   if wr.first != w1 {
      t.Logf("Watch.first is wrong after first Add()\n")
      t.Fail()
   }
   if wr.first.Prev != nil {
      t.Logf("Watch.first.prev is non-nil after first Add()\n")
      t.Fail()
   }
   if wr.last != w1 {
      t.Logf("Watch.last is wrong after first Add()\n")
      t.Fail()
   }
   if wr.last.Next != nil {
      t.Logf("Watch.last.next is non-nil after first Add()\n")
      t.Fail()
   }
   if wr.current != w1 {
      t.Logf("Watch.current is wrong after first Add()\n")
      t.Fail()
   }
   //fmt.Printf("%s\n", wr)
   w2 := NewWatch()
   e2 := NewEntry([]int{2,3,4})
   w2.E = e2
   wr.Add(w2)
   if wr.first != w1 {
      t.Logf("Watch.first is wrong after second Add()\n")
      t.Fail()
   }
   if wr.first.Prev != nil {
      t.Logf("Watch.first.prev is non-nil after first Add()\n")
      t.Fail()
   }
   if wr.last != w2 {
      t.Logf("Watch.last is wrong after second Add()\n")
      t.Fail()
   }
   if wr.last.Next != nil {
      t.Logf("Watch.last.next is non-nil after first Add()\n")
      t.Fail()
   }
   if wr.current != w1 {
      t.Logf("Watch.current is wrong after second Add()\n")
      t.Fail()
   }
   //fmt.Printf("%s\n", wr)
   w3 := NewWatch()
   e3 := NewEntry([]int{1,2,5})
   w3.E = e3
   wr.Add(w3)
   if wr.first != w1 {
      t.Logf("Watch.first is wrong after second Add()\n")
      t.Fail()
   }
   if wr.first.Prev != nil {
      t.Logf("Watch.first.prev is non-nil after first Add()\n")
      t.Fail()
   }
   if wr.last != w3 {
      t.Logf("Watch.last is wrong after second Add()\n")
      t.Fail()
   }
   if wr.last.Next != nil {
      t.Logf("Watch.last.next is non-nil after first Add()\n")
      t.Fail()
   }
   if wr.current != w1 {
      t.Logf("Watch.current is wrong after second Add()\n")
      t.Fail()
   }
   //fmt.Printf("%s\n", wr)
}


func TestNextPrevCurrent(t *testing.T) {
   wr := NewWatchList()
   w1 := NewWatch()
   w1.E = NewEntry([]int{1,2,3})
   w2 := NewWatch()
   w2.E = NewEntry([]int{2,3,4})
   w3 := NewWatch()
   w3.E = NewEntry([]int{3,4,5})

   if wr.Current() != nil || wr.Next() != nil || wr.Prev() != nil {
      t.Logf("Current/Next/Prev initializes incorrectly\n")
      t.Fail()
   }
   wr.Add(w1)
   wr.Add(w2)
   wr.Add(w3)

   if wr.Current() != w1 {
      t.Logf("Current is incorrect\n")
      t.Fail()
   }
   if wr.Next() != w2 {
      t.Logf("Next is incorrect\n")
      t.Fail()
   }
   if wr.Current() != w2 {
      t.Logf("Current after Next is incorrect\n")
      t.Fail()
   }
   if wr.Next() != w3 {
      t.Logf("Next2 is incorrect\n")
      t.Fail()
   }
   if wr.Current() != w3 {
      t.Logf("Current after Next2 is incorrect\n")
      t.Fail()
   }
   if wr.Next() != nil {
      t.Logf("Next3 is non-nil\n")
      t.Fail()
   }
   if wr.Last() != w3 {
      t.Logf("Last is incorrect\n")
      t.Fail()
   }
   if wr.Current() != w3 {
      t.Logf("Current after last is incorrect\n")
      t.Fail()
   }
   if wr.Prev() != w2 {
      t.Logf("Prev() after last is incorrect\n")
      t.Fail()
   }
   if wr.Current() != w2 {
      t.Logf("Current after prev is incorrect\n")
      t.Fail()
   }
   if wr.Prev() != w1 {
      t.Logf("Prev2 is incorrect\n")
      t.Fail()
   }
   if wr.Current() != w1 {
      t.Logf("Current after prev2 is incorrect\n")
      t.Fail()
   }
   if wr.Prev() != nil {
      t.Logf("Prev2 is incorrect\n")
      t.Fail()
   }
   if wr.Current() != nil {
      t.Logf("Current after prev2 is incorrect\n")
      t.Fail()
   }
   if wr.First() != w1 {
      t.Logf("First is incorrect\n")
      t.Fail()
   }
   if wr.Current() != w1 {
      t.Logf("Current after First is incorrect\n")
      t.Fail()
   }
}



func TestRemove(t *testing.T) {
   wr := NewWatchList()
   w1 := NewWatch()
   w1.E = NewEntry([]int{1,2,3})
   w2 := NewWatch()
   w2.E = NewEntry([]int{2,3,4})
   w3 := NewWatch()
   w3.E = NewEntry([]int{3,4,5})


   // Remove from an empty List
   if wr.Remove() != nil {
      t.Logf("WatchList.Remove returns something even though list was empty\n")
      t.Fail()
   }

   wr.Add(w1)
   // Remove only element
   if wr.Remove() != w1 {
      t.Logf("WatchList.Remove returns the wrong thing\n")
      t.Fail()
   }
   if wr.first != nil {
      t.Logf("WatchList.first incorrect after removing from 1-element list\n")
      t.Fail()
   }
   if wr.last != nil {
      t.Logf("WatchList.last incorrect after removing from 1-element list\n")
      t.Fail()
   }
   if wr.current != nil {
      t.Logf("WatchList.current incorrect after removing from 1-element list\n")
      t.Fail()
   }

   wr.Clear()
   wr.Add(w1)
   wr.Add(w2)
   // Remove first element
   if wr.Remove() != w1 {
      t.Logf("WatchList.Remove e1 returns the wrong thing\n")
      t.Fail()
   }
   if wr.first != w2 {
      t.Logf("WatchList.first incorrect after removing e1 2-element list\n")
      t.Fail()
   }
   if wr.last != w2 {
      t.Logf("WatchList.last incorrect after removing e1 2-element list\n")
      t.Fail()
   }
   if wr.current != w2 {
      t.Logf("WatchList.current incorrect after removing e1 2-element list\n")
      t.Fail()
   }

   wr.Clear()
   wr.Add(w1)
   wr.Add(w2)
   wr.Next()
   // Remove Second element
   if wr.Remove() != w2 {
      t.Logf("WatchList.Remove e2 returns the wrong thing\n")
      t.Fail()
   }
   if wr.first != w1 {
      t.Logf("WatchList.first incorrect after removing e2 2-element list\n")
      t.Fail()
   }
   if wr.last != w1 {
      t.Logf("WatchList.last incorrect after removing e2 2-element list\n")
      t.Fail()
   }
   if wr.current != nil {
      t.Logf("WatchList.current incorrect after removing e2 2-element list\n")
      t.Fail()
   }

   wr.Clear()
   wr.Add(w1)
   wr.Add(w2)
   wr.Add(w3)
   wr.Next()
   // Remove Second element
   if wr.Remove() != w2 {
      t.Logf("WatchList.Remove e2 on 3e list returns the wrong thing\n")
      t.Fail()
   }
   if wr.first != w1 {
      t.Logf("WatchList.first incorrect after removing e2 from 3-element list\n")
      t.Fail()
   }
   if wr.last != w3 {
      t.Logf("WatchList.last incorrect after removing e2 from 3-element list\n")
      t.Fail()
   }
   if wr.current != w3 {
      t.Logf("WatchList.current incorrect after removing e2 from 3-element list\n")
      t.Fail()
   }
}


