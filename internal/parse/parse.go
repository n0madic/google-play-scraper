package parse

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Float parse from string
func Float(str string) float64 {
	re := regexp.MustCompile(`\d+[.,]\d+`)
	f, err := strconv.ParseFloat(strings.Replace(re.FindString(str), ",", ".", 1), 64)
	if err != nil {
		return 0
	}
	return f
}

// ID parse return ID value from url
func ID(path string) (id string) {
	p := strings.Split(path, "?")
	if len(p) == 2 {
		m, err := url.ParseQuery(p[1])
		if err == nil {
			id = m.Get("id")
		}
	}
	return
}

// Int parse from string
func Int(str string) int {
	re := regexp.MustCompile(`[^0-9 ]+`)
	i, err := strconv.Atoi(re.ReplaceAllString(str, ""))
	if err != nil {
		return 0
	}
	return i
}

// Int64 parse from string
func Int64(str string) int64 {
	re := regexp.MustCompile(`[^0-9 ]+`)
	i64, err := strconv.ParseInt(re.ReplaceAllString(str, ""), 10, 64)
	if err != nil {
		return 0
	}
	return i64
}
