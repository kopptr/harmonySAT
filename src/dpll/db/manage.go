package db

import (
	"dpll/assignment/guess"
	"errors"
	"fmt"
)

type ClauseDBMS byte

const (
	None ClauseDBMS = iota
	Queue
	BerkMin
)

var manageFuncs = [...]func(*DB, *guess.Guess, *Manager){none, queue, berkmin}

type Manager struct {
	Manage     func(*DB, *guess.Guess, *Manager)
   strat ClauseDBMS
	MaxLearned uint
}

func NewManager() (m *Manager) {
	m = new(Manager)
	m.SetStrat(None)
	return
}

func (m *Manager) SetStrat(d ClauseDBMS) {
	m.Manage = manageFuncs[d]
   m.strat = d
}

func (m *Manager) Strat() ClauseDBMS {
   return m.strat
}

// Performs the basic management that is not specific to any particular
// strategy, for instance, dividing the counts in the VSIDS counter
func (m *Manager) basic(db *DB, g *guess.Guess) {
	db.Counts.DivCounts(uint(3))
}

func none(db *DB, g *guess.Guess, m *Manager) {
	m.basic(db, g)
	return
}

func queue(db *DB, g *guess.Guess, m *Manager) {
	m.basic(db, g)
	for db.NLearned() > m.MaxLearned {
		db.DelEntry(db.End)
	}
}

// View db as a queue. From the first 15/16 newest clauses, delete anything with
// more than 42 literals. From the last 1/16, delete anything with more than 8.
// This is a simplification of Berkmin. We'd need to add the notion of clause
// activity in order to do it properly. We'd also need restarts.
// Call it a TODO
func berkmin(cdb *DB, g *guess.Guess, m *Manager) {
	var (
		beginning = (g.NAssigned() / 16) * 15
		count     = 0
		tmp       *Entry
	)
	for e := cdb.Learned; e != nil; e = e.Next {
		count++
		if count > beginning {
			if len(e.Clause.Lits) > 8 {
				tmp = e.Prev
				cdb.DelEntry(e)
				e = tmp
			}
		} else {
			if len(e.Clause.Lits) > 42 {
				tmp = e.Prev
				cdb.DelEntry(e)
				e = tmp
			}
		}
	}
}

// Manager needs to satisfy the flag.Value interface
func (m Manager) String() (s string) {
   return m.strat.String()
}
func (d ClauseDBMS) String() (s string) {
   switch d {
   case None: s = "none"
   case Queue: s = "queue"
   case BerkMin: s = "berkmin"
   default: s = "unimplemented"
   }
   return
}

func (m *Manager) Set(s string) error {
	switch s {
	case "":
		return nil
	case "none":
		m.SetStrat(None)
	case "queue":
		m.SetStrat(Queue)
	case "berkmin":
		m.SetStrat(BerkMin)
	default:
		return errors.New(fmt.Sprintf("\"Set\" given invalid value: %s", s))
	}
	return nil
}
