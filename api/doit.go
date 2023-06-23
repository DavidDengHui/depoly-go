package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
var status_data = make(map[string]interface{})

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	token := r.URL.Query().Get("token")
	if token == "" {
		token = "Z2hwX01PZDNidlo3aGRJbUpUTDJzQWR3V1VXNDRxZnBDdDBVWUhpNw=="
	}
	hookName := r.URL.Query().Get("hook_name")
	if r.Method == "POST" {
		var temp map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&temp)
		if err != nil {
			log.Println(err)
			return
		}
		token = ToList(temp)["password"]
		hookName = ToList(temp)["hook_name"]
	}
	if r.Header.Get("HTTP_X_GITHUB_EVENT") != "" {
		hookName = r.Header.Get("HTTP_X_GITHUB_EVENT")
	}
	if token != "" {
		if token == "get_info" {
			if hookName != "" {
				if hookName == "bilibili" {
					status_data["status"] = "success"
					status_data["code"] = "1101"
					status_data["doit"] = "https://api2.hnest.eu.org/get_api?url=https://api.bilibili.com/x/web-interface/popular?ps=1&pn=1"
					status_data["callback"] = hookName
				} else {
					status_data["code"] = "1004"
					status_data["doit"] = token + " | " + hookName
					status_data["callback"] = "INVALID_HOOK"
				}
			} else {
				status_data["code"] = "1003"
				status_data["doit"] = token
				status_data["callback"] = "INVALID_HOOK"
			}
		} else {
			tokenBytes, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				log.Println(err)
				status_data["code"] = "1002"
				status_data["doit"] = token
			} else {
				token = string(tokenBytes)
				if hookName != "" {
					if hookName == "push_hooks" {
						status_data["status"] = "success"
						status_data["code"] = "1104"
						status_data["doit"] = token + " | " + hookName
						username := r.URL.Query().Get("username")
						repopath := r.URL.Query().Get("repopath")
						reponame := r.URL.Query().Get("reponame")
						url := fmt.Sprintf("https://api.github.com/repos/%s/%s/dispatches", repopath, reponame)
						data := map[string]string{"event_type": r.URL.Query().Get("event_type")}
						dataJSON, err := json.Marshal(data)
						if err != nil {
							log.Println(err)
							return
						}
						headers := map[string]interface{}{
							"Accept":        "application/json",
							"Authorization": fmt.Sprintf("token %s", token),
							"User-Agent":    username,
							"Content-Type":  "application/json",
							"Content-Length": fmt.Sprintf("%d", len(dataJSON)),
						}
						// 创建一个自定义的请求对象，设置请求方法为post，请求的url为url，请求的数据为dataJSON
						req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(dataJSON)))
						if err != nil {
								fmt.Println(err)
								return
						}
						// 设置请求的标头，遍历headers变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
						for k, v := range headers {
								req.Header.Set(k, fmt.Sprint(v))
						}
						// 创建一个http.Client对象
						client := &http.Client{}
						// 使用http.Client.Do方法发送请求，并获取响应
						resp, err := client.Do(req)
						if err != nil {
								fmt.Println(err)
								return
						}
						defer resp.Body.Close()
						status_data["callback"] = map[string]interface{}{
							"url":     url,
							"data":    data,
							"headers": headers,
							"state":   resp.StatusCode,
						}
					} else {
						status_data["status"] = "error"
						status_data["code"] = "1008"
						status_data["doit"] = token + " | " + hookName
						status_data["callback"] = "INVALID_HOOK"
					}
				} else {
					status_data["code"] = "1007"
					status_data["doit"] = token
					status_data["callback"] = "NO_HOOK"
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json_data, err := json.Marshal(status_data)
	if err != nil {
			fmt.Println(err)
			return
	}
	w.Write(json_data)
}

func ToList(cont map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range cont {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}
