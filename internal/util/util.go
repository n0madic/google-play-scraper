package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/k3a/html2text"
	"github.com/tidwall/gjson"
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

// ExtractInitData from Google HTML
func ExtractInitData(html []byte) map[string]string {
	scriptRegex := regexp.MustCompile(`>AF_initDataCallback[\s\S]*?<\/script`)
	keyRegex := regexp.MustCompile(`(ds:\d*?)'`)
	valueRegex := regexp.MustCompile(`return ([\s\S]*?)\n}}\);<\/`)

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
func GetJSONArray(data, path string) []gjson.Result {
	return gjson.Get(data, path).Array()
}

// GetJSONValue by path
func GetJSONValue(data, path string) string {
	return gjson.Get(data, path).String()
}

// GetInitData from Google HTML
func GetInitData(req *http.Request) (map[string]string, error) {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search error: %s", resp.Status)
	}

	html, err := ioutil.ReadAll(resp.Body)
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
