package netutils

import (
	"bytes"
	"cloud-maque-common/utils"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// post提交参数
func GetWebRequestPostJson(url string, bytesData []byte) (result string, err error) {

	if len(url) < 10 {
		if err != nil {
			return "", utils.NewError(500, "请求地址不正确")
		}
	}
	postdata := bytes.NewReader(bytesData)
	resp, err := http.Post(url, "application/json;charset=UTF-8", postdata)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	result = string(body)

	return

}

// post提交参数
func GetWebRequestPost(url string) (result string, err error) {
	b := []byte("{}")

	result, err = GetWebRequestPostJson(url, b)

	return

}

// 以GET获取网络数据
func GetWebRequestGet(url string) (result string, err error) {

	connectTimeout := 5 * time.Second

	if len(url) < 10 {
		if err != nil {
			return "", utils.NewError(500, "请求地址不正确")
		}
	}

	if url[0:5] != "https" {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}
		result = string(body)

	} else {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr, Timeout: connectTimeout}
		resp, err := client.Get(url)

		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		result = string(body)
	}

	return

}

func GetWebRequestGetWithHeader(url string, header string, ip string) (result string, err error) {

	connectTimeout := 5 * time.Second

	if len(url) < 10 {
		if err != nil {
			return "", utils.NewError(500, "请求地址不正确")
		}
	}

	resp := &http.Response{}

	reqest, err := http.NewRequest("GET", url, nil)
	if header != "" {
		var param []interface{}
		err = json.Unmarshal([]byte(header), &param)

		if err != nil {
			return
		}

		if len(param) > 0 {
			for i := 0; i < len(param); i++ {
				var vv = param[i].(map[string]interface{})

				for key, value := range vv {
					reqest.Header.Add(key, value.(string))
				}
			}

		}

		//reqest.Header.Add("Cookie", "xxxxxx")
		//reqest.Header.Add("User-Agent", "xxx")
		//reqest.Header.Add("X-Requested-With", "xxxx")
	}

	if ip != "" {
		reqest.Header.Set("X-Real-Ip", ip)
		reqest.Header.Set("X-Forwarded-For", ip)
	}
	if url[0:5] != "https" {
		//提交请求
		resp, err = (&http.Client{}).Do(reqest)
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr, Timeout: connectTimeout}
		resp, err = client.Do(reqest)
	}

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result = string(body)

	return

}

func GetWebRequestPostWithHeader(url string, header string, data string) (result string, err error) {

	connectTimeout := 5 * time.Second

	if len(url) < 10 {
		if err != nil {
			return "", utils.NewError(500, "请求地址不正确")
		}
	}

	resp := &http.Response{}

	postdata := bytes.NewReader([]byte(data))
	reqest, err := http.NewRequest("POST", url, postdata)
	if header != "" {
		var param []interface{}
		err = json.Unmarshal([]byte(header), &param)

		if err != nil {
			return
		}

		if len(param) > 0 {
			for i := 0; i < len(param); i++ {
				var vv = param[i].(map[string]interface{})

				for key, value := range vv {
					reqest.Header.Add(key, value.(string))
				}
			}

		}

		//reqest.Header.Add("Cookie", "xxxxxx")
		//reqest.Header.Add("User-Agent", "xxx")
		//reqest.Header.Add("X-Requested-With", "xxxx")
	}

	if url[0:5] != "https" {
		//提交请求
		resp, err = (&http.Client{}).Do(reqest)
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr, Timeout: connectTimeout}
		resp, err = client.Do(reqest)
	}

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result = string(body)

	return

}
