package cnf

import (
	"sort"
	"bytes"
	"fmt"
)

// Potential attributes of a clause
type ClauseStats struct {
	binary   bool // Having exactly two literals
	ternary  bool // Having exactly three literals
	horn     bool // Having <= one positive literal
	definite bool // Having exactly one positive literal
}

// A clause, or collection of literals. Contains stats about itself that are
// calculated at creation. Sorted by Lit.Val.
type Clause struct {
	ClauseStats
	Lits []Lit
}

// Given a slice of Lits, calculates the stats for them as if they were a
// clause.
func (cs *ClauseStats) SetStats(lits []Lit) {
	cs.binary = false
	cs.ternary = false
	cs.horn = false
	cs.definite = false
	pos, neg, cnt := 0, 0, len(lits)

	for _, l := range lits {
		if l.Pol == Pos {
			pos++
		} else {
			neg++
		}
	}
	if cnt == 2 {
		cs.binary = true
	} else if cnt == 3 {
		cs.ternary = true
	}
	if pos <= 1 {
		cs.horn = true
		if pos == 1 {
			cs.definite = true
		}
	}
}

// Given a collection of integers, converts them into a sorted clause with
// statistics.
func NewClause(vars []int) (c *Clause) {
	c = new(Clause)
	c.Lits = make([]Lit, len(vars))
	for i, v := range vars {
		c.Lits[i] = NewLit(v)
	}
	sort.Sort(c)
	c.SetStats(c.Lits)
	return
}

// Returns true if a particular literal is contained in the clause, as well as
// its polarity.
func (c *Clause) Contains(val uint) (isContained bool, polarity byte) {
	// sort.Search returns the lowest index that satisfies the fn
	in := sort.Search(len(c.Lits), func(i int) bool {
		return c.Lits[i].Val >= val
	})
	if in == len(c.Lits) || c.Lits[in].Val != val {
		return false, 0
	}
	return true, c.Lits[in].Pol
}

// String representation
func (c *Clause) String() string {
	buffer := bytes.NewBufferString("")
	for _, l := range c.Lits {
		fmt.Fprintf(buffer, "%s ", l)
	}
	fmt.Fprintf(buffer, "\n")
	return string(buffer.Bytes())
}

// These functions provide access to the statistical properties of the Clause
func (c *Clause) IsBinary() bool {
	return c.binary
}
func (c *Clause) IsTernary() bool {
	return c.ternary
}
func (c *Clause) IsHorn() bool {
	return c.horn
}
func (c *Clause) IsDefinite() bool {
	return c.definite
}

/* These functions are implemented to implement sort.Interface.
 * They are not useful otherwise useful. It is important to implement
 * sort.Interface so that we can do a binary search for a Lit in the
 * Contains() function.
 */
// Number of elements in collection
func (c Clause) Len() int {
	return len(c.Lits)
}
// Compares the literals' values, disregards polarity
func (c Clause) Less(i, j int) bool {
	return c.Lits[i].Val < c.Lits[j].Val
}
// Swaps two literals
func (c Clause) Swap(i, j int) {
	tmp := Lit{c.Lits[i].Val, c.Lits[i].Pol}
	c.Lits[i].Val = c.Lits[j].Val
	c.Lits[i].Pol = c.Lits[j].Pol
	c.Lits[j].Val = tmp.Val
	c.Lits[j].Pol = tmp.Pol
}
