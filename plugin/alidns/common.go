package alidns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ParamInfo struct {
	Key string
	Val string
}

type ParamList []ParamInfo

func (p ParamList) Len() int {
	return len(p)
}

func (p ParamList) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ParamList) Less(i int, j int) bool {
	return p[i].Key < p[j].Key
}

func specialUrlEncode(ctx string) string {
	ctx = url.QueryEscape(ctx)
	ctx = strings.Replace(ctx, "+", "%20", -1)
	ctx = strings.Replace(ctx, "*", "%2A", -1)
	ctx = strings.Replace(ctx, "%7E", "~", -1)
	return ctx
}

func sign(key string, ctx string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(ctx))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

const (
	API_VERSION       = "2015-01-09"
	SIGNATURE_VERSION = "1.0"
	SIGNATURE_METHOD  = "HMAC-SHA1"
)

func BuildParam(accessKey string, accessSecret string, method string, params []ParamInfo) string {
	vals := make([]ParamInfo, 0)
	for _, p := range params {
		if p.Key != "Signature" {
			vals = append(vals, p)
		}
	}

	ts := time.Now().UTC().Format("2006-01-02 15:04:05")

	vals = append(vals, ParamInfo{"SignatureMethod", SIGNATURE_METHOD})
	vals = append(vals, ParamInfo{"SignatureNonce", fmt.Sprintf("%d", time.Now().UnixNano())})
	vals = append(vals, ParamInfo{"AccessKeyId", accessKey})
	vals = append(vals, ParamInfo{"SignatureVersion", SIGNATURE_VERSION})
	vals = append(vals, ParamInfo{"Timestamp", ts + "Z"})
	vals = append(vals, ParamInfo{"Format", "JSON"})
	vals = append(vals, ParamInfo{"Version", API_VERSION})

	sort.Sort(ParamList(vals))

	ctx := ""
	for idx, p := range vals {
		if idx == 0 {
			ctx += fmt.Sprintf("%s=%s", specialUrlEncode(p.Key), specialUrlEncode(p.Val))
		} else {
			ctx += fmt.Sprintf("&%s=%s", specialUrlEncode(p.Key), specialUrlEncode(p.Val))
		}
	}

	signature := sign(accessSecret+"&", method+"&"+specialUrlEncode("/")+"&"+specialUrlEncode(ctx))
	signature = specialUrlEncode(signature)
	return ctx + "&Signature=" + signature
}
