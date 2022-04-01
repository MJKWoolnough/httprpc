# httprpc
--
    import "vimagination.zapto.org/httprpc"

Package httprpc creates an HTTP POST endpoint that wraps an RPC server

## Usage

#### func  Handle

```go
func Handle(server *rpc.Server, serverCodec func(io.ReadWriteCloser) rpc.ServerCodec, readLimit int64, contentType string) http.Handler
```
Handle returns a new http.Handler that wraps an RPC server
