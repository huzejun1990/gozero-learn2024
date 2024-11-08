// @Author huzejun 2024/11/8 22:04:00
package biz

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/zeromicro/go-zero/core/iox"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/codec"
	"github.com/zeromicro/go-zero/core/fs"
)

const (
	pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCyeDYV2ieOtNDi6tuNtAbmUjN9
pTHluAU5yiKEz8826QohcxqUKP3hybZBcm60p+rUxMAJFBJ8Dt+UJ6sEMzrf1rOF
YOImVvORkXjpFU7sCJkhnLMs/kxtRzcZJG6ADUlG4GDCNcZpY/qELEvwgm2kCcHi
tGC2mO8opFFFHTR0aQIDAQAB
-----END PUBLIC KEY-----`
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCyeDYV2ieOtNDi6tuNtAbmUjN9pTHluAU5yiKEz8826QohcxqU
KP3hybZBcm60p+rUxMAJFBJ8Dt+UJ6sEMzrf1rOFYOImVvORkXjpFU7sCJkhnLMs
/kxtRzcZJG6ADUlG4GDCNcZpY/qELEvwgm2kCcHitGC2mO8opFFFHTR0aQIDAQAB
AoGAcENv+jT9VyZkk6karLuG75DbtPiaN5+XIfAF4Ld76FWVOs9V88cJVON20xpx
ixBphqexCMToj8MnXuHJEN5M9H15XXx/9IuiMm3FOw0i6o0+4V8XwHr47siT6T+r
HuZEyXER/2qrm0nxyC17TXtd/+TtpfQWSbivl6xcAEo9RRECQQDj6OR6AbMQAIDn
v+AhP/y7duDZimWJIuMwhigA1T2qDbtOoAEcjv3DB1dAswJ7clcnkxI9a6/0RDF9
0IEHUcX9AkEAyHdcegWiayEnbatxWcNWm1/5jFnCN+GTRRFrOhBCyFr2ZdjFV4T+
acGtG6omXWaZJy1GZz6pybOGy93NwLB93QJARKMJ0/iZDbOpHqI5hKn5mhd2Je25
IHDCTQXKHF4cAQ+7njUvwIMLx2V5kIGYuMa5mrB/KMI6rmyvHv3hLewhnQJBAMMb
cPUOENMllINnzk2oEd3tXiscnSvYL4aUeoErnGP2LERZ40/YD+mMZ9g6FVboaX04
0oHf+k5mnXZD7WJyJD0CQQDJ2HyFbNaUUHK+lcifCibfzKTgmnNh9ZpePFumgJzI
EfFE5H+nzsbbry2XgJbWzRNvuFTOLWn4zM+aFyy9WvbO
-----END RSA PRIVATE KEY-----`
	body = ""
)

var key = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")

func TestContentSecurity(t *testing.T) {
	tests := []struct {
		name        string
		mode        string
		extraKey    string
		extraSecret string
		extraTime   string
		err         error
		code        int
	}{
		{
			name:      "encrypted",
			mode:      "1",
			extraTime: "3600",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			limit := int64(1024)
			src := []byte("")
			rb := bytes.NewBuffer(src)
			wb := new(bytes.Buffer)
			r, err := http.NewRequest(http.MethodGet, "http://localhost:8888/v1/user/info",
				io.TeeReader(iox.LimitTeeReader(http.NoBody, wb, limit), rb))
			r.Header.Set("Content-Type", "application/json")
			assert.Nil(t, err)

			timestamp := time.Now().Unix()
			bodySign := computeBodySignature(r)
			fmt.Println(bodySign)
			contentOfSign := strings.Join([]string{
				strconv.FormatInt(timestamp+3600, 10),
				http.MethodGet,
				r.URL.Path,
				r.URL.RawQuery,
				bodySign,
			}, "\n")
			sign := hs256(key, contentOfSign)
			content := strings.Join([]string{
				"version=v1",
				"type=" + test.mode,
				fmt.Sprintf("key=%s", base64.StdEncoding.EncodeToString(key)) + test.extraKey,
				"time=" + strconv.FormatInt(timestamp+3600, 10),
			}, "; ")

			encrypter, err := codec.NewRsaEncrypter([]byte(pubKey))
			if err != nil {
				log.Fatal(err)
			}

			output, err := encrypter.Encrypt([]byte(content))
			if err != nil {
				log.Fatal(err)
			}

			encryptedContent := base64.StdEncoding.EncodeToString(output)
			join := strings.Join([]string{
				fmt.Sprintf("key=%s", fingerprint(pubKey)),
				"secret=" + encryptedContent + test.extraSecret,
				"signature=" + sign,
			}, "; ")
			fmt.Println(join)
			r.Header.Set("X-Content-Security", join)

			file, err := fs.TempFilenameWithText(priKey)
			assert.Nil(t, err)
			defer os.Remove(file)

		})
	}
}

func computeBodySignature(r *http.Request) string {
	var dup io.ReadCloser
	r.Body, dup = iox.DupReadCloser(r.Body)
	sha := sha256.New()
	io.Copy(sha, r.Body)
	r.Body = dup
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func fingerprint(key string) string {
	h := md5.New()
	io.WriteString(h, key)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func hs256(key []byte, body string) string {
	h := hmac.New(sha256.New, key)
	io.WriteString(h, body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
