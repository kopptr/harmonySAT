package db

import (
   "cnf"
   "fmt"
   "bytes"
)

// A Watch Node, an item in a list that corresponds to a db entry. Designed to
// be used such that each entry has two watches in the database, and can keep
// track of two-literal watches.
type Watch struct {
   Next *Watch
   Prev *Watch
   E *Entry
   Watching cnf.Lit
}

// Returns a new watch Variable.
func NewWatch() (w *Watch) {
   w = new(Watch)
   w.Watching.Val = 0
   w.Watching.Pol = 0
   w.Next = nil
   w.Prev = nil
   w.E = nil
   return
}

// A doubly-linked list of Watches
// The DB will use one Watch Ring per lit per polarity
type WatchList struct {
   first *Watch
   last *Watch
   current *Watch
}

// Returns a new empty WatchList
func NewWatchList() (wr *WatchList) {
   wr = new(WatchList)
   wr.first = nil
   wr.last = nil
   wr.current = nil
   return
}

// Adds a watch to the list
func (wr *WatchList) Add(w *Watch) {
   if wr.first == nil { // empty
      wr.first = w
      wr.last = w
      wr.current = w
   } else {
      w.Prev = wr.last
      w.Next = nil
      wr.last.Next = w
      wr.last = w
   }
}

// Removes the Watch currently pointed to from the list, and returns it.
func (wr *WatchList) Remove() (w *Watch) {
   w = wr.current
   if w != nil {
      if wr.current != wr.first {
         wr.current.Prev.Next = wr.current.Next
      } else {
         wr.first = wr.current.Next
      }
      if wr.current != wr.last {
         wr.current.Next.Prev = wr.current.Prev
         wr.current = wr.current.Next
      } else {
         wr.last = wr.current.Prev
         wr.current = nil
      }
      w.Next = nil
      w.Prev = nil
   }
   return
}

// DO NOT USE
func (w *Watch) Pluck() {
   if w.Prev != nil {
      w.Prev.Next = w.Next
   }
   if w.Next != nil {
      w.Next.Prev = w.Prev
   }
   w.Next = nil
   w.Prev = nil
}

// Removes all watches from the list
func (wr *WatchList) Clear() {
   for wr.First(); wr.Remove() != nil; {}
}


// Returns the watch currently pointed to by the list
func (wr *WatchList) Current() (w *Watch) {
   w = wr.current
   return
}

// Moves the current pointer to the previous watch, and returns it.
func (wr *WatchList) Prev() (w *Watch) {
   if wr.current != nil {
      wr.current = wr.current.Prev
   }
   w = wr.current
   return
}

// Moves the current pointer to the next watch, and returns it.
func (wr *WatchList) Next() (w *Watch) {
   if wr.current != nil {
      wr.current = wr.current.Next
   }
   w = wr.current
   return
}

// Moves the current pointer to the first position, and returns it
func (wr *WatchList) First() (w *Watch) {
   wr.current = wr.first
   w = wr.current
   return
}

// Moves the current pointer to the last position, and returns it
func (wr *WatchList) Last() (w *Watch) {
   wr.current = wr.last
   w = wr.current
   return
}

// String representation
func (wr WatchList) String() string {
   if wr.first == nil {
      return ""
   }
	buffer := bytes.NewBufferString("")
   for w := wr.first; w != nil; w = w.Next {
		fmt.Fprintf(buffer, "%s\n", w.E.Clause.String())
	}
	return string(buffer.Bytes())
}

