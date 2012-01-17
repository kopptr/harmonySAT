package db

import (
   "cnf"
)

/* An entry in the database. Essentially a clause with metadata about
 * two-literal watches.
 */
type Entry struct {
   *cnf.Clause
   Watches [2]Watch
}

// Creates a new Entry. Initializes the clause data and allocates the watches.
// Does not add the watches to any data structure.
func NewEntry(vars []int) (e *Entry) {
   e = new(Entry)
   e.Clause = cnf.NewClause(vars)
   e.Watches[0].New()
   e.Watches[1].New()
   return
}


