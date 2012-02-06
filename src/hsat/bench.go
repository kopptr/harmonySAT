package main

import (
   "dpll"
   "dpll/db"
   "time"
   "os"
   "fmt"
   "log"
   "encoding/gob"
   "config"
)


func benchmarkFormula(formulaFile string, texFile string, gobFile string) {
   var (
      bestDuration time.Duration
      bestBr dpll.BranchRule
      bestDbms db.ClauseDBMS
      e config.Entry
   )
   branches := [...]dpll.BranchRule{dpll.Ordered, dpll.Random, dpll.Vsids, dpll.Moms}
   dbms := [...]db.ClauseDBMS{db.Queue, db.BerkMin}

   tex,err := os.OpenFile(texFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
   if err != nil {
      log.Fatal(err)
   }
   gobWriter, err := os.OpenFile(gobFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
   if err != nil {
      log.Fatal(err)
   }
   defer gobWriter.Close()
   defer tex.Close()

   gobE := gob.NewEncoder(gobWriter)

   writeTableHeader(tex)
   fmt.Fprintf(tex,analyzeFormula(formulaFile))
   writeBetweenTables(tex)

   for _,b := range branches {
      for _, m := range dbms {
         // Start the timeout checker
         timeout := make(chan bool, 1)
         go func() {
            // 20 minutes
            time.Sleep(1200e9)
            timeout <- true
         }()

         // Run the bench
         before := time.Now()
         g := runBench(file, b, m, timeout)
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
   fmt.Printf("best config: %s, %s\n", bestBr, bestDbms)
   fmt.Printf("time: %s", bestDuration)

   // Write the gob
   db, _, err := initSolver(file)
   if err != nil {
      log.Fatal(err)
   }
   e.Proportions = *(config.NewProportions(db))
   e.Config.Dbms = bestDbms
   e.Config.Branch = bestBr
   gobE.Encode(e)

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
