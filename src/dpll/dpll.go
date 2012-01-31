package dpll

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
)
/*
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
	ok := db.Bcp(a.Guess(), *l, indent(a), m)

	if ok {
		g = Dpll(db, a, b, m)
		if g != nil {
			return g
		}
	}

	// try the reverse polarity
	l = a.PopAssign()
	l.Flip()
	a.PushAssign(l.Val, l.Pol)
	ok = db.Bcp(a.Guess(), *l, indent(a), m)
	if ok {
		g = Dpll(db, a, b, m)
		if g != nil {
			return g
		}
	}

	a.PopAssign()
	return nil
}
*/


type dpllStackNode struct {
   l *cnf.Lit
   Flipped bool
}

func Dpll(cdb *db.DB, a *assignment.Assignment, b *Brancher, m *db.Manager) *guess.Guess {

   nVar := a.Guess().Len()
   aStack := make([]dpllStackNode, nVar)
   top := -1

   for {
      top++
      aStack[top].l = b.Decide(cdb, a)
      aStack[top].Flipped = false
      a.PushAssign(aStack[top].l.Val, aStack[top].l.Pol)

      for {
         status := cdb.Bcp(a.Guess(), *aStack[top].l, indent(a), m)
         if status == db.Conflict {
            // BackTrack
            for aStack[top].Flipped == true {
               top--
               a.PopAssign()
            }
            if top < 0 {
               return nil
            }
            // Flip the assignment
            a.PopAssign()
            aStack[top].l.Flip()
            aStack[top].Flipped = true
            a.PushAssign(aStack[top].l.Val, aStack[top].l.Pol)
         } else if status == db.Sat {
            return a.Guess()
         } else {
            break
         }
      }
   }
   panic("Dpll is broken")
}


func indent(a *assignment.Assignment) string {
	s := ""
	for i := uint(0); i < a.Depth(); i++ {
		s += "\t"
	}
	return s
}
