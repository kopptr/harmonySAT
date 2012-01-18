package db

import (
   "bytes"
   "fmt"
)

// Statistics for a database of clauses
type Stats struct {
	Binary   uint // Having exactly two literals
	Ternary  uint // Having exactly three literals
	Horn     uint // Having <= one positive literal
	Definite uint // Having exactly one positive literal
}

// A database of clauses.
type DB struct {
   Stats
   Given *Entry
   Learned *Entry
   NGiven uint
   NLearned uint
   WatchRings []Watch
   learning bool
}

func (db *DB) New(nVars int) {
   db.Binary, db.Ternary, db.Horn, db.Definite = 0,0,0,0
   db.NLearned, db.NGiven = 0,0
   db.WatchRings = make([]Watch, 2*nVars)
   for i := range db.WatchRings {
      db.WatchRings[i].New()
   }
}

// Adds an entry to the database. If called before the call to StartLearning(),
// it adds it to the section of given clauses. If added after, it adds it to the
// set of learned clauses.
func (db *DB) AddEntry(vars []int) {
   var insertionPt *Entry
   e := NewEntry(vars)

   // Insert into list
   if !db.learning {
      insertionPt = db.Given
      db.NGiven++
   } else {
      insertionPt = db.Learned
      insertionPt.Prev.Next = e
      db.NLearned++
   }
   e.Prev = insertionPt.Prev
   e.Next = insertionPt
   insertionPt.Prev = e

   // Add to watches
   for i := range e.Watches {
      watch := db.WatchRings[e.Lits[i].Val]
      e.Watches[i].Watching.Val = e.Lits[i].Val
      e.Watches[i].Watching.Pol = e.Lits[i].Pol
      e.Watches[i].Next = watch.Next
      e.Watches[i].Prev = &watch
      watch.Next.Prev = &e.Watches[i]
      watch.Next = &e.Watches[i]
   }

   // Update DB stats
   if e.IsBinary() {
      db.Binary++
   }
   if e.IsTernary() {
      db.Ternary++
   }
   if e.IsHorn() {
      db.Horn++
   }
   if e.IsDefinite() {
      db.Definite++
   }
}

func (db *DB) DelEntry(e *Entry) {
   // Bookkeeping
   db.NLearned--
   // Remove from watches
   for i := range e.Watches {
      e.Watches[i].Prev.Next = e.Watches[i].Next
      e.Watches[i].Next.Prev = e.Watches[i].Prev
   }
   // Remove from List
   if e.Next != nil {
      e.Next.Prev = e.Prev
   }
   if e.Prev != nil { // Should be redundant
      e.Prev.Next = e.Next
   }
   // Update stats
   if e.IsBinary() {
      db.Binary--
   }
   if e.IsTernary() {
      db.Ternary--
   }
   if e.IsHorn() {
      db.Horn--
   }
   if e.IsDefinite() {
      db.Definite--
   }
   // The caller just needs to get rid of his reference to the object
   // i.e. set e = nil, then the gc should get it.
}


func (db *DB) StartLearning() {
   db.learning = true
}

func (db *DB) String() string {
	buffer := bytes.NewBufferString("")
   fmt.Fprintf(buffer, "Given:\n")
   for e := db.Given; e != nil; e = e.Next {
      if e == db.Learned {
         fmt.Fprintf(buffer, "Learned:\n")
      }
		fmt.Fprintf(buffer, "%s\n", e.Clause)
	}
	fmt.Fprintf(buffer, "\n")
	return string(buffer.Bytes())
}

