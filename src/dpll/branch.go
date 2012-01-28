package dpll

import (
   "dpll/assignment"
   "dpll/assignment/guess"
   "dpll/db"
   "dpll/db/cnf"
   "math/rand"
   "errors"
   "fmt"
)

type BranchRule byte
const (
   Ordered BranchRule = iota
   Random
)

var branchFuncs = [...]func(*db.DB, *assignment.Assignment) *cnf.Lit { ordered, random }

type Brancher struct {
   Decide func(*db.DB, *assignment.Assignment) *cnf.Lit
}

func NewBrancher() (b *Brancher) {
   b = new(Brancher)
   b.SetRule(Ordered)
   return
}

func (b *Brancher) SetRule(r BranchRule) {
   b.Decide = branchFuncs[r]
}

func ordered(db *db.DB, a *assignment.Assignment) (l *cnf.Lit) {
	// find the first in-order unassigned literal
	for i := uint(1); i <= a.Len(); i++ {
		if p, _ := a.Get(i); p == guess.Unassigned {
			return &cnf.Lit{i, cnf.Pos}
		}
	}
	return &cnf.Lit{0, 0}
}

func random(db *db.DB, a *assignment.Assignment) (l *cnf.Lit) {
   sign := byte(rand.Int() % 2)
   val := uint((rand.Int() % int(a.Len()))+1)
   for i := val; i != val; i++ {
      if v,_ := a.Get(val); v == guess.Unassigned {
         return &cnf.Lit{i,sign}
      } else if i == a.Len() {
         i = 0
      }
   }
   return &cnf.Lit{0,0}
}

// Brancher needs to satisfy the flag.Value interface
func (b Brancher) String() string {
   return ""
}

func (b *Brancher) Set(s string) error {
   switch s {
   case "": return nil
   case "ordered": b.SetRule(Ordered)
   case "random": b.SetRule(Random)
   default: return errors.New(fmt.Sprintf("\"Set\" given invalid value: %s", s))
   }
   return nil
}

