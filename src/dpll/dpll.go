package dpll

import (
   "time"
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
)

type dpllStackNode struct {
	l       *cnf.Lit
	Flipped bool
}

func Dpll(cdb *db.DB, a *assignment.Assignment, b *Brancher, m *db.Manager, adapt *Adapter) *guess.Guess {

	nVar := a.Guess().Len()
	aStack := make([]dpllStackNode, nVar)
	top := -1

	for {
		top++
		aStack[top].l = b.Decide(cdb, a)
		aStack[top].Flipped = false
		a.PushAssign(aStack[top].l.Val, aStack[top].l.Pol)

		for {
			status := cdb.Bcp(a.Guess(), *aStack[top].l, m)
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

            // Reconfigure if neccessary
            if adapt != nil {
               adapt.Reconfigure(cdb, b, m)
            }
			} else if status == db.Sat {
				return a.Guess()
			} else {
				break
			}
		}
	}
	panic("Dpll is broken")
}


func DpllTimeout(cdb *db.DB, a *assignment.Assignment, b *Brancher, m *db.Manager, adapt *Adapter, timeout <-chan time.Time) *guess.Guess {

	nVar := a.Guess().Len()
	aStack := make([]dpllStackNode, nVar)
	top := -1

	for {

		top++
		aStack[top].l = b.Decide(cdb, a)
		aStack[top].Flipped = false
		a.PushAssign(aStack[top].l.Val, aStack[top].l.Pol)

		for {
         select {
         case <-timeout:
            return nil
         default:
         }

			status := cdb.Bcp(a.Guess(), *aStack[top].l, m)
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
            // Reconfigure if neccessary
            if adapt != nil {
               adapt.Reconfigure(cdb, b, m)
            }

			} else if status == db.Sat {
				return a.Guess()
			} else {
				break
			}
		}
	}
	panic("Dpll is broken")
}
