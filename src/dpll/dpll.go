package dpll

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
	"fmt"
)

func Dpll(db *db.DB, a *assignment.Assignment) *guess.Guess {
	var g *guess.Guess

	l := decide(db, a)
	if l.Eq(&cnf.Lit{0, 0}) {
		return a.Guess()
	}
	a.PushAssign(l.Val, l.Pol)
	fmt.Printf("PUSH: %s\n", l)
	ok := db.Bcp(a.Guess(), *l)

	if ok {
		g = Dpll(db, a)
		if g != nil {
			return g
		}
	}

	// try the reverse polarity
	a.PopAssign()
	fmt.Printf("POP: %s\n", l)
	l.Flip()
	a.PushAssign(l.Val, l.Pol)
	fmt.Printf("PUSH: %s\n", l)
	ok = db.Bcp(a.Guess(), *l)
	if ok {
		g = Dpll(db, a)
		if g != nil {
			return g
		}
	}
	return nil
}

func decide(db *db.DB, a *assignment.Assignment) (l *cnf.Lit) {
	// find the first in-order unassigned literal
	for i := uint(1); i <= a.Len(); i++ {
		if p, _ := a.Get(i); p == guess.Unassigned {
			return &cnf.Lit{i, cnf.Pos}
		}
	}
	return &cnf.Lit{0, 0}
}
