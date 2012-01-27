package dpll

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
	"log"
)

func Dpll(db *db.DB, a *assignment.Assignment) *guess.Guess {
	var g *guess.Guess

	l := decide(db, a)
	if l.Eq(&cnf.Lit{0, 0}) {
		return a.Guess()
	}
	a.PushAssign(l.Val, l.Pol)
	log.Printf("%sPUSH: %s\n", indent(a), l)
	ok := db.Bcp(a.Guess(), *l)
        log.Printf("Guess: %s%s\n", indent(a), a.Guess())

	if ok {
		g = Dpll(db, a)
		if g != nil {
			return g
		}
	}

	// try the reverse polarity
	a.PopAssign()
	log.Printf("%sPOP: %s\n", indent(a), l)
        log.Printf("Guess: %s%s\n", indent(a), a.Guess())
	l.Flip()
	a.PushAssign(l.Val, l.Pol)
	log.Printf("%sPUSH: %s\n", indent(a), l)
        log.Printf("Guess: %s%s\n", indent(a), a.Guess())
	ok = db.Bcp(a.Guess(), *l)
	if ok {
		g = Dpll(db, a)
		if g != nil {
			return g
		}
	}

	a.PopAssign()
	log.Printf("%sPOP: %s\n", indent(a), l)
        log.Printf("Guess: %s%s\n", indent(a), a.Guess())
	return nil
}

func indent(a *assignment.Assignment) string {
        s := ""
        for i := uint(0); i < a.Depth(); i++ {
                s += "\t"
        }
        return s
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
