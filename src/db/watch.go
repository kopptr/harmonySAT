package db

import (
   "cnf"
)

type Watch struct {
   next *Watch
   prev *Watch
   watching cnf.Lit
}

func (w *Watch) New() {
   w.next = w
   w.prev = w
   w.watching.Val = 0
   w.watching.Pol = 0
}

func (w *Watch) isDummy() bool {
   return w.watching.Val == 0
}
