package wrap

import "context"

type EndPoint func(*Wrap) error

type Wrap struct {
	Methods   []EndPoint
	Fn        EndPoint
	index     int
	MethodTag string
	ctx       context.Context
}

func (w *Wrap) Next() (err error) {
	w.index++
	for s := len(w.Methods); w.index <= s+1; w.index++ {
		if w.index > s {
			return w.Fn(w)
		}
		err = w.Methods[w.index-1](w)
		if err != nil {
			return
		}
	}
	return
}

func (w *Wrap) Reset() {
	w.index = 0
}

func (w *Wrap) SetCtx(c context.Context) {
	w.ctx = c
}

func (w *Wrap) GetCtx() context.Context {
	return w.ctx
}
