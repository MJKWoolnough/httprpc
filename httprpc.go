// Package httprpc creates an HTTP endpoint that wraps an RPC server
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
	contentType string
}

// Handle returns a new http.Handler that wraps an RPC server
func Handle(server *rpc.Server, serverCodec func(io.ReadWriteCloser) rpc.ServerCodec, readLimit int64, contentType string) http.Handler {
	if server == nil {
		server = rpc.DefaultServer
	}
	return httpRPC{
		server:      server,
		serverCodec: serverCodec,
		readLimit:   readLimit,
		contentType: contentType,
	}
}

type conn struct {
	io.Reader
	io.Closer
	io.Writer
}

func (h httpRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.contentType != "" {
		w.Header().Set("Content-Type", h.contentType)
	}
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
