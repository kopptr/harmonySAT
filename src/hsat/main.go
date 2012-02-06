package main

import (
   "config"
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
	benchmark bool
	adaptive bool
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
	flag.BoolVar(&benchmark, "b", false, "True for benchmark output.")
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

   // Special switches for training behavior
   if analyze {
      a := analyzeFormula(file)
      fmt.Printf("%s\n", a)
      return
   } else if benchmark {
      benchmarkFormula(file, "tex.tex", "gob.gob")
      return
   } else if adaptive {
      runAdaptiveSolver(file)
   } else {
      runNormalSolver(file, branch, manage, quiet)
   }

	return
}

func analyzeFormula(file string) string {
   // Initialize the cdb and assignment
   db, _, err := initSolver(file)
   if err != nil {
      log.Fatal(err)
   }

   return config.AnalyzeTexString(db)
}


func runBench(file string, b dpll.BranchRule, d db.ClauseDBMS, timeout chan bool) *guess.Guess {
   // Prepare the specific run
   br := dpll.NewBrancher()
   br.SetRule(b)
   ma := db.NewManager()
   ma.SetStrat(d)

   // Initialize the cdb and assignment
   cdb, a, err := initSolver(file)
   if err != nil {
      log.Fatal(err)
   }
   // Set the proper max db size
   ma.MaxLearned = cdb.NGiven() / 3

   g := dpll.DpllTimeout(cdb, a, br, ma, timeout)
   return g
}

func runAdaptiveSolver(file string) {
   // Initialize the cdb and assignment
   db, _, err := initSolver(file)
   if err != nil {
      log.Fatal(err)
   }
	// Set the proper max db size
	manage.MaxLearned = db.NGiven() / 3

}

func runNormalSolver(file string, b *dpll.Brancher, m *db.Manager, quiet bool) {
   // Initialize the cdb and assignment
   db, a, err := initSolver(file)
   if err != nil {
      log.Fatal(err)
   }
	// Set the proper max db size
	manage.MaxLearned = db.NGiven() / 3

   // Run the Dpll!
	g := dpll.Dpll(db, a, b, m)

   printResults(g, db, !quiet)
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

