# httprpc

[![CI](https://github.com/MJKWoolnough/httprpc/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/httprpc/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/httprpc.svg)](https://pkg.go.dev/vimagination.zapto.org/httprpc)

--
    import "vimagination.zapto.org/httprpc"

Package httprpc creates an HTTP POST endpoint that wraps an RPC server.

## Highlights

 - Wrap a standard Go `rpc.Server` in an HTTP handler.
 - Can customise transfer codec, Content-Type header, and read limit size.

## Usage

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"

	"vimagination.zapto.org/httprpc"
)

type Math struct{}

func (Math) Add(args [2]int, sum *int) error {
	*sum = args[0] + args[1]

	return nil
}

type bufCloser struct {
	bytes.Buffer
}

func (bufCloser) Close() error { return nil }

type bodyWriter struct {
	io.ReadCloser
}

func (bodyWriter) Write(_ []byte) (int, error) { return 0, io.ErrNoProgress }

func main() {
	rs := rpc.NewServer()
	rs.Register(new(Math))

	var mux http.ServeMux

	mux.Handle("/", httprpc.Handle(rs, func(conn io.ReadWriteCloser) rpc.ServerCodec { return jsonrpc.NewServerCodec(conn) }, 1<<10, "application/json"))

	srv := httptest.NewServer(&mux)

	var buf bufCloser

	jsonrpc.NewClientCodec(&buf).WriteRequest(&rpc.Request{ServiceMethod: "Math.Add"}, [2]int{1, 2})

	resp, err := srv.Client().Post(srv.URL+"/", "application/json", &buf)
	if err != nil {
		fmt.Println(err)

		return
	}

	var sum int

	c := jsonrpc.NewClientCodec(&bodyWriter{resp.Body})
	c.ReadResponseHeader(&rpc.Response{})
	c.ReadResponseBody(&sum)

	fmt.Println(sum)

	// Output:
	// 3
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/httprpc
