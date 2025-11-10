package util

import (
	"io"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/Yoak3n/gulu/logger"
)

func ClientWithProxy(proxy string) *http.Client {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	parsed, err := url.Parse(proxy)
	if proxy != "" {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(parsed),
		}
	} else if err != nil {
		logger.Logger.Errorf("Parsing URL error: %v", err)
	}

	return client
}
func GetRequestWithCookie(addr string) *http.Request {
	uri, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	// TODO
	cookie := ""
	header := http.Header{
		"Cookie": []string{cookie},
		"User-Agent": []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		},
	}
	err = Sign(uri)
	if err != nil {
		return nil
	}
	req := &http.Request{
		Method: "GET",
		URL:    uri,
		Header: header,
	}
	return req
}

func RequestGetWithAll(addr string, proxy string) []byte {
	client := ClientWithProxy(proxy)
	req := GetRequestWithCookie(addr)
	if req == nil {
		return nil
	}
	uri := req.URL
	logger.Logger.Debugln(uri)
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorf("Request GetWithAll Error And Retrying: %s", err.Error())
		return requestRetry(req, 1)
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 412 {
			logger.Logger.Errorf("目前已被风控")
		}
		return requestRetry(req, 1)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return body
}

func requestRetry(req *http.Request, count int) []byte {
	if count == 10 {
		return nil
	}
	// TODO
	proxy := ""
	time.Sleep(time.Duration(math.Min(300.0, math.Pow(2, float64(count)))) * time.Second)
	client := ClientWithProxy(proxy)
	res, err := client.Do(req)
	if err != nil {
		return requestRetry(req, count+1)
	}
	if res.StatusCode != 200 {
		return requestRetry(req, count+1)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return body
}
