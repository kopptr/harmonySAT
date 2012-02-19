package dpll

import (
	"bytes"
	"dpll/db"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

type Proportions struct {
	Binary   float64 // Having exactly two literals
	Ternary  float64 // Having exactly three literals
	Horn     float64 // Having <= one positive literal
	Definite float64 // Having exactly one positive literal
	Lowest   float64
	Low      float64
	High     float64
	Highest  float64
}

func NewLProportions(cdb *db.DB) *Proportions {
	p := new(Proportions)
	total := float64(cdb.NLearned())
	p.Binary = float64(cdb.LStats.Binary) / total
	p.Ternary = float64(cdb.LStats.Ternary) / total
	p.Horn = float64(cdb.LStats.Horn) / total
	p.Definite = float64(cdb.LStats.Definite) / total
	cs := cdb.GetCountStats()
	p.Highest = cs.Highest
	p.High = cs.High
	p.Low = cs.Low
	p.Lowest = cs.Lowest
	return p
}

func NewGProportions(cdb *db.DB) *Proportions {
	p := new(Proportions)
	total := float64(cdb.NGiven())
	p.Binary = float64(cdb.GStats.Binary) / total
	p.Ternary = float64(cdb.GStats.Ternary) / total
	p.Horn = float64(cdb.GStats.Horn) / total
	p.Definite = float64(cdb.GStats.Definite) / total
	cs := cdb.GetCountStats()
	p.Highest = cs.Highest
	p.High = cs.High
	p.Low = cs.Low
	p.Lowest = cs.Lowest
	return p
}

type Config struct {
	Dbms   db.ClauseDBMS
	Branch BranchRule
}

type Entry struct {
	Config
	Proportions
}

type Adapter struct {
	entries   []Entry
	nChanges  int
	firstCall bool
}

func NewAdapter(jsonFile string) *Adapter {
	var e error
	a := new(Adapter)
	a.entries = make([]Entry, 2)
	jsonReader, err := os.Open(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	jsonD := json.NewDecoder(jsonReader)
	for i := range a.entries {
		e = jsonD.Decode(&(a.entries[i]))
		if e != nil {
			break
		}
	}

	a.nChanges = -1 // The first choice doesn't count as a change.
	a.firstCall = true

	return a
}

func (a *Adapter) Reconfigure(cdb *db.DB, b *Brancher, m *db.Manager) {
	var (
		p         *Proportions
		bestI     = -1
		bestD     = math.Inf(1) // +infty, all distances should be smaller
		d         float64
		originalB = b.Rule()
		originalM = m.Strat()
	)

	if a.firstCall {
		p = NewGProportions(cdb)
	} else {
		p = NewLProportions(cdb)
		if cdb.NLearned() < 3 && !a.firstCall {
			return
		}
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
		//fmt.Printf("Changed rules from {%s,%s} to {%s,%s}\n", originalB,originalM,b.Rule(),m.Strat())
		a.nChanges++
	}
	a.firstCall = false
	return
}

func (a *Adapter) NChanges() int {
	return a.nChanges
}

func EuclideanDist(p1 *Proportions, p2 *Proportions) float64 {
	return math.Abs(math.Sqrt(
		math.Pow((p1.Binary-p2.Binary), 2.0) +
			math.Pow((p1.Ternary-p2.Ternary), 2.0) +
			math.Pow((p1.Horn-p2.Horn), 2.0) +
			math.Pow((p1.Definite-p2.Definite), 2.0) +
			math.Pow((p1.Lowest-p2.Lowest), 2.0) +
			math.Pow((p1.Low-p2.Low), 2.0) +
			math.Pow((p1.High-p2.High), 2.0) +
			math.Pow((p1.Highest-p2.Highest), 2.0)))
}

func AnalyzeTexString(cdb *db.DB) string {

	p := NewGProportions(cdb)
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "%.0f & %.0f & %.0f & %.0f & ", p.Binary*float64(100), p.Ternary*float64(100), p.Horn*float64(100), p.Definite*float64(100))

	cs := cdb.GetCountStats()
	fmt.Fprintf(buffer, "%f & %f & %f & %f ", cs.Highest, cs.High, cs.Low, cs.Lowest)
	return string(buffer.Bytes())
}

func (p Proportions) String() string {
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "Binary: %f, Ternary: %f, Horn: %f, Definite: %f\n", p.Binary, p.Ternary, p.Horn, p.Definite)
	return string(buffer.Bytes())
}
