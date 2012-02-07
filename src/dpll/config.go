package dpll

import (
   "math"
   "dpll/db"
   "fmt"
   "bytes"
   "os"
   "log"
   "encoding/json"
)

type Proportions struct {
	Binary   float64 // Having exactly two literals
	Ternary  float64 // Having exactly three literals
	Horn     float64 // Having <= one positive literal
	Definite float64 // Having exactly one positive literal
}

func NewProportions(cdb *db.DB) *Proportions {
   p := new(Proportions)
   if cdb.IsLearning() {
      total := float64(cdb.NLearned())
      p.Binary = float64(cdb.LStats.Binary)/total
      p.Ternary = float64(cdb.LStats.Ternary)/total
      p.Horn = float64(cdb.LStats.Horn)/total
      p.Definite = float64(cdb.LStats.Definite)/total
   } else {
      total := float64(cdb.NGiven())
      p.Binary = float64(cdb.GStats.Binary)/total
      p.Ternary = float64(cdb.GStats.Ternary)/total
      p.Horn = float64(cdb.GStats.Horn)/total
      p.Definite = float64(cdb.GStats.Definite)/total
   }
   return p
}

type Config struct {
   Dbms db.ClauseDBMS
   Branch BranchRule
}

type Entry struct {
   Config
   Proportions
}

type Adapter struct {
   entries  []Entry
   nChanges int
}

func NewAdapter(jsonFile string) *Adapter {
   var e error
   a := new(Adapter)
   a.entries = make([]Entry,2)
   jsonReader, err:= os.Open(jsonFile)
   if err != nil {
      log.Fatal(err)
   }
   jsonD := json.NewDecoder(jsonReader)
   for i := range a.entries {
      e = jsonD.Decode(&(a.entries[i]))
      if e != nil {
         break;
      }
   }

   a.nChanges = -1 // The first choice doesn't count as a change.

   return a
}


func (a *Adapter) Reconfigure(cdb *db.DB, b *Brancher, m *db.Manager) {
   var (
      p = NewProportions(cdb)
      bestI = -1
      bestD = math.Inf(1) // +infty, all distances should be smaller
      d float64
      originalB = b.Rule()
      originalM = m.Strat()
   )

   if cdb.NLearned() < 3 {
      return
   }

   // Find the best match
   for i := range a.entries {
      d = EuclideanDist(p, &a.entries[i].Proportions)
      if d < bestD {
         bestD = d
         bestI = i
      }
   }

   // Apply it
   if originalM != a.entries[bestI].Config.Dbms || originalB != a.entries[bestI].Config.Branch {
      m.SetStrat(a.entries[bestI].Config.Dbms)
      b.SetRule(a.entries[bestI].Config.Branch)
      fmt.Printf("Changed rules from {%s,%s} to {%s,%s}\n", originalB,originalM,b.Rule(),m.Strat())
      a.nChanges++
   }
   return
}

func (a *Adapter) NChanges() int {
   return a.nChanges
}


func EuclideanDist(p1 *Proportions, p2 *Proportions) float64 {
   return math.Abs(math.Sqrt(math.Pow((p1.Binary-p2.Binary), 2.0) +
   (math.Pow((p1.Ternary-p2.Ternary), 2.0)) +
   (math.Pow((p1.Horn-p2.Horn), 2.0)) +
   (math.Pow((p1.Definite-p2.Definite), 2.0))))
}

func AnalyzeTexString(cdb *db.DB) string {

   p := NewProportions(cdb)
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "\\begin{tabular}{|c|c|}\\hline\n")
	fmt.Fprintf(buffer, "Type & Number\\\\\\hline\\hline\n")
	fmt.Fprintf(buffer, "Binary & %f\\\\\\hline\n", p.Binary)
	fmt.Fprintf(buffer, "Ternary & %f\\\\\\hline\n", p.Ternary)
	fmt.Fprintf(buffer, "Horn & %f\\\\\\hline\n", p.Horn)
	fmt.Fprintf(buffer, "Definite & %f\\\\\\hline\n", p.Definite)
	fmt.Fprintf(buffer, "\\end{tabular}")
	return string(buffer.Bytes())
}

func (p Proportions) String() string {
   buffer := bytes.NewBufferString("")
   fmt.Fprintf(buffer, "Binary: %f, Ternary: %f, Horn: %f, Definite: %f\n", p.Binary, p.Ternary, p.Horn, p.Definite)
   return string(buffer.Bytes())
}
