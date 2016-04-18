package util

import (
	"net/http"
	"io/ioutil"
	"net/url"
	log "github.com/Sirupsen/logrus"
	_ "bytes"
	_ "strconv"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"bytes"
	"strconv"

)

func HttpGet(url string) (code int, jsonObj interface{}) {
	code, body := httpGet(url, map[string]string{})
	if body != "" {
		jsonObj = JsonParse(body)
	}
	return
}

func HttpPost(url string, data url.Values) (code int, jsonObj interface{}) {
	code, body := httpPost(url, data, map[string]string{})
	if body != "" {
		jsonObj = JsonParse(body)
	}
	return
}

func HttpGetExt(url string, additionalHeaders map[string]string) (code int, jsonObj interface{}) {
	code, body := httpGet(url, additionalHeaders)
	if body != "" {
		jsonObj = JsonParse(body)
	}
	return
}

func HttpPostExt(url string, data url.Values, additionalHeaders map[string]string) (code int, jsonObj interface{}) {
	code, body := httpPost(url, data, additionalHeaders)
	if body != "" {
		jsonObj = JsonParse(body)
	}
	return
}

func httpGet(url string, additionalHeaders map[string]string) (code int, body string) {
	log.Debugf("Sending GET %s", url)
	var client = http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("Could not create new request, reason %v", err)
	}
	for k,v := range additionalHeaders {
		req.Header.Add(k, v)
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Could not process GET request, reason %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Could not parse response, reason %v", err)
	}
	code, body = resp.StatusCode, string(bodyBytes)
	log.Debugf("Response is %d, %s", code, body)
	return
}

func httpPost(url string, data url.Values, additionalHeaders map[string]string) (code int, body string) {
	log.Debugf("Sending POST %s", url)
	var client = http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Panicf("Could not create new request, reason %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	for k,v := range additionalHeaders {
		req.Header.Add(k, v)
	}

	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Could not process POST request, reason %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Could not parse response, reason %v", err)
	}
	code, body = resp.StatusCode, string(bodyBytes)
	log.Debugf("Response is %d, %s", code, body)
	return
}

func FormValueReq(r *http.Request, key string, msgWhenError string) string {
	var res = r.FormValue(key)
	if res == "" {
		panic(CreateHttpError(http.StatusBadRequest, fmt.Sprintf("Argument %s is not given, %s", key, msgWhenError)))
	}
	return res
}

func FormValueOpt(r *http.Request, key string, defaultValue string) string {
	var res = r.FormValue(key)
	if res == "" {
		return defaultValue
	}
	return res
}

func QueryValueReq(r *http.Request, key string, msgWhenError string) string {
	values := r.URL.Query()
	res := values.Get(key)
	if res == "" {
		panic(CreateHttpError(http.StatusBadRequest, fmt.Sprintf("Argument %s is not given, %s", key, msgWhenError)))
	}
	return res
}

func QueryValueOpt(r *http.Request, key string, defaultValue string) string {
	values := r.URL.Query()
	res := values.Get(key)
	if res == "" {
		return defaultValue;
	}
	return res
}

func ParamByNameReq(ps *httprouter.Params, key string, msgWhenError string) string {
	var res = ps.ByName(key)
	if res == "" {
		panic(CreateHttpError(http.StatusBadRequest, fmt.Sprintf("Argument %s is not given, %s", key, msgWhenError)))
	}
	return res
}

func ParamByNameOpt(ps *httprouter.Params, key string, defaultValue string) string {
	var res = ps.ByName(key)
	if res == "" {
		return defaultValue
	}
	return res
}


