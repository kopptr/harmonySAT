package db

import (
	"dpll/assignment/guess"
	"dpll/assignment"
	"dpll/db/cnf"
	"errors"
	"fmt"
)

type CountStats struct {
	Highest float64
	High    float64
	Low     float64
	Lowest  float64
}

type pair struct {
   count int
   lit int
}

type LitCounts struct {
	counts []int
   vmtf []pair
}

func NewLitCounts(nVar int) (lc *LitCounts) {
	lc = new(LitCounts)
	lc.counts = make([]int, nVar*2)
	lc.vmtf = make([]pair, nVar*2)
   for i := 0; i < nVar; i++ {
      lc.vmtf[i].count = 0;
      lc.vmtf[i].lit = i+1;
   }
   for i := nVar; i < 2*nVar; i++ {
      lc.vmtf[i].count = 0;
      lc.vmtf[i].lit = -1*(i-nVar+1);
   }
	return
}

func (db *DB) GetCountStats() (cs *CountStats) {
	var (
		percentage                         float64
		c75to100, c50to74, c25to49, c1to24 int
		nClauses                           = float64(db.nLearned + db.nGiven)
		nLits                              = float64(len(db.Counts.counts))
	)
	cs = new(CountStats)

	for _, n := range db.Counts.counts {
		percentage = float64(n) / nClauses
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
	cs.Highest = float64(c75to100) / nLits
	cs.High = float64(c50to74) / nLits
	cs.Low = float64(c25to49) / nLits
	cs.Lowest = float64(c1to24) / nLits
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

func (lc *LitCounts) AddGiven(vars []int) {
	for _, v := range vars {
		if v < 0 {
			lc.counts[(-1*v)-1]++
			lc.vmtf[(-1*v)-1].count++
		} else if v > 0 {
			lc.counts[(len(lc.counts)/2)+v-1]++
			lc.vmtf[(len(lc.counts)/2)+v-1].count++
		}
	}
}

func (lc *LitCounts) sortVmtf() {
   var (
      biggest int
      tmp pair
   )
   for i := range lc.vmtf {
      biggest = i
      for j := i; j < len(lc.vmtf); j++ {
         if lc.vmtf[i].count > lc.vmtf[biggest].count {
            biggest = i
         }
      }
      tmp = lc.vmtf[i]
      lc.vmtf[i] = lc.vmtf[biggest]
      lc.vmtf[biggest] = tmp
   }
}

const varsToMove = 8
func (lc *LitCounts) reorderVmtf(c []int) {
   var (
      n int
      place = 0
      tmp pair
   )
   if len(c) < varsToMove {
      n = len(c)
   } else {
      n = varsToMove
   }
   for i := range lc.vmtf {
      for j := range c {
         if c[i] == lc.vmtf[j].lit {
            tmp = lc.vmtf[place]
            lc.vmtf[place] = lc.vmtf[j]
            lc.vmtf[j] = tmp
            place++
            break
         }
      }
      if place == n {
         break
      }
   }
}

func (lc *LitCounts) GetNextVmtf(a *assignment.Assignment) (l *cnf.Lit) {
   var lit uint
   for i := range lc.vmtf {
      if lc.vmtf[i].lit < 0 {
         lit = uint(-1*lc.vmtf[i].lit)
      } else {
         lit = uint(lc.vmtf[i].lit)
      }
      if val, err := a.Get(lit); val == guess.Unassigned && err == nil {
         if lc.vmtf[i].lit < 0 {
            return &cnf.Lit{uint(lc.vmtf[i].lit*-1), cnf.Pos}
         } else {
            return &cnf.Lit{uint(lc.vmtf[i].lit), cnf.Pos}
         }
      }
   }
   panic("getNextVmtf is broken")
}

func (lc *LitCounts) Add(vars []int) {
	for _, v := range vars {
		if v < 0 {
			lc.counts[(-1*v)-1]++
		} else if v > 0 {
			lc.counts[(len(lc.counts)/2)+v-1]++
		}
	}
   lc.reorderVmtf(vars)
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
