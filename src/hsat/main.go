package main

import (
	"dimacs"
	"dpll"
	"dpll/assignment"
	"dpll/db"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

// flags
var (
	seed    int64
	file    string
	logFile string
	cpuprof string
	quiet   bool
	branch  *dpll.Brancher = dpll.NewBrancher()
	manage  *db.Manager    = db.NewManager()
)

func main() {

	flag.Int64Var(&seed, "seed", time.Now().Unix(), "random number generator seed")
	flag.StringVar(&file, "file", "", "dimacs file containing formula")
	flag.StringVar(&logFile, "log", "", "Log output file")
	flag.StringVar(&cpuprof, "cpuprofile", "", "write cpu profile to file")
	flag.BoolVar(&quiet, "q", false, "True for quiet output. States \"SAT\" or \"UNSAT\"")
	flag.Var(branch, "branch", "DPLL branching rule")
	flag.Var(manage, "dbms", "DPLL clause database management strategy")
	flag.Parse()

	rand.Seed(seed)

	err := initLogging(logFile)
	if err != nil {
		log.Fatal(err)
	}

	if cpuprof != "" {
		f, err := os.Create(cpuprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the db
	db, nVars, err := dimacs.DimacsToDb(f)
	f.Close()
	if err != nil {
		log.Fatal("Failed to parse input correctly: %s\n", err.Error())
	}

	db.StartLearning()
   if !quiet {
      fmt.Printf("c Loaded %d clauses into the database\n", db.NGiven())
   }

	// Initialize the assignment
	a := assignment.NewAssignment(nVars)

	// Set the proper max db size
	manage.MaxLearned = db.NGiven() / 3

	// DPLL!
	g := dpll.Dpll(db, a, branch, manage)
	if g == nil {
		if quiet {
			fmt.Printf("UNSAT\n")
		} else {
			fmt.Printf("s UNSAT\n")
		}
	} else {
		ok := db.Verify(g)
		if ok {
			if quiet {
				fmt.Printf("SAT\n")
			} else {
				fmt.Printf("c Solution verified\n")
			}
		} else {
			if quiet {
				fmt.Printf("UNSAT\n")
			} else {
				fmt.Printf("ERROR: Solution could not be verified\n")
			}
		}
		if !quiet {
			fmt.Printf("s SAT\n")
			fmt.Printf("%s\n", g)
		}
	}
	return
}

func initLogging(s string) error {
	if s != "" {
		// No prefix to logged strings
		log.SetFlags(0)
		log.SetPrefix("")
		lf, err := os.Create(logFile)
		if err != nil {
			return err
		}
		log.SetOutput(lf)
	}
	return nil
}
