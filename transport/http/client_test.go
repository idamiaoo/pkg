package http

import (
	"fmt"
	"net/http"
	"testing"
)

type testCallOption struct {
}

func (testCallOption) Before(r *http.Request) error {
	fmt.Println("before")
	return nil
}

func (testCallOption) After(r *http.Response) error {
	fmt.Println("after")
	return nil
}

func TestParseTarget(t *testing.T) {
	e1 := "www.ucloud.cn/a/b/c"
	u, err := parseTarget(e1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("scheme=%s, authority=%s, endpoint=%s", u.Scheme, u.Authority, u.Endpoint)
}
