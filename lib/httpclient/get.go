package httpclient

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/eviltomorrow/rogue/lib/simplifiedchinese"
)

// DefaultHeader default header
var DefaultHeader = map[string]string{
	"Connection":                "keep-alive",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"Accept-Encoding":           "gzip, deflate",
	"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,da;q=0.7,pt;q=0.6,ja;q=0.5",
}

func Get(url string, timeout time.Duration, header map[string]string) (string, error) {
	var client = &http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new http request failure, nest error: %v", err)
	}

	for key, val := range header {
		request.Header.Add(key, val)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("client do http request failure, nest error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		buf, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("http status code: %d, message: %s", response.StatusCode, buf)
	}

	var buffer []byte
	contentEncode := response.Header.Get("Content-Encoding")
	switch {
	case contentEncode == "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return "", fmt.Errorf("panic: gzip new reader failure, nest error: %v", err)
		}
		defer reader.Close()

		buffer, err = ioutil.ReadAll(reader)
		if err != nil {
			return "", fmt.Errorf("panic: read all data buffer failure, nest error: %v", err)
		}

	default:
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("panic: read all data buffer failure, nest error: %v", err)
		}
		buffer = buf
	}

	var data string
	contentType := response.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, simplifiedchinese.GB18030.String()):
		data = simplifiedchinese.BytesToString(simplifiedchinese.GB18030, buffer)
	case strings.Contains(contentType, simplifiedchinese.GBK.String()):
		data = simplifiedchinese.BytesToString(simplifiedchinese.GBK, buffer)
	default:
		data = simplifiedchinese.BytesToString(simplifiedchinese.UTF8, buffer)
	}

	return data, nil
}
