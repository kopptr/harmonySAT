package main

import (
   "rand"
   "time"
   "flag"
   "os"
   "fmt"
   "dimacs"
   "dpll"
   "dpll/assignment"
)

var (
   seed = flag.Int64("seed", time.Nanoseconds(), "random number generator seed")
   file = flag.String("file", "", "dimacs file containing formula")
)

func main() {

   flag.Parse()
   rand.Seed(*seed)

   f, err := os.Open(*file)
   if err != nil {
      fmt.Printf("Error opening file\n")
      return
   }

   db, nVars, ok := dimacs.DimacsToDb(f)
   if db == nil || !ok {
      fmt.Printf("Failed to parse input correctly.\n")
   }
   fmt.Printf("s read in %d clauses\n", db.NGiven())

   db.StartLearning()

   a := assignment.NewAssignment(nVars)
   g := dpll.Dpll(db, a)
   if g == nil {
      fmt.Printf("s UNSAT\n")
   } else {
      fmt.Printf("s SAT\n")
      fmt.Printf("%s\n", g)
   }
   return
}
