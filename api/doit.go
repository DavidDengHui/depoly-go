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
	var PayLoad map[string]interface{}
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&PayLoad)
		if err != nil {
			fmt.Println(err)
			return
		}
		if vaule, ok := PayLoad["password"]; ok {
			token = vaule.(string)
		}
		if vaule, ok := PayLoad["hook_name"]; ok {
			hookName = vaule.(string)
		}
	}
	if r.Header.Get("HTTP_X_GITHUB_EVENT") != "" {
		hookName = r.Header.Get("HTTP_X_GITHUB_EVENT")
	}
	if token != "" {
		if token == "get_info" {
			if hookName != "" {
				if hookName == "bilibili" {
					api := "api.bilibili.com"
					path := "/x/web-interface/"
					typ := "popular"
					if r.URL.Query().Get("type") != "" {
						switch r.URL.Query().Get("type") {
						case "rank":
							typ = "ranking/v2?rid=0&type=all"
						case "rank01":
							path = "/pgc/web/rank/"
							typ = "list?day=3&season_type=1"
						case "rank02":
							path = "/pgc/season/rank/web/"
							typ = "list?day=3&season_type=4"
						case "rank03":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=168&type=all"
						case "rank04":
							path = "/pgc/season/rank/web/"
							typ = "list?day=3&season_type=3"
						case "rank05":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=1&type=all"
						case "rank06":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=3&type=all"
						case "rank07":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=129&type=all"
						case "rank08":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=4&type=all"
						case "rank09":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=36&type=all"
						case "rank10":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=188&type=all"
						case "rank11":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=234&type=all"
						case "rank12":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=223&type=all"
						case "rank13":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=160&type=all"
						case "rank14":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=211&type=all"
						case "rank15":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=217&type=all"
						case "rank16":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=119&type=all"
						case "rank17":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=155&type=all"
						case "rank18":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=5&type=all"
						case "rank19":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=181&type=all"
						case "rank20":
							path = "/pgc/season/rank/web/"
							typ = "list?day=3&season_type=2"
						case "rank21":
							path = "/pgc/season/rank/web/"
							typ = "list?day=3&season_type=5"
						case "rank22":
							path = "/pgc/season/rank/web/"
							typ = "list?day=3&season_type=7"
						case "rank23":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=0&type=origin"
						case "rank24":
							path = "/x/web-interface/ranking/"
							typ = "v2?rid=0&type=rookie"
							break
						default:
							path = "/x/web-interface/"
							typ = "popular"
							break
						}
					} else {
						if r.URL.Query().Get("ps") != "" {
								typ += "?ps=" + r.URL.Query().Get("ps")
								if r.URL.Query().Get("pn") != "" {
										typ += "&pn=" + r.URL.Query().Get("pn")
								} else {
										typ += "&pn=1"
								}
						}
					}
					url := fmt.Sprintf("https://%s%s%s", api, path, typ)
					status_data["doit"] = url
					resp, err := http.Get(url)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						get_data, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							fmt.Println(err)
							return
						}
						data_json := make(map[string]interface{})
						err = json.Unmarshal(get_data, &data_json)
						if err != nil {
							fmt.Println(err)
							return
						}
						status_data["callback"] = data_json
						status_data["code"] = "1101"
						status_data["status"] = "success"
					} else {
						status_data["callback"] = fmt.Sprintf("%s[%d]", hookName, resp.StatusCode)
						status_data["code"] = "1010"
						status_data["status"] = "error"
					}
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
					if hookName == "deployment_status" {
						status_data["code"] = "1009"
						status_data["doit"] = fmt.Sprintf("%s | %s", token, hookName)
						username := r.URL.Query().Get("username")
						repopath := r.URL.Query().Get("repopath")
						reponame := r.URL.Query().Get("reponame")
						state := r.URL.Query().Get("state")
						page_url := r.URL.Query().Get("url")
						if r.Method == "POST" {
							if vaule, ok := PayLoad["deployment_status"].(map[string]interface{})["state"]; ok {
								state = vaule.(string)
							}
							if vaule, ok := PayLoad["deployment_status"].(map[string]interface{})["environment_url"]; ok {
								page_url = vaule.(string)
							}
						}
						if state == "success" {
							status_data["status"] = "success"
							status_data["code"] = "1103"
							url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/master", repopath, reponame)
							response, err := http.Get(url)
							if err != nil {
								fmt.Println(err)
								return
							}
							defer response.Body.Close()
							var data map[string]interface{}
							err = json.NewDecoder(response.Body).Decode(&data)
							if err != nil {
								fmt.Println(err)
								return
							}
							sha := data["sha"].(string)
							commitMsg := data["commit"].(map[string]interface{})["message"].(string)
							url = fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s/comments", repopath, reponame, sha)
							headers := map[string]interface{}{
								"Content-Type":    "application/json",
								"Authorization": fmt.Sprintf("token %s", token),
								"User-Agent":    username,
							}
							SendData := map[string]string{
								"body" : fmt.Sprintf("# Successfully deployed with \n > %s\n## Following the Pages URL:\n### [%s](%s)", commitMsg, page_url, page_url),
							}
							dataJSON, err := json.Marshal(SendData)
							if err != nil {
								log.Println(err)
								return
							}
							req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(dataJSON)))
							if err != nil {
									fmt.Println(err)
									return
							}
							for k, v := range headers {
									req.Header.Set(k, fmt.Sprint(v))
							}
							client := &http.Client{}
							resp, err := client.Do(req)
							if err != nil {
									fmt.Println(err)
									return
							}
							defer resp.Body.Close()
							var callback map[string]interface{}
							err = json.NewDecoder(resp.Body).Decode(&callback)
							if err != nil {
								fmt.Println(err)
								return
							}
							status_data["callback"] = callback
						} else {
							status_data["callback"] = fmt.Sprintf("[state]%s", state)
						}
					} else if hookName == "push_hooks" {
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
						req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(dataJSON)))
						if err != nil {
								fmt.Println(err)
								return
						}
						for k, v := range headers {
								req.Header.Set(k, fmt.Sprint(v))
						}
						client := &http.Client{}
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