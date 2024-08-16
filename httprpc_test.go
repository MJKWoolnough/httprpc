package httprpc

import (
	"bytes"
	"io"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"testing"
)

type RPCT struct{}

func (RPCT) TrimSpace(str string, s *string) error {
	*s = strings.TrimSpace(str)
	return nil
}

func TestRPC(t *testing.T) {
	rpc.Register(RPCT{})
	srv := httptest.NewServer(Handle(nil, jsonrpc.NewServerCodec, 0, "application/json; charset=utf-8"))
	var buf bytes.Buffer
	for n, test := range []struct {
		Input, Output string
	}{
		{
			"{\"method\":\"RPCT.TrimSpace\",\"params\":[\" 123 \"],\"id\":0}\n",
			"{\"id\":0,\"result\":\"123\",\"error\":null}\n",
		},
		{
			"{\"method\":\"RPCT.TrimSpace\",\"params\":[\" 123 \"],\"id\":0}\n{\"method\":\"RPCT.TrimSpace\",\"params\":[\" 123 \"],\"id\":0}",
			"{\"id\":0,\"result\":\"123\",\"error\":null}\n{\"id\":0,\"result\":\"123\",\"error\":null}\n",
		},
		{
			"{\"method\":\"RPCT.TrimSpace\",\"params\":[\" 123 ABC \"],\"id\":1}\n{\"method\":\"RPCT.TrimSpace\",\"params\":[\"\\tABC 123\\t\"],\"id\":2}\n",
			"{\"id\":2,\"result\":\"ABC 123\",\"error\":null}\n{\"id\":1,\"result\":\"123 ABC\",\"error\":null}\n",
		},
	} {
		resp, _ := srv.Client().Post(srv.URL+"/", "application/json", strings.NewReader(test.Input))
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		if str := buf.String(); str != test.Output {
			t.Errorf("test %d: expecting %q, got %q", n+1, test.Output, str)
		}
		buf.Reset()
	}
	srv.Close()
}
