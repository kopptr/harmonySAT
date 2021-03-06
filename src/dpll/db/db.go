package db

import (
	"bytes"
	"dpll/assignment/guess"
	"dpll/db/cnf"
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
	GStats     Stats
	LStats     Stats
	Counts     *LitCounts
	Given      *Entry
	Learned    *Entry
	End        *Entry
	nGiven     uint
	nLearned   uint
	WatchLists []*WatchList
	learning   bool
}

func (db *DB) AnalyzeTexString() string {
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "\\begin{tabular}{|c|c|}\\hline\n")
	fmt.Fprintf(buffer, "Type & Number\\\\\\hline\\hline\n")
	fmt.Fprintf(buffer, "Binary & %d\\\\\\hline\n", db.GStats.Binary)
	fmt.Fprintf(buffer, "Ternary & %d\\\\\\hline\n", db.GStats.Ternary)
	fmt.Fprintf(buffer, "Horn & %d\\\\\\hline\n", db.GStats.Horn)
	fmt.Fprintf(buffer, "Definite & %d\\\\\\hline\n", db.GStats.Definite)
	fmt.Fprintf(buffer, "\\end{tabular}\n")
	return string(buffer.Bytes())
}

func (db *DB) NLearned() uint {
	return db.nLearned
}

func (db *DB) NGiven() uint {
	return db.nGiven
}

func NewDB(nVars int) (db *DB) {
	db = new(DB)
	db.GStats.Binary, db.GStats.Ternary, db.GStats.Horn, db.GStats.Definite = 0, 0, 0, 0
	db.LStats.Binary, db.LStats.Ternary, db.LStats.Horn, db.LStats.Definite = 0, 0, 0, 0
	db.Learned, db.Given = nil, nil
	db.nLearned, db.nGiven = 0, 0
	db.WatchLists = make([]*WatchList, 2*nVars)
	for i := range db.WatchLists {
		db.WatchLists[i] = NewWatchList()
	}
	db.Counts = NewLitCounts(nVars)
	return
}

// Adds an entry to the database. If called before the call to StartLearning(),
// it adds it to the section of given clauses. If added after, it adds it to the
// set of learned clauses.
func (db *DB) AddEntry(vars []int, shouldSort bool) {
	e := NewEntry(vars, shouldSort)

	// Insert into list
	if !db.learning {
		// If this is the first clause, point the learned at it so we don't have
		// to traverse the database when we start learning to append to it.
		if db.nGiven == 0 {
			db.Learned = e
			db.Given = e
			e.Next = nil
			e.Prev = nil
		} else {
			// Insert given clauses at front
			e.Prev = nil
			e.Next = db.Given
			db.Given.Prev = e
			db.Given = e
		}
		db.nGiven++
	} else {
		if db.nLearned == 0 {
			// If this is the first learned clause, db.Learned points to the last
			// given clause in the List.
			e.Prev = db.Learned
			e.Next = nil
			db.Learned.Next = e
			db.Learned = e
			db.End = e
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
		wl := db.GetWatchList(e.Lits[i])
		e.Watches[i].E = e
		e.Watches[i].Watching.Val = e.Lits[i].Val
		e.Watches[i].Watching.Pol = e.Lits[i].Pol
		wl.Add(e.Watches[i])
	}

	// Update DB stats
	if db.learning {
		if e.IsBinary() {
			db.LStats.Binary++
		}
		if e.IsTernary() {
			db.LStats.Ternary++
		}
		if e.IsHorn() {
			db.LStats.Horn++
		}
		if e.IsDefinite() {
			db.LStats.Definite++
		}
	} else {
		if e.IsBinary() {
			db.GStats.Binary++
		}
		if e.IsTernary() {
			db.GStats.Ternary++
		}
		if e.IsHorn() {
			db.GStats.Horn++
		}
		if e.IsDefinite() {
			db.GStats.Definite++
		}
	}

	// Update Lit Counts
   if db.learning {
      db.Counts.Add(vars)
   } else {
      db.Counts.AddGiven(vars)
   }
}

func (db *DB) DelEntry(e *Entry) {
	// Bookkeeping
	db.nLearned--
	// Remove from watches
	for i := range e.Watches {
		db.Pluck(e.Watches[i])
	}
	// Remove from List
	if e == db.End {
		if e == db.Learned {
			db.End = nil
		} else {
			db.End = e.Prev
		}
	}

	if e == db.Learned {
		db.Learned = e.Next
	}

	if db.nLearned == 0 {
		db.Learned = e.Prev
		db.End = nil
	}

	if e.Next != nil {
		e.Next.Prev = e.Prev
	}
	if e.Prev != nil { // Should be redundant
		e.Prev.Next = e.Next
	}
	// Update stats
	if e.IsBinary() {
		db.LStats.Binary--
	}
	if e.IsTernary() {
		db.LStats.Ternary--
	}
	if e.IsHorn() {
		db.LStats.Horn--
	}
	if e.IsDefinite() {
		db.LStats.Definite--
	}
	// The caller just needs to get rid of his reference to the object
	// i.e. set e = nil, then the gc should get it.
}

func (db *DB) GetWatchList(l cnf.Lit) *WatchList {
	i := int(l.Val - 1)
	if l.Pol == cnf.Pos {
		i += len(db.WatchLists) / 2
	}
	return db.WatchLists[i]
}

func (db *DB) StartLearning() {
	db.learning = true
   db.Counts.sortVmtf()
}

func (db *DB) IsLearning() bool {
	return db.learning
}

func (db *DB) String() string {
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "Given:\n")
	for e := db.Given; e != nil; e = e.Next {
		if e == db.Learned && db.nLearned > 0 {
			fmt.Fprintf(buffer, "Learned:\n")
		}
		fmt.Fprintf(buffer, "%s\n", e.Clause)
	}
	fmt.Fprintf(buffer, "Watches:\n")
	watchNum := -1
	for _, wl := range db.WatchLists {
		fmt.Fprintf(buffer, "Watching %d:\n", watchNum)
		fmt.Fprintf(buffer, "%s", wl)
		if watchNum < 0 {
			watchNum--
		} else {
			watchNum++
		}
		if watchNum < (len(db.WatchLists) / 2 * -1) {
			watchNum = 1
		}
	}
	return string(buffer.Bytes())
}

// Returns true iff the Guess satisfies the DB
func (db *DB) Verify(g *guess.Guess) bool {
	for e := db.Given; e != db.Learned; e = e.Next {
		for _, l := range e.Lits {
			if p, _ := g.Get(l.Val); p == l.Pol {
				goto found
			}
		}
		return false
	found:
	}
	return true
}
