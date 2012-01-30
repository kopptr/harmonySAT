package db

import (
   "dpll/db/cnf"
   "dpll/assignment/guess"
   "errors"
   "fmt"
)

// TODO reimplement as a heap
type LitCounts struct {
   counts []uint
}


func NewLitCounts(nVar int) (lc *LitCounts) {
   lc = new(LitCounts)
   lc.counts = make([]uint, nVar*2)
   return
}

func (lc *LitCounts) Get(l *cnf.Lit) (uint, error) {
   if l.Val > uint(len(lc.counts)/2) || l.Val < 1 || (l.Pol != cnf.Pos && l.Pol != cnf.Neg) {
      return 0, errors.New(fmt.Sprintf("Called LitCounts.Get with %s",l))
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

func (lc *LitCounts) DivCounts(divisor uint) {
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
   for i,v := range lc.counts {
      if int(v) > bestV {
         if i < len(lc.counts)/2 { // negative polarity
            if p, _ := g.Get(uint(i+1)); p == guess.Unassigned {
               bestV = int(v)
               bestI = i
            }
         } else { // positive polarity
            if p,_ := g.Get(uint(i-(len(lc.counts)/2)+1)); p == guess.Unassigned {
                  bestV = int(v)
                  bestI = i
            }
         }
      }
   }
   if bestV == -1 {
           return &cnf.Lit{0,0}
   } else if bestI < len(lc.counts)/2 {
      return &cnf.Lit{uint(bestI+1),guess.Neg}
   } else {
      return &cnf.Lit{uint(bestI+1-(len(lc.counts)/2)),guess.Pos}
   }
   panic("LitCount.Max is horribly broken\n")
}



