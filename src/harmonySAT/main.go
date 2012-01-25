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
   seed = flag.Int64("seed", time.Nanoseconds(), "seed for random number generator")
)

func main() {

   flag.Parse()
   rand.Seed(*seed)

   f, err := os.Open("dimacs/test/1.cnf")
   if err != nil {
      fmt.Printf("Error opening file\n")
   }

   db, nVars, ok := dimacs.DimacsToDb(f)
   if db == nil || !ok {
      fmt.Printf("Failed to parse input correctly.\n")
   }
   db.StartLearning()

   a := assignment.NewAssignment(nVars)
   g := dpll.Dpll(db, a)
   fmt.Printf("%s\n", g)
   return
}
