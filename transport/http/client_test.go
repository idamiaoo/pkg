package http

import (
	"fmt"
	"io"
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

func TestUxiao(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://uxiao.ucloudadmin.com/uqa/apisource?api=GetCertificateList", nil)
	req.AddCookie(&http.Cookie{
		Name: "INNER_AUTH_TOKEN", Value: "61xtNteyilCaGkYFNBLBviF9rPGQYWbPN2zz4bq9jMqsKMcxq296M1mhutbjD3bNNOE9msxjLaziRwX4dY5UiGxoWew=",
	})
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	t.Log(string(body))
}
