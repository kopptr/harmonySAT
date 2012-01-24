package dpll

import (
   "dpll/db"
   "dpll/db/cnf"
   "dpll/assignment"
   "dpll/assignment/guess"
   "fmt"
)

func Dpll( db *db.DB, a *assignment.Assignment ) *guess.Guess {
   var g *guess.Guess

   fmt.Printf("0\n")
   l := decide(db, a)
   fmt.Printf("1\n")
   if l.Eq(&cnf.Lit{0,0}) {
      return a.Guess()
   }
   fmt.Printf("2\n")
   a.PushAssign(l.Val, l.Pol)
   fmt.Printf("3\n")
   ok := db.Bcp(a.Guess(), *l)
   fmt.Printf("4\n")

   if ok {
      g = Dpll(db, a)
      if g != nil {
         return g
      }
   }

   // try the reverse polarity
   a.PopAssign()
   l.Flip()
   a.PushAssign(l.Val, l.Pol)
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
   for i := a.Depth(); i < a.Len(); i++ {
      if a.Get(i+1) == guess.Unassigned {
         return &cnf.Lit{i+1, cnf.Pos}
      }
   }
   return &cnf.Lit{0,0}
}
