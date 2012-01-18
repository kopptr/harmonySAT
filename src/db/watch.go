package db

import (
   "cnf"
)

type Watch struct {
   Next *Watch
   Prev *Watch
   Watching cnf.Lit
}

func NewWatch() (w *Watch) {
   w = new(Watch)
   w.Watching.Val = 0
   w.Watching.Pol = 0
   w.Next = w
   w.Prev = w
   return
}

func (w *Watch) IsDummy() bool {
   return w.Watching.Val == 0
}
