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
)

var manageFuncs = [...]func(*DB, *guess.Guess, *Manager){none, queue}

type Manager struct {
	Manage     func(*DB, *guess.Guess, *Manager)
	MaxLearned uint
}

func NewManager() (m *Manager) {
	m = new(Manager)
	m.SetStrat(None)
	return
}

func (m *Manager) SetStrat(d ClauseDBMS) {
	m.Manage = manageFuncs[d]
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
		db.DelEntry(db.Learned)
	}
}

// Manager needs to satisfy the flag.Value interface
func (m Manager) String() string {
	return ""
}
func (m *Manager) Set(s string) error {
	switch s {
	case "":
		return nil
	case "none":
		m.SetStrat(None)
	case "queue":
		m.SetStrat(Queue)
	default:
		return errors.New(fmt.Sprintf("\"Set\" given invalid value: %s", s))
	}
	return nil
}
