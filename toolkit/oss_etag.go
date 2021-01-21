package toolkit

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 获取oss资源的md5值
func GetOSSObjectMD5(urlStr string) (string, error){
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Head(urlStr)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 404 {
		return "", errors.New("资源不存在于oss上")
	}

	etag := resp.Header.Get("ETag")
	if etag == "" {
		return "", errors.New("无法获取oss资源的md5")
	}


	if invalidMD5 := strings.Contains(etag, "-"); invalidMD5 == false  {
		return strings.ToLower(strings.ReplaceAll(etag, `"`, "")), nil
	}

	h := md5.New()

	resp, _ = client.Get(urlStr)
	defer resp.Body.Close()

	if _, err = io.Copy(h, resp.Body); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), err

}