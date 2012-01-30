package dpll

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db/cnf"
	"dpll/db"
	"log"
)

func Dpll(db *db.DB, a *assignment.Assignment, b *Brancher, m *db.Manager) *guess.Guess {
	var g *guess.Guess

	l := b.Decide(db, a)
        // Assignment is full
	if l.Eq(&cnf.Lit{0, 0}) {
                if db.Verify(a.Guess()) {
                        // Done
                        return a.Guess()
                } else {
                        // Backtrack
                        return nil
                }
	}
	a.PushAssign(l.Val, l.Pol)
	log.Printf("%sPUSH: %s\n", indent(a), l)
	ok := db.Bcp(a.Guess(), *l, indent(a), m)
        log.Printf("Guess: %s%s\n", indent(a), a.Guess())

	if ok {
		g = Dpll(db, a, b, m)
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
	ok = db.Bcp(a.Guess(), *l, indent(a), m)
	if ok {
		g = Dpll(db, a, b, m)
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

