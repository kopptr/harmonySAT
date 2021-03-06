package db

import (
	"dpll/assignment/guess"
	"dpll/db/cnf"
)

type BcpReturn int

const (
	Okay BcpReturn = iota
	Conflict
	Sat
)

func (db *DB) Bcp(g *guess.Guess, lit cnf.Lit, m *Manager) BcpReturn {
	lq := newLitQ() // queue of literals to be assigned
	lq.PushBack(lit)
	return db.LQBcp(g, lq, m)
}

// Performs BCP on the clause database until no unit clauses remain.
// Returns false iff the formula is conflicted.
func (db *DB) LQBcp(g *guess.Guess, lq *LitQ, m *Manager) BcpReturn {

	// For each new literal in the queue
	for l, ok := lq.PopFront(); ok; l, ok = lq.PopFront() {
		// Each clause watching the literal was just satisfied.
		// For each clause watching the reverse polarity of the literal
		reverse := cnf.Lit{l.Val, l.Pol}
		reverse.Flip()
		wl := db.GetWatchList(reverse)
		for wl.First(); wl.Current() != nil; wl.Next() {
			// We need to watch something else iff the other watch is unsatisfied
			// Check if it's satisfied
			otherWatch := wl.Current().Other()
			if p, _ := g.Get(otherWatch.Watching.Val); p == otherWatch.Watching.Pol {
				// The whole clause is therefore satisfied
				continue
			}
			// We must try to find a new literal to watch
			found := false // found a new literal to watch
			// for each literal in the clause
			for _, newL := range wl.Current().E.Clause.Lits {
				// If the other watch is watching it, this one cannot
				if otherWatch.Watching.Eq(&newL) {
					continue
				}
				// If it is assigned in the correct polarity or unassigned
				// Watch it
				if p, _ := g.Get(newL.Val); p == newL.Pol || p == guess.Unassigned {
					w := wl.Current()
					db.Pluck(w)
					w.Watching.Pol = newL.Pol
					w.Watching.Val = newL.Val
					newWl := db.GetWatchList(newL)
					newWl.Add(w)
					found = true
					break
				}
			}
			// If we found nothing new to watch, it's either a new unit clause or a
			// conflict
			if !found {
				// If unit clause
				if p, _ := g.Get(otherWatch.Watching.Val); p == guess.Unassigned {
					// Add it to the queue
					lq.PushBack(otherWatch.Watching)
					g.Set(otherWatch.Watching.Val, otherWatch.Watching.Pol)
				} else {
					// CONFLICT
					db.AddConflictEntry(g)
					if m != nil {
						m.Manage(db, g, m)
					}
					return Conflict
				}
			}
		}
	}
	if g.NAssigned() == g.Len() {
		if db.Verify(g) {
			return Sat
		} else {
			return Conflict
		}
	}
	return Okay
}

func (db *DB) AddConflictEntry(g *guess.Guess) {
	db.AddEntry(g.Vars(true), false)
}

type LitQ struct {
	First *litQNode
	Last  *litQNode
}

type litQNode struct {
	Next *litQNode
	L    cnf.Lit
}

func newLitQNode(l cnf.Lit) (lqn *litQNode) {
	lqn = new(litQNode)
	lqn.L.Val = l.Val
	lqn.L.Pol = l.Pol
	lqn.Next = nil
	return
}

func newLitQ() (lq *LitQ) {
	lq = new(LitQ)
	lq.First = nil
	lq.Last = nil
	return
}

func (lq *LitQ) PushBack(l cnf.Lit) {
	lqn := newLitQNode(l)
	if lq.First == nil { // empty
		lq.First = lqn
		lq.Last = lqn
	} else {
		lq.Last.Next = lqn
		lq.Last = lqn
	}
}

func (lq *LitQ) PopFront() (l cnf.Lit, ok bool) {
	if lq.First == nil { // empty
		l = cnf.Lit{0, 0}
		ok = false
	} else {
		l = lq.First.L
		ok = true
		lq.First = lq.First.Next
	}
	return
}
