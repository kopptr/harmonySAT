package db

import (
   "bytes"
   "fmt"
   "cnf"
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
   nGiven uint
   nLearned uint
   WatchRings []*Watch
   learning bool
}

func (db *DB) NLearned() uint {
   return db.nLearned
}
func (db *DB) NGiven() uint {
   return db.nGiven
}

func NewDB(nVars int) (db *DB) {
   db = new(DB)
   db.Binary, db.Ternary, db.Horn, db.Definite = 0,0,0,0
   db.Learned, db.Given = nil, nil
   db.nLearned, db.nGiven = 0,0
   db.WatchRings = make([]*Watch, 2*nVars)
   for i := range db.WatchRings {
      db.WatchRings[i] = NewWatch()
//      db.WatchRings[i].Init()
   }
   return
}

// Adds an entry to the database. If called before the call to StartLearning(),
// it adds it to the section of given clauses. If added after, it adds it to the
// set of learned clauses.
func (db *DB) AddEntry(vars []int) {
   e := NewEntry(vars)

   // Insert into list
   if !db.learning {
      // If this is the first clause, point the learned at it so we don't have
      // to traverse the database when we start learning to append to it.
      if db.nGiven == 0 {
         db.Learned = e
         db.Given = e
      }
      // Insert given clauses at front
      e.Prev = nil
      if db.nGiven != 0 {
         e.Next = db.Given
      }
      db.Given.Prev = e
      db.Given = e
      db.nGiven++
   } else {
      if db.nLearned == 0 {
         // If this is the first learned clause, db.Learned points to the last
         // given clause in the List.
         e.Prev = db.Learned
         e.Next = nil
         db.Learned.Next = e
         db.Learned = e
      } else {
         // Otherwise insert at the back of the given/front of the learned.
         e.Next = db.Learned
         e.Prev = db.Learned.Prev
         db.Learned.Prev.Next = e
         db.Learned.Prev = e
         db.Learned = e
      }
      db.nLearned++
   }

   // Add to watches
   for i := range e.Watches {
      var indx uint
      if e.Lits[i].Pol == cnf.Neg {
         indx = e.Lits[i].Val
      } else {
         indx = e.Lits[i].Val + uint((len(db.WatchRings)/2))
      }

      watch := db.WatchRings[indx]
      e.Watches[i].Watching.Val = e.Lits[i].Val
      e.Watches[i].Watching.Pol = e.Lits[i].Pol
      e.Watches[i].Next = watch.Next
      e.Watches[i].Prev = watch
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
   db.nLearned--
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
      if e == db.Learned && db.learning {
         fmt.Fprintf(buffer, "Learned:\n")
      }
		fmt.Fprintf(buffer, "%s", e.Clause)
	}
	fmt.Fprintf(buffer, "\n")
	return string(buffer.Bytes())
}

