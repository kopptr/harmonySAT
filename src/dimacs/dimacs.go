package dimacs

import (
   "io"
   "dpll/db"
   "dimacs/scanner"
)

type tType int
const (
   tErr tType = iota
   tComment
   tClause
   tLit
)




func DimacsToDb(r io.Reader) (clauseDB *db.DB, nVar int, ok bool) {
   s := scanner.NewScanner(r)

   ok = eatDimacsComments(s)
   if !ok {
      return nil, 0, false
   }
   nVar, nClauses, ok := eatFormulaInfo(s)
   if !ok {

      return nil, 0, false
   }
   clauseDB = db.NewDB(nVar)
   n, ok := eatClauses(s, clauseDB)
   if !ok || n != nClauses {
      return nil, 0, false
   }

   return
}

func eatClauses(r *scanner.Scanner, clauseDB *db.DB) (n int, ok bool) {
   for n = 0; r.HasNextLine(); n++ {
      clause := []int{}
      for i := r.NextInt(); i != 0; i = r.NextInt() {
         clause = append(clause,i)
      }
      clauseDB.AddEntry(clause)
   }
   return n, true
}

func eatFormulaInfo(r *scanner.Scanner) (nVar int, nClauses int, ok bool) {
   // cnf
   c := r.Next()
   if c != "cnf" {
      panic("p line did not have cnf string\n")
      return -1, -1, false
   }

   // nVar
   nVar = r.NextInt()
   if nVar <= 0 {
      return -1, -1, false
   }

   // nClauses
   nClauses = r.NextInt()
   if nClauses <= 0 {
      return -1, -1, false
   }

   return nVar, nClauses, true
}

func eatDimacsComments(r *scanner.Scanner) bool {
   for {
      c := r.Next()
      if c != "c" {
         if c != "p" {
            panic("First non-comment line does not begine with a p\n")
            return false
         } else {
            return true
         }
      }
      _ = r.NextLine()
   }
   // Should never get here
   return false
}
