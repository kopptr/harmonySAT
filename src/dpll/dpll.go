package dpll

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
)

type dpllStackNode struct {
	l       *cnf.Lit
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
			} else if status == db.Sat {
				return a.Guess()
			} else {
				break
			}
		}
	}
	panic("Dpll is broken")
}


func DpllTimeout(cdb *db.DB, a *assignment.Assignment, b *Brancher, m *db.Manager, timeout chan bool) *guess.Guess {

	nVar := a.Guess().Len()
	aStack := make([]dpllStackNode, nVar)
	top := -1

	for {
      select {
      case <-timeout:
         return nil
      default:
      }

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
			} else if status == db.Sat {
				return a.Guess()
			} else {
				break
			}
		}
	}
	panic("Dpll is broken")
}
