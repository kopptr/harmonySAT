package main

import (
	"dimacs"
	"dpll"
	"dpll/assignment"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	seed = flag.Int64("seed", time.Now().Unix(), "random number generator seed")
	file = flag.String("file", "", "dimacs file containing formula")
	quiet = flag.Bool("q", false,
      "True for quiet output. States \"SAT\" or \"UNSAT\"")
)

func main() {

	flag.Parse()
	rand.Seed(*seed)

	f, err := os.Open(*file)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	db, nVars, ok := dimacs.DimacsToDb(f)
	if db == nil || !ok {
		fmt.Printf("Failed to parse input correctly.\n")
		return
	}
	db.StartLearning()

	a := assignment.NewAssignment(nVars)
	g := dpll.Dpll(db, a)
	if g == nil {
      if *quiet {
         fmt.Printf("UNSAT\n")
      } else {
         fmt.Printf("s UNSAT\n")
      }
	} else {
		ok := db.Verify(g)
		if ok {
         if *quiet {
            fmt.Printf("SAT\n")
         } else {
            fmt.Printf("c Solution verified\n")
         }
		} else {
         if *quiet {
            fmt.Printf("UNSAT\n")
         } else {
            fmt.Printf("ERROR: Solution could not be verified\n")
         }
		}
      if !*quiet {
         fmt.Printf("s SAT\n")
         fmt.Printf("%s\n", g)
      }
	}
	return
}
