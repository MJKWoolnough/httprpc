package httprpc

import (
	"io"
	"net/http"
	"net/rpc"
)

type httpRPC struct {
	server      *rpc.Server
	serverCodec func(io.ReadWriteCloser) rpc.ServerCodec
	readLimit   int64
}

func Handle(server *rpc.Server, serverCodec func(io.ReadWriteCloser) rpc.ServerCodec, readLimit int64) http.Handler {
	if server == nil {
		server = rpc.DefaultServer
	}
	return httpRPC{
		server:      server,
		serverCodec: serverCodec,
		readLimit:   readLimit,
	}
}

type conn struct {
	io.Reader
	io.Closer
	io.Writer
}

func (h httpRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := conn{
		Reader: r.Body,
		Closer: r.Body,
		Writer: w,
	}
	if h.readLimit > 0 {
		c.Reader = io.LimitReader(r.Body, h.readLimit)
	}
	h.server.ServeCodec(&wrapper{
		ServerCodec: h.serverCodec(&c),
	})
}
