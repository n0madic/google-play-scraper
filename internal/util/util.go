package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/k3a/html2text"
	"github.com/tidwall/gjson"
)

var (
	scriptRegex = regexp.MustCompile(`>AF_initDataCallback[\s\S]*?<\/script`)
	keyRegex    = regexp.MustCompile(`(ds:\d*?)'`)
	valueRegex  = regexp.MustCompile(`data:([\s\S]*?), sideChannel: {}}\);<\/`)
)

// AbsoluteURL return absolute url
func AbsoluteURL(base, path string) (string, error) {
	p, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	b, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return b.ResolveReference(p).String(), nil
}

// BatchExecute for PlayStoreUi
func BatchExecute(country, language, payload string) (string, error) {
	url := "https://play.google.com/_/PlayStoreUi/data/batchexecute"

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	q := req.URL.Query()
	q.Add("authuser", "0")
	q.Add("bl", "boq_playuiserver_20190424.04_p0")
	q.Add("gl", country)
	q.Add("hl", language)
	q.Add("soc-app", "121")
	q.Add("soc-platform", "1")
	q.Add("soc-device", "1")
	q.Add("rpcids", "qnKhOb")
	req.URL.RawQuery = q.Encode()

	body, err := DoRequest(req)
	if err != nil {
		return "", err
	}

	var js [][]interface{}
	err = json.Unmarshal(bytes.TrimLeft(body, ")]}'"), &js)
	if err != nil {
		return "", err
	}
	if len(js) < 1 || len(js[0]) < 2 {
		return "", fmt.Errorf("invalid size of the resulting array")
	}
	if js[0][2] == nil {
		return "", nil
	}

	return js[0][2].(string), nil
}

// ExtractInitData from Google HTML
func ExtractInitData(html []byte) map[string]string {
	data := make(map[string]string)
	scripts := scriptRegex.FindAll(html, -1)
	for _, script := range scripts {
		key := keyRegex.FindSubmatch(script)
		value := valueRegex.FindSubmatch(script)
		if len(key) > 1 && len(value) > 1 {
			data[string(key[1])] = string(value[1])
		}
	}
	return data
}

// GetJSONArray by path
func GetJSONArray(data string, paths ...string) []gjson.Result {
	for _, path := range paths {
		if gjson.Get(data, path).Exists() {
			return gjson.Get(data, path).Array()
		}
	}
	return nil
}

// GetJSONValue with multiple path
func GetJSONValue(data string, paths ...string) string {
	for _, path := range paths {
		if gjson.Get(data, path).Exists() {
			return gjson.Get(data, path).String()
		}
	}
	return ""
}

// DoRequest by HTTP and read all
func DoRequest(req *http.Request) ([]byte, error) {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

// GetInitData from Google HTML
func GetInitData(req *http.Request) (map[string]string, error) {
	html, err := DoRequest(req)
	if err != nil {
		return nil, err
	}

	return ExtractInitData(html), nil
}

// HTMLToText return plain text from HTML
func HTMLToText(html string) string {
	html2text.SetUnixLbr(true)
	return html2text.HTML2Text(html)
}
