package db

import (
	"dpll/assignment/guess"
	"dpll/db/cnf"
	"errors"
	"fmt"
)

type CountStats struct {
   P75to100 float64
   P50to74 float64
   P25to49 float64
   P1to24 float64
}

// TODO reimplement as a heap
type LitCounts struct {
	counts []int
}

func NewLitCounts(nVar int) (lc *LitCounts) {
	lc = new(LitCounts)
	lc.counts = make([]int, nVar*2)
	return
}

func (db *DB) GetCountStats() (cs *CountStats) {
   var (
      percentage float64
      c75to100, c50to74, c25to49, c1to24 int
      nClauses = float64(db.nLearned + db.nGiven)
      nLits = float64(len(db.Counts.counts))
   )
   cs = new(CountStats)

   for _, n := range db.Counts.counts {
      percentage = float64(n)/nClauses
      if percentage >= 0.03 {
         c75to100++
      } else if percentage >= 0.015 {
         c50to74++
      } else if percentage >= 0.0075 {
         c25to49++
      } else {
         c1to24++
      }
   }
   cs.P75to100 = float64(c75to100) / nLits
   cs.P50to74 = float64(c50to74) / nLits
   cs.P25to49 = float64(c25to49) / nLits
   cs.P1to24 = float64(c1to24) / nLits
   return
}

func (lc *LitCounts) Get(l *cnf.Lit) (int, error) {
	if l.Val > uint(len(lc.counts)/2) || l.Val < 1 || (l.Pol != cnf.Pos && l.Pol != cnf.Neg) {
		return 0, errors.New(fmt.Sprintf("Called LitCounts.Get with %s", l))
	}
	if l.Pol == cnf.Pos {
		return lc.counts[(2*l.Val)-1], nil
	} else {
		return lc.counts[l.Val-1], nil
	}
	panic("LitCount.Max is horribly broken\n")
}

func (lc *LitCounts) Add(vars []int) {
	for _, v := range vars {
		if v < 0 {
			lc.counts[(-1*v)-1]++
		} else if v > 0 {
			lc.counts[(len(lc.counts)/2)+v-1]++
		}
	}
}

func (lc *LitCounts) DivCounts(divisor int) {
	for i := range lc.counts {
		lc.counts[i] /= divisor
	}
}

func (lc *LitCounts) Max(g *guess.Guess) *cnf.Lit {
	var (
		bestI int
		bestV int
	)

	bestV = -1
	for i, v := range lc.counts {
		if int(v) > bestV {
			if i < len(lc.counts)/2 { // negative polarity
				if p, _ := g.Get(uint(i + 1)); p == guess.Unassigned {
					bestV = int(v)
					bestI = i
				}
			} else { // positive polarity
				if p, _ := g.Get(uint(i - (len(lc.counts) / 2) + 1)); p == guess.Unassigned {
					bestV = int(v)
					bestI = i
				}
			}
		}
	}
	if bestV == -1 {
		return &cnf.Lit{0, 0}
	} else if bestI < len(lc.counts)/2 {
		return &cnf.Lit{uint(bestI + 1), guess.Neg}
	} else {
		return &cnf.Lit{uint(bestI + 1 - (len(lc.counts) / 2)), guess.Pos}
	}
	panic("LitCount.Max is horribly broken\n")
}
