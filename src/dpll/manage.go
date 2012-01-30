package dpll

import (
   "dpll/assignment"
   "dpll/db"
   "errors"
   "fmt"
)

type ClauseDBMS byte
const (
        None ClauseDBMS = iota
        Queue
)

var manageFuncs = [...]func(*db.DB, *assignment.Assignment, *Manager) {none, queue}

type Manager struct {
        Manage func(*db.DB, *assignment.Assignment, *Manager)
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

func none(*db.DB, *assignment.Assignment, *Manager) {
        return
}

func queue(db *db.DB, a *assignment.Assignment, m *Manager) {
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
        case "": return nil
        case "none": m.SetStrat(None)
        case "queue": m.SetStrat(Queue)
        default: return errors.New(fmt.Sprintf("\"Set\" given invalid value: %s", s))
        }
        return nil
}
