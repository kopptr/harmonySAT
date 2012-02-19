package dimacs

import (
	"dpll/assignment"
	"dpll/assignment/guess"
	"dpll/db"
	"dpll/db/cnf"
	"errors"
	"io"
	"scanner"
	"strings"
   "fmt"
)

func DimacsToDb(r io.Reader) (clauseDB *db.DB, a *assignment.Assignment, err error) {
	s := scanner.NewScanner(r)

	pLine := matchDimacsComments(s)

	nVar, _, err := matchFormulaInfo(pLine)
	if err != nil {
		return nil, nil, err
	}
	clauseDB = db.NewDB(nVar)
	a = assignment.NewAssignment(nVar)
	n, err := matchClauses(s, clauseDB, a)
	if err != nil {
		return nil, nil, err
	}
	if n == -1 {
		return nil, nil, errors.New("s Unsatisfiable")
	}
	/*
	else if n != nClauses {
			return nil, nil, errors.New(fmt.Sprintf("Read %d/%d clauses.", n, nClauses))
		}
	* This is not an error. Unit clauses are assigned & discarded.
	*/
	return
}

// Matches all of the clauses in the database, and inserts them.
func matchClauses(r *scanner.Scanner, clauseDB *db.DB, a *assignment.Assignment) (n int, err error) {
	foundUnit := false
	n = 0
	lq := new(db.LitQ)
	for r.HasNextLine() {
		clause := []int{}
		for i := r.NextInt(); i != 0; i = r.NextInt() {
			clause = append(clause, i)
		}
		// Add unit clauses directly to assignment
		if len(clause) == 1 {
			foundUnit = true
			if clause[0] < 1 {
				if foo, _ := a.Guess().Get(uint((clause[0] * -1))); foo == guess.Pos {
               fmt.Printf("Conflicting unit clauses: %d\n", clause[0])
					return -1, nil
				}
				a.Guess().Set(uint((clause[0] * -1)), cnf.Neg)
				lq.PushBack(cnf.Lit{uint((clause[0] * -1)), cnf.Neg})
			} else {
				if foo, _ := a.Guess().Get(uint(clause[0])); foo == guess.Neg {
               fmt.Printf("Conflicting unit clauses: %d\n", clause[0])
					return -1, nil
				}
				a.Guess().Set(uint(clause[0]), cnf.Pos)
				lq.PushBack(cnf.Lit{uint(clause[0]), cnf.Pos})
			}
		} else {
			clauseDB.AddEntry(clause, true)
			n++
		}
	}
	// If we found unit clauses, we should BCP
	if foundUnit {
		if clauseDB.LQBcp(a.Guess(), lq, nil) == db.Conflict {
         fmt.Printf("conflict before dpll\n")
			// Unsatisfiable
			return -1, nil
		}
	}

	return n, nil
}

// Matches and returns the formula info.
func matchFormulaInfo(s string) (nVar int, nClauses int, err error) {
	// cnf
	r := scanner.NewScanner(strings.NewReader(s))
	c := r.Next()
	if c != "p" {
		return -1, -1, errors.New("First non-comment line does not begin with p.\n")
	}

	c = r.Next()
	if c != "cnf" {
		return -1, -1, errors.New("p line did not have cnf string\n")
	}

	// nVar
	nVar = r.NextInt()
	if nVar <= 0 {
		return -1, -1, errors.New("Invalid or incorrectly formatted nVar field")
	}

	// nClauses
	nClauses = r.NextInt()
	if nClauses <= 0 {
		return -1, -1, errors.New("Invalid or incorrectly formatted nClauses field")
	}

	return nVar, nClauses, nil
}

// Consumes the comment lines at the beginning of the dimacs Reader, if any.
func matchDimacsComments(r *scanner.Scanner) string {
	for {
		c := r.NextLine()
		if c[0] != 'c' {
			return c
		}
	}
	panic("Should never get here")
}
