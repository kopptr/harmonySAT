package main

import (
   "dpll"
   "dpll/db"
   "time"
   "os"
   "fmt"
   "log"
   "encoding/json"
   "dpll/assignment/guess"
   "runtime"
)


func benchmarkFormula(formulaFile string, texFile string, jsonFile string) {
   var (
      bestDuration time.Duration
      bestBr dpll.BranchRule
      bestDbms db.ClauseDBMS
      e dpll.Entry
   )
   branches := [...]dpll.BranchRule{dpll.Ordered, dpll.Random, dpll.Vsids, dpll.Moms}
   dbms := [...]db.ClauseDBMS{db.Queue, db.BerkMin}

   tex,err := os.OpenFile(texFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
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

   writeTableHeader(tex)
   fmt.Fprintf(tex,analyzeFormula(formulaFile))
   writeBetweenTables(tex)

   for _,b := range branches {
      for _, m := range dbms {
         // Run the bench
         before := time.Now()
         g := runBench(file, b, m)
         after := time.Now()
         if g == nil {
            fmt.Fprintf(tex, "%s & %s & TO\\\\\\hline\n",b,m)
         } else {
            thisRun := after.Sub(before)
            fmt.Fprintf(tex, "%s & %s & %s\\\\\\hline\n",b,m,thisRun)
            if thisRun.Nanoseconds() < bestDuration.Nanoseconds() || bestDbms == 0 {
               bestDuration = thisRun
               bestBr = b
               bestDbms = m
            }
         }
      }
   }
   writeTableFooter(tex, formulaFile)

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
   fmt.Fprintf(tex,"\\begin{table}[ht!]\n\\centering\n\\subfloat[][]{\n")
}

func writeBetweenTables(tex *os.File) {
   fmt.Fprintf(tex,"}\n")
   fmt.Fprintf(tex,"\\subfloat[][]{\n")
   fmt.Fprintf(tex,"\\begin{tabular}{|c|c|c|}\\hline\n")
   fmt.Fprintf(tex,"Branch & DBMS & Time\\\\\\hline\\hline\n")
}


func writeTableFooter(tex *os.File, fName string) {
   fmt.Fprintf(tex,"\\end{tabular}\n}\n")
   fmt.Fprintf(tex,"\\caption{%s}\n", fName)
   fmt.Fprintf(tex,"\\end{table}\n")
}

func runBench(file string, b dpll.BranchRule, d db.ClauseDBMS) *guess.Guess {

   runtime.GOMAXPROCS(3)
   timeout := time.After(20*time.Minute)

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

