package dimacs

import (
	"dpll/assignment"
	"dpll/db"
	"dpll/db/cnf"
	"errors"
	"io"
	"strings"
	"scanner"
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
	_, err = matchClauses(s, clauseDB, a)
	if err != nil {
		return nil, nil, err
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
        n = 0
	for r.HasNextLine() {
		clause := []int{}
		for i := r.NextInt(); i != 0; i = r.NextInt() {
			clause = append(clause, i)
		}
		// Add unit clauses directly to assignment
		if len(clause) == 1 {
			if clause[0] < 1 {
				a.Guess().Set(uint((clause[0] * -1)), cnf.Neg)
			} else {
				a.Guess().Set(uint(clause[0]), cnf.Pos)
			}
		} else {
			clauseDB.AddEntry(clause, true)
                        n++
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
