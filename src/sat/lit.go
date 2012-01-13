package sat

import (
   "strconv"
)

/* Represents a propositional variable, or literal. Val is a unique id for the
 * variable, Pol is it's polarity, Pos(itive), Neg(ative), or Unassigned.
 */
type Lit struct {
   Val uint
   Pol byte
}

const (
   Unassigned byte = 0
   Pos byte = 1
   Neg byte = 2
)

func (l Lit) String() (s string) {
   if l.Pol == Pos {
      s = strconv.Uitoa(l.Val)
   } else if l.Pol == Neg {
      s = "-" + strconv.Uitoa(l.Val)
   }
   return
}

/* Returns true if the Lit's polarity is Pos or Neg
 */
func (l *Lit) IsSet() bool {
   return l.Pol == Pos || l.Pol == Neg
}

/* Reverses the polarity of the Lit. If Lit.Pol is unassigned, does nothing.
 */
func (l *Lit) Flip() {
   if l.Pol == Pos {
      l.Pol = Neg
   } else if l.Pol == Neg {
      l.Pol = Pos
   }
}


