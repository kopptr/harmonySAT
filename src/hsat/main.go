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
	seed       int64
	file       string
	logFile    string
	cpuprof    string
	quiet      bool
	analyze    bool
	extraStats bool
	chooseOnce bool
	benchmark  bool
	adaptive   string
	branch     *dpll.Brancher = dpll.NewBrancher()
	manage     *db.Manager    = db.NewManager()
)

func main() {

	flag.Int64Var(&seed, "seed", time.Now().Unix(), "random number generator seed")
	flag.StringVar(&file, "file", "", "dimacs file containing formula")
	flag.StringVar(&logFile, "log", "", "Log output file")
	flag.StringVar(&cpuprof, "cpuprofile", "", "write cpu profile to file")
	flag.BoolVar(&quiet, "q", false, "True for quiet output. States \"SAT\" or \"UNSAT\"")
	flag.BoolVar(&analyze, "e", false, "True for examination output. If true, will not actually run solver.")
	flag.BoolVar(&benchmark, "b", false, "True for benchmark output.")
	flag.BoolVar(&extraStats, "s", false, "True for extra statistics.")
	flag.BoolVar(&chooseOnce, "c", false, "True for choosing once, then not adapting.")
	flag.StringVar(&adaptive, "a", "", "Path to analysis json file. Enables adaptive solving")
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
	if benchmark {
		if adaptive == "" {
			benchmarkFormula(file, "output.tex", "analysis.json")
		} else {
			testFormula(file, "adapt-output.tex", adaptive)
		}
		return
	} else if analyze {
		a := analyzeFormula(file)
		fmt.Printf("%s\n", a)
		return
	} else if adaptive != "" {
		runAdaptiveSolver(file, adaptive, chooseOnce, extraStats, quiet)
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

	return dpll.AnalyzeTexString(db)
}

func runAdaptiveSolver(file string, json string, chooseOnce bool, extraStats, quiet bool) {
	// Initialize the cdb and assignment
	cdb, a, err := initSolver(file)
	if err != nil {
		log.Fatal(err)
	}
	// Set the proper max db size
	manage.MaxLearned = cdb.NGiven() / 3
	// Read the data from the file and make the adapter
	adapt := dpll.NewAdapter(json, chooseOnce, extraStats)
	// Use the adapter to set the initial state
	b := dpll.NewBrancher()
	adapt.Reconfigure(cdb, b, manage)

	g := dpll.Dpll(cdb, a, b, manage, adapt)

	printResults(g, cdb, adapt, !quiet)
}

func runNormalSolver(file string, b *dpll.Brancher, m *db.Manager, quiet bool) {
	// Initialize the cdb and assignment
	db, a, err := initSolver(file)
	if err != nil {
		log.Fatal(err)
	}
	// Set the proper max db size
	m.MaxLearned = db.NGiven() / 3

	// Run the Dpll!
	g := dpll.Dpll(db, a, b, m, nil)

	printResults(g, db, nil, !quiet)
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
		return nil, nil, err
	}

	// Initialize the db
	cdb, a, err = dimacs.DimacsToDb(f)
	f.Close()
	if err != nil {
		return nil, nil, err
	}

	cdb.StartLearning()

	return
}

func runBaseSolver() {

}

func printResults(g *guess.Guess, cdb *db.DB, a *dpll.Adapter, verbose bool) {
	if verbose {
		printVerboseResults(g, cdb, a)
	} else {
		printQuietResults(g, cdb)
	}
}

func printQuietResults(g *guess.Guess, cdb *db.DB) {
	if g == nil {
		fmt.Printf("UNSAT\n")
	} else if !cdb.Verify(g) {
		fmt.Printf("UNKNOWN\n")
	} else {
		fmt.Printf("SAT\n")
	}
}

func printVerboseResults(g *guess.Guess, cdb *db.DB, a *dpll.Adapter) {
	if g == nil {
		fmt.Printf("s UNSAT\n")
	} else if !cdb.Verify(g) {
		fmt.Printf("c ERROR: Solution could not be verified\n")
		fmt.Printf("s UNKNOWN\n")
		fmt.Printf("%s\n", g)
	} else {
		fmt.Printf("c Solution verified\n")
		if adaptive != "" {
			fmt.Printf("c Adaptive solver changed strategies %d times.\n", a.NChanges())
		}
		fmt.Printf("s SAT\n")
		fmt.Printf("s %s\n", g)
	}
}
