package db

import (
   "cnf"
)

type Watch struct {
   Next *Watch
   Prev *Watch
   Watching cnf.Lit
}

func (w *Watch) New() {
   w.Next = w
   w.Prev = w
   w.Watching.Val = 0
   w.Watching.Pol = 0
}

func (w *Watch) isDummy() bool {
   return w.Watching.Val == 0
}
