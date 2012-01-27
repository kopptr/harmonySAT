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
	"log"
)

var (
	seed  = flag.Int64("seed", time.Now().Unix(), "random number generator seed")
	file  = flag.String("file", "", "dimacs file containing formula")
        logFile = flag.String("log", "hsat.log", "Log output file")
	quiet = flag.Bool("q", false,
		"True for quiet output. States \"SAT\" or \"UNSAT\"")
)

func main() {

	flag.Parse()
	rand.Seed(*seed)

        err := initLogging()
        if err != nil {
                fmt.Printf("Failed to open log file: %s\n", err.Error())
                return
        }

	f, err := os.Open(*file)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	db, nVars, err := dimacs.DimacsToDb(f)
	if err != nil {
		fmt.Printf("Failed to parse input correctly: %s\n", err.Error())
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

func initLogging() error {
        // No prefix to logged strings
        log.SetFlags(0);
        log.SetPrefix("");
        lf, err := os.Create(*logFile)
        if err != nil {
                return err
        }
        log.SetOutput(lf)
        return nil
}
