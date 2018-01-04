package httprpc

import (
	"net/rpc"
	"sync"
)

type wrapper struct {
	rpc.ServerCodec
	wg sync.WaitGroup
}

func (w *wrapper) ReadRequestHeader(r *rpc.Request) error {
	err := w.ServerCodec.ReadRequestHeader(r)
	if err == nil {
		w.wg.Add(1)
	}
	return err
}

func (w *wrapper) WriteResponse(r *rpc.Response, i interface{}) error {
	err := w.WriteResponse(r, i)
	w.wg.Done()
	return err
}

func (w *wrapper) Close() error {
	w.wg.Wait()
	return w.ServerCodec.Close()
}
