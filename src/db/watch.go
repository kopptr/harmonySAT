package db

type watch struct {
   next *watch
   prev *watch
   watching Lit
}

func (w *watch) New() {
   w.next = w
   w.prev = w
   w.watching.Val = 0
   w.watching.Pol = Unassigned
}

func (w *watch) isDummy() bool {
   return w.watching.Val == 0
}
