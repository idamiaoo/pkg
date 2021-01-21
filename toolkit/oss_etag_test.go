package toolkit

import "testing"

// 有效etag
func TestGetOSSObjectMD5(t *testing.T) {
	gotMd5, _ := GetOSSObjectMD5("http://wisroom-video.oss-cn-shenzhen.aliyuncs.com/1O8a3KwPcYPC0dbbC0gAb8E7tY9")
	wantMd5 := "68a688e65d912f9f4a35c53f3f99a314"

	if gotMd5 != wantMd5 {
		t.Errorf("got %q want %q", gotMd5, wantMd5)
	}
}

// 无效etag
func TestGetOSSObjectMD52(t *testing.T) {
	gotMd5, _ := GetOSSObjectMD5("https://wisroom-images-test.oss-cn-shenzhen.aliyuncs.com/K2B/Unit1/Session2/1QK1kk6rT6AkjZHL7XNuOr5lc6P.jpg")
	wantMd5 := "8f9acdd0364c8ca65c73331e870323cb"

	if gotMd5 != wantMd5 {
		t.Errorf("got %q want %q", gotMd5, wantMd5)
	}

}
