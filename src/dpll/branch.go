package dpll

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
	"errors"
	"fmt"
	"math/rand"
)

type BranchRule byte

const (
	Ordered BranchRule = iota
	Random
	Vsids
	Moms
)

var branchFuncs = [...]func(*db.DB, *assignment.Assignment) *cnf.Lit{ordered, random, vsids, moms}

type Brancher struct {
	Decide func(*db.DB, *assignment.Assignment) *cnf.Lit
   rule BranchRule
}

func NewBrancher() (b *Brancher) {
	b = new(Brancher)
	b.SetRule(Ordered)
	return
}

func (b *Brancher) SetRule(r BranchRule) {
	b.Decide = branchFuncs[r]
   b.rule = r
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
	sign := byte((rand.Int() % 2) + 1)
	val := uint((rand.Int() % int(a.Len())) + 1)
	for i := val; i <= a.Len(); i++ {
		if v, _ := a.Get(i); v == guess.Unassigned {
			return &cnf.Lit{i, sign}
		}
	}
	for i := uint(1); i < val; i++ {
		if v, _ := a.Get(i); v == guess.Unassigned {
			return &cnf.Lit{i, sign}
		}
	}
	return &cnf.Lit{0, 0}
}

func vsids(db *db.DB, a *assignment.Assignment) (l *cnf.Lit) {
	l = db.Counts.Max(a.Guess())
	return
}

// We'll use a modified MOMS rule that only looks at binary clauses
func moms(db *db.DB, a *assignment.Assignment) (l *cnf.Lit) {
	var (
		counts   = make([]int, a.Guess().Len())
		eCount   int
		vals     [2]uint
		g        = a.Guess()
		biggest  = 0
		biggestI = -1
	)

	// Count up the total lits in binary clauses
	for e := db.Learned; e != nil; e = e.Next {
		eCount = 0
		for _, l := range e.Clause.Lits {
			if v, _ := g.Get(l.Val); v == guess.Unassigned {
				eCount++
				if eCount > 2 {
					break
				}
				vals[eCount-1] = l.Val
			}
		}
		if eCount <= 2 {
			for i := 0; i < eCount; i++ {
				counts[vals[i]-1]++
			}
		}
	}

	// Search for the biggest
	for i, v := range counts {
		if v > biggest {
			biggest = v
			biggestI = i
		}
	}

	if biggestI == -1 {
		return random(db, a)
	} else {
		return &cnf.Lit{uint(biggestI + 1), byte((rand.Int() % 2) + 1)}
	}
	panic("MOMS is broken")
}

// Brancher needs to satisfy the flag.Value interface
func (b Brancher) String() string {
   return b.rule.String()
}

func (b BranchRule) String() (s string) {
   switch b {
   case Ordered: s = "ordered"
   case Random: s = "random"
   case Vsids: s = "vsids"
   case Moms: s = "moms"
   default: s = "unimplemented"
   }
   return
}

func (b *Brancher) Set(s string) error {
	switch s {
	case "":
		return nil
	case "ordered":
		b.SetRule(Ordered)
	case "random":
		b.SetRule(Random)
	case "vsids":
		b.SetRule(Vsids)
	case "moms":
		b.SetRule(Moms)
	default:
		return errors.New(fmt.Sprintf("\"Set\" given invalid value: %s", s))
	}
	return nil
}
