package main

import (
	"dpll"
	"dpll/assignment/guess"
	"dpll/db"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func benchmarkFormula(formulaFile string, texFile string, jsonFile string) {
	var (
		bestDuration time.Duration
		bestBr       dpll.BranchRule
		bestDbms     db.ClauseDBMS
		e            dpll.Entry
	)
	branches := [...]dpll.BranchRule{dpll.Ordered, dpll.Random, dpll.Vsids, dpll.Moms, dpll.Vmtf}
	dbms := [...]db.ClauseDBMS{db.Queue, db.BerkMin}

	tex, err := os.OpenFile(texFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}
	jsonWriter, err := os.OpenFile(jsonFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonWriter.Close()
	defer tex.Close()

	jsonE := json.NewEncoder(jsonWriter)

	fmt.Fprintf(tex, analyzeFormula(formulaFile))

	for _, b := range branches {
		for _, m := range dbms {
			// Run the bench
			before := time.Now()
			g := runBench(file, b, m)
			after := time.Now()
			if g == nil {
				fmt.Fprintf(tex, "TO & ")
			} else {
				thisRun := after.Sub(before)
				fmt.Fprintf(tex, "& %s ", thisRun)
				if thisRun.Nanoseconds() < bestDuration.Nanoseconds() || bestDbms == 0 {
					bestDuration = thisRun
					bestBr = b
					bestDbms = m
				}
			}
		}
	}
	fmt.Fprintf(tex, "\\\\\\hline")

	// Write the json
	if bestDbms != 0 {
		db, _, err := initSolver(file)
		if err != nil {
			log.Fatal(err)
		}
		e.Proportions = *(dpll.NewGProportions(db))
		e.Config.Dbms = bestDbms
		e.Config.Branch = bestBr
		jsonE.Encode(e)
	}

}

func writeTableHeader(tex *os.File) {
	fmt.Fprintf(tex, "\\begin{table}[ht!]\n\\centering\n\n")
}

func writeBetweenTables(tex *os.File) {
	fmt.Fprintf(tex, "}\n")
	fmt.Fprintf(tex, "\\subfloat[][]{\n")
	fmt.Fprintf(tex, "\\begin{tabular}{|c|c|c|}\\hline\n")
	fmt.Fprintf(tex, "Branch & DBMS & Time\\\\\\hline\\hline\n")
}

func writeTableFooter(tex *os.File, fName string) {
	fmt.Fprintf(tex, "\\end{tabular}\n}\n")
	fmt.Fprintf(tex, "\\caption{%s}\n", fName)
	fmt.Fprintf(tex, "\\end{table}\n")
}

func runBench(file string, b dpll.BranchRule, d db.ClauseDBMS) *guess.Guess {

	runtime.GOMAXPROCS(3)
	timeout := time.After(20 * time.Minute)

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

	g := dpll.DpllTimeout(cdb, a, br, ma, nil, timeout)
	return g
}

func runAdaptiveBench(file string, jsonFile string, chooseOnce bool, extraStats bool) (*guess.Guess, *dpll.Adapter) {

	runtime.GOMAXPROCS(3)
	timeout := time.After(20 * time.Minute)

	// Initialize the cdb and assignment
	cdb, a, err := initSolver(file)
	if err != nil {
		log.Fatal(err)
	}
	// Set the proper max db size
	manage.MaxLearned = cdb.NGiven() / 3
	// Read the data from the file and make the adapter
	adapt := dpll.NewAdapter(jsonFile, chooseOnce, extraStats)
	// Use the adapter to set the initial state
	b := dpll.NewBrancher()
	m := db.NewManager()
	adapt.Reconfigure(cdb, b, m)

	g := dpll.DpllTimeout(cdb, a, b, m, adapt, timeout)
	return g, adapt
}

func testFormula(formulaFile string, texFile string, jsonFile string) {
	branches := [...]dpll.BranchRule{dpll.Ordered, dpll.Random, dpll.Vsids, dpll.Moms}
	dbms := [...]db.ClauseDBMS{db.Queue, db.BerkMin}

	tex, err := os.OpenFile(texFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}
	defer tex.Close()

	fmt.Fprintf(tex, analyzeFormula(formulaFile))

	for _, b := range branches {
		for _, m := range dbms {
			// Run the bench
			before := time.Now()
			g := runBench(file, b, m)
			after := time.Now()
			if g == nil {
				fmt.Fprintf(tex, "& TO ")
			} else {
				thisRun := after.Sub(before)
				fmt.Fprintf(tex, "& %s ", thisRun)
			}
		}
	}

	chooseOnce := true
	extraStats := false
	for i := 0; i < 4; i++ {
		// get all 4 combos
		if i > 1 {
			chooseOnce = false
		}
		if i%2 == 1 {
			extraStats = !extraStats
		}
		// Run the Adaptive bench
		before := time.Now()
		g, a := runAdaptiveBench(file, jsonFile, chooseOnce, extraStats)
		after := time.Now()
		if g == nil {
         if chooseOnce {
            fmt.Fprintf(tex, "& TO ")
         } else {
            fmt.Fprintf(tex, "& TO & --- ")
         }
		} else {
			thisRun := after.Sub(before)
         if chooseOnce {
            fmt.Fprintf(tex, "& %s ", thisRun)
         } else {
            fmt.Fprintf(tex, "& %s & %d", thisRun, a.NChanges())
         }
		}
	}
	fmt.Fprintf(tex, "\\\\\\hline")

}
