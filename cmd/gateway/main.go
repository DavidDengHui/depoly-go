package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/apex/gateway"
)

// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
var status_data = make(map[string]interface{})

var (
	port = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/doit", DoitHandler)
	http.HandleFunc("/get_api", GetApiHandler)
	http.HandleFunc("/get_img", GetImgHandler)
	http.HandleFunc("/readme", ReadmeHandler)
	http.HandleFunc("/send_api", SendApiHandler)
	http.HandleFunc("/get_web", GetWebHandler)

	listener := gateway.ListenAndServe
	portStr := "n/a"

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	log.Fatal(listener(portStr, nil))
}

func ToList(cont map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range cont {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

func DoitHandler(w http.ResponseWriter, r *http.Request) {
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
							fmt.Fprintf(w, fmt.Sprintf("%s|%s[%s]->[%s]:[%s]", token, hookName, state, page_url, PayLoad))
							return
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

func GetApiHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	// 定义一个url变量
	var url string
	// 先获取get请求中设置的url参数赋值给url
	url = r.URL.Query().Get("url")
	// 再判断如果是post方式的请求
	if r.Method == "POST" {
		// 则判断请求的body数据中有无url参数
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 如果有则将url参数值覆盖给url
		if data["url"] != "" {
			url = data["url"]
		}
	}
	// 最后再判断变量url是否有值
	if url != "" {
		// 增加将url值赋给status_data["doit"]
		status_data["doit"] = url
		// 有则向该url发送get请求
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		// 如果请求成功了直接输出返回的数据
		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(body)
		} else {
			// 否则输出json数据{"status":"error","code":"1002","doit":"获取到的url值","callback":"INVALID_URL_加请求状态"}
			w.Header().Set("Content-Type", "application/json")
			status_data["code"] = "1002"
			status_data["callback"] = fmt.Sprintf("INVALID_URL_%d", resp.StatusCode)
			json_data, err := json.Marshal(status_data)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(json_data)
		}
	} else {
		// 判断变量url没有值时，status_data["doit"]="NO_URL"、status_data["callback"]="INVALID_HOOK"
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
		// 然后将数组status_data转换为json数据后输出
		w.Header().Set("Content-Type", "application/json")
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Write(json_data)
    }
}

func GetImgHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	// 设置url、type、filename三个字符串变量
	var url, typ, filename string
	// 从get请求中获取url、type参数值
	url = r.URL.Query().Get("url")
	typ = r.URL.Query().Get("type")
	// filename初始化为"get_img"
	filename = "get_img"
	// 判断请求方式如果是POST，则将请求body部分中设置的url和type参数值覆盖
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			return
		}
		if data["url"] != "" {
			url = data["url"]
		}
		if data["type"] != "" {
			typ = data["type"]
		}
	}
	// 如果没有获取到url值，则status_data["doit"]="NO_URL"、status_data["callback"]="INVALID_HOOK"，将status_data转换为json输出
	if url == "" {
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
		w.Header().Set("Content-Type", "application/json")
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Write(json_data)
        return
    }
    // 如果获取到url值，则将url中最后一个/字符之后的内容赋值给filename
    filename = url[strings.LastIndex(url, "/")+1:]
    // 再判断type是否有值，没有则将type赋值为"jpg"
    if typ == "" {
        typ = "jpg"
        // 再继续判断url结尾如果有.字符，则将从.之后开始到非字母结束的部分赋值给type，将filename中.字符及以后的部分删除
        if strings.Contains(filename, ".") {
            re := regexp.MustCompile(`\.[a-zA-Z]+`)
            match := re.FindString(filename)
            if match != "" {
                typ = match[1:]
                // filename = strings.TrimSuffix(filename, match)
                filename = filename[:strings.Index(filename, ".")]
            }
        }
    }
    // 然后向url发送get请求，将返回的图片显示到屏幕上，mimetype为image/type，设置ctrl+s保存图片时名字为filename.type
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    w.Header().Set("Content-Type", fmt.Sprintf("image/%s", typ))

    // 设置ctrl+s保存图片时名字为filename.type
    w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s.%s", filename, typ))

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    w.Write(body)
}

func ReadmeHandler(w http.ResponseWriter, r *http.Request) {
	// 重定向页面到/
	http.Redirect(w, r, "/", http.StatusFound)
}

func SendApiHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	// 定义url、type、header、send_data四个变量
	var url, typ string
	var header, send_data map[string]interface{}
	// 判断请求方法
	if r.Method == "GET" {
		// 从get请求中获取url、type、header、data参数值
		url = r.URL.Query().Get("url")
		typ = strings.ToUpper(r.URL.Query().Get("type"))
		header = JsonToMap(r.URL.Query().Get("header"))
		send_data = JsonToMap(r.URL.Query().Get("data"))
	} else if r.Method == "POST" {
		// 从post请求的body中获取url、type、header、data参数值
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			return
		}
		url = data["url"]
		typ = strings.ToUpper(data["type"])
		header = JsonToMap(data["header"])
		send_data = JsonToMap(data["data"])
	}

	// 判断如果没有获取到url值，则返回错误信息并结束程序
	if url == "" {
		status_data["code"] = "1001"
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(json_data)
        return
    }

    // 将status_data["doit"]赋值为一个多维数组，其中包含"url"键值为url值，"type"键值为type值，"headers"键值为header数组，"data"键值为send_data数组
    status_data["doit"] = map[string]interface{}{
        "url": url,
        "type": typ,
        "headers": header,
        "data": send_data,
    }

    // 根据type的方式向url发送请求，携带header和send_data，并返回响应数据
    var resp *http.Response
    var err error
    if typ == "POST" {
        // 如果type是post，使用http.Post方法发送请求，将send_data转换为json格式的字节切片作为body参数
        json_data, err := json.Marshal(send_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        // return
        // 创建一个自定义的请求对象，设置请求方法为post，请求的url为url，请求的数据为json_data
        req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(json_data)))
        if err != nil {
            fmt.Println(err)
            return
        }
        // 设置请求的标头，遍历header变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
        for k, v := range header {
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
        // 读取响应的状态码和数据，并输出到页面上
        status := resp.StatusCode
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(err)
            return
        }
        data := string(body)
        // 设置status_data["status"]为"success"，status_data["code"]为"1102-响应状态码"
        status_data["status"] = "success"
        status_data["code"] = fmt.Sprintf("1102-%d", status)
        // 读取响应的body数据，并将其赋值给status_data["callback"]
        status_data["callback"] = data
    } else {
        // 如果type不是post，默认使用get方法发送请求，将send_data转换为查询字符串作为url参数
        query := make(map[string]string)
        for k, v := range send_data {
            query[k] = fmt.Sprint(v)
        }
        // 判断url中是否包含?字符，有则使用&连接，否则使用?连接
        if strings.Contains(url, "?") {
            resp, err = http.Get(url + "&" + StrEncode(query))
        } else {
            resp, err = http.Get(url + "?" + StrEncode(query))
        }
        if err != nil {
            fmt.Println(err)
            return
        }
        defer resp.Body.Close()
        // 读取响应的body数据，并将其赋值给status_data["callback"]
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(err)
            return
        }
        status_data["callback"] = string(body)

        // 设置status_data["status"]为"success"，status_data["code"]为"1101-响应状态码"
        status_data["status"] = "success"
        status_data["code"] = fmt.Sprintf("1101-%d", resp.StatusCode)
    }
    // 将status_data转换为json数据并返回
    json_data, err := json.Marshal(status_data)
    if err != nil {
        fmt.Println(err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(json_data)
    return
}

// 定义一个辅助函数，将json格式的字符串转换为map类型
func JsonToMap(s string) map[string]interface{} {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return m
}

func StrEncode(query map[string]string) string {
	// 定义一个字符串切片，用于存储编码后的键值对
	var pairs []string
	// 遍历query变量中的键值对
	for k, v := range query {
		// 使用url.QueryEscape函数对键和值进行编码，替换特殊字符为%XX序列
		k = url.QueryEscape(k)
		v = url.QueryEscape(v)
		// 将编码后的键和值用=连接起来，形成一个键值对字符串
		pair := k + "=" + v
		// 将键值对字符串追加到字符串切片中
		pairs = append(pairs, pair)
	}
	// 使用strings.Join函数将字符串切片中的元素用&连接起来，形成一个查询字符串
	queryStr := strings.Join(pairs, "&")
	// 返回查询字符串
	return queryStr
}

func GetWebHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	// 设置url、type、filename三个字符串变量
	var url string
	// 从get请求中获取url、type参数值
	url = r.URL.Query().Get("url")
	// 判断请求方式如果是POST，则将请求body部分中设置的url和type参数值覆盖
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			return
		}
		if data["url"] != "" {
			url = data["url"]
		}
	}
	// 如果没有获取到url值，则status_data["doit"]="NO_URL"、status_data["callback"]="INVALID_HOOK"，将status_data转换为json输出
	if url == "" {
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
		w.Header().Set("Content-Type", "application/json")
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Write(json_data)
        return
    }

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 将返回的数据直接写入到页面
		io.Copy(w, resp.Body)
}