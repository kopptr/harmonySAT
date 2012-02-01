package dimacs

import (
	"dpll/db"
	"errors"
	"fmt"
	"io"
	"scanner"
)

func DimacsToDb(r io.Reader) (clauseDB *db.DB, nVar int, err error) {
	s := scanner.NewScanner(r)

	err = matchDimacsComments(s)
	if err != nil {
		return nil, 0, nil
	}
	nVar, nClauses, err := matchFormulaInfo(s)
	if err != nil {
		return nil, 0, err
	}
	clauseDB = db.NewDB(nVar)
	n, err := matchClauses(s, clauseDB)
	if err != nil {
		return nil, 0, err
	} else if n != nClauses {
		return nil, 0, errors.New(fmt.Sprintf("Read %d/%d clauses.", n, nClauses))
	}

	return
}

// Matches all of the clauses in the database, and inserts them.
func matchClauses(r *scanner.Scanner, clauseDB *db.DB) (n int, err error) {
	for n = 0; r.HasNextLine(); n++ {
		clause := []int{}
		for i := r.NextInt(); i != 0; i = r.NextInt() {
			clause = append(clause, i)
		}
		clauseDB.AddEntry(clause, true)
	}
	return n, nil
}

// Matches and returns the formula info.
func matchFormulaInfo(r *scanner.Scanner) (nVar int, nClauses int, err error) {
	// cnf
	c := r.Next()
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
func matchDimacsComments(r *scanner.Scanner) error {
	for {
		c := r.Next()
		if c != "c" {
			if c != "p" {
				return errors.New("First non-comment line does not begine with a p\n")
			} else {
				return nil
			}
		}
		_ = r.NextLine()
	}
	panic("Should never get here")
}
