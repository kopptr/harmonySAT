package main

import (
	"dimacs"
	"dpll"
	"dpll/assignment"
	"dpll/assignment/guess"
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
	analyze bool
	branch  *dpll.Brancher = dpll.NewBrancher()
	manage  *db.Manager    = db.NewManager()
)

func main() {

	flag.Int64Var(&seed, "seed", time.Now().Unix(), "random number generator seed")
	flag.StringVar(&file, "file", "", "dimacs file containing formula")
	flag.StringVar(&logFile, "log", "", "Log output file")
	flag.StringVar(&cpuprof, "cpuprofile", "", "write cpu profile to file")
	flag.BoolVar(&quiet, "q", false, "True for quiet output. States \"SAT\" or \"UNSAT\"")
	flag.BoolVar(&analyze, "a", false, "True for analysis output. If true, will not actually run solver.")
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

   // Initialize the cdb and assignment
   db, a, err := initSolver(file)
   if err != nil {
      log.Fatal(err)
   }
	if analyze {
		fmt.Printf("%s", db.AnalyzeTexString())
		return
	}


	// Set the proper max db size
	manage.MaxLearned = db.NGiven() / 3

	// DPLL!
	g := dpll.Dpll(db, a, branch, manage)

   printResults(g, db, !quiet)

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

func initSolver(file string) (cdb *db.DB, a *assignment.Assignment, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil,nil,err
	}

	// Initialize the db
	cdb, nVars, err := dimacs.DimacsToDb(f)
	f.Close()
	if err != nil {
		return nil,nil,err
	}

	cdb.StartLearning()

	// Initialize the assignment
	a = assignment.NewAssignment(nVars)
   return
}

func runBaseSolver() {

}


func printResults(g *guess.Guess, cdb *db.DB, verbose bool) {
   if verbose {
      printVerboseResults(g, cdb)
   } else {
      printQuietResults(g, cdb)
   }
}

func printQuietResults(g *guess.Guess, cdb *db.DB) {
   if g == nil {
      fmt.Printf("UNSAT\n")
   } else if !cdb.Verify(g) {
      fmt.Printf("UNKOWN\n")
   } else {
      fmt.Printf("SAT\n")
   }
}

func printVerboseResults(g *guess.Guess, cdb *db.DB) {
	if g == nil {
      fmt.Printf("s UNSAT\n")
	} else if !cdb.Verify(g) {
      fmt.Printf("c ERROR: Solution could not be verified\n")
      fmt.Printf("s UNKNOWN\n")
      fmt.Printf("%s\n", g)
   } else {
      fmt.Printf("c Solution verified\n")
      fmt.Printf("s SAT\n")
      fmt.Printf("%s\n", g)
	}
}

