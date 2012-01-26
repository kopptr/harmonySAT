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
	fmt.Printf("c read in %d clauses\n", db.NGiven())

	db.StartLearning()

	a := assignment.NewAssignment(nVars)
	g := dpll.Dpll(db, a)
	if g == nil {
		fmt.Printf("s UNSAT\n")
	} else {
		ok := db.Verify(g)
		if ok {
			fmt.Printf("c Solution verified\n")
		} else {
			fmt.Printf("ERROR: Solution could not be verified\n")
		}
		fmt.Printf("s SAT\n")
		fmt.Printf("%s\n", g)
	}
	return
}
