package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/apex/gateway"
)

// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
var status_data = make(map[string]string)

var (
	port = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/get_api", getApiHandler)
	http.HandleFunc("/get_img", getImgHandler)
	http.HandleFunc("/readme", readmeHandler)
	http.HandleFunc("/send_api", sendApiHandler)
	listener := gateway.ListenAndServe
	portStr := "n/a"

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	log.Fatal(listener(portStr, nil))
}

func getApiHandler(w http.ResponseWriter, r *http.Request) {
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

func getImgHandler(w http.ResponseWriter, r *http.Request) {
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

func readmeHandler(w http.ResponseWriter, r *http.Request) {
	// 重定向页面到/
	http.Redirect(w, r, "/", http.StatusFound)
}

func sendApiHandler(w http.ResponseWriter, r *http.Request) {
	// 定义url、type、header、send_data四个变量
	var url, typ string
	var header, send_data map[string]interface{}
	// 判断请求方法
	if r.Method == "GET" {
			// 从get请求中获取url、type、header、data参数值
			url = r.URL.Query().Get("url")
			typ = r.URL.Query().Get("type")
			header = jsonToMap(r.URL.Query().Get("header"))
			send_data = jsonToMap(r.URL.Query().Get("data"))
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
			typ = data["type"]
			header = jsonToMap(data["header"])
			send_data = jsonToMap(data["data"])
	} else {
			// 如果请求方法不是get或post，返回错误信息
			status_data["status"] = "error"
			status_data["code"] = "1001"
			status_data["doit"] = "send_api"
			status_data["callback"] = fmt.Sprintf("INVALID_METHOD_%s", r.Method)
			w.Header().Set("Content-Type", "application/json")
			json_data, err := json.Marshal(status_data)
			if err != nil {
					fmt.Println(err)
					return
			}
			w.Write(json_data)
			return
	}

	// 将url、type、header、send_data赋值给status_data["doit"]
	m := map[string]interface{}{
		"url":     url,
		"type":    typ,
		"headers": header,
		"data":    send_data,
	}
  // 使用json.Marshal把map转换成一个JSON字符串
  s, err := json.Marshal(m)
  if err != nil {
      fmt.Println(err)
      return
  }
	status_data["doit"] = string(s)
	w.Header().Set("Content-Type", "application/json")
	json_data, err := json.Marshal(status_data)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(json_data)
	return
	

	// 如果url有值，根据type的方式向url发送请求，携带header和send_data，并返回响应数据
	if url != "" {
			var resp *http.Response
			var err error
			if typ == "POST" {
					// 如果type是post，使用http.Post方法发送请求，将send_data转换为json格式的字节切片作为body参数
					json_data, err := json.Marshal(send_data)
					if err != nil {
							fmt.Println(err)
							return
					}
					resp, err = http.Post(url, "application/json", ioutil.NopCloser(bytes.NewReader(json_data)))
					if err != nil {
							fmt.Println(err)
							return
					}
					defer resp.Body.Close()
					
					// 设置请求的header，遍历header变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
					for k, v := range header {
							resp.Header.Set(k, fmt.Sprint(v))
					}
					
			} else {
					// 如果type不是post，默认使用get方法发送请求，将send_data转换为查询字符串作为url参数
					query := make(map[string]string)
					for k, v := range send_data {
							query[k] = fmt.Sprint(v)
					}
					resp, err = http.Get(url + "?" + str_encode(query))
					if err != nil {
							fmt.Println(err)
							return
					}
					defer resp.Body.Close()
					
					// 设置请求的header，遍历header变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
					for k, v := range header {
							resp.Header.Set(k, fmt.Sprint(v))
					}
					
			}

			
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
	
	// 将status_data转换为json数据并返回
	w.Header().Set("Content-Type", "application/json")
	json_data, err := json.Marshal(status_data)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(json_data)
	} else {
			// 如果url没有值，返回错误信息
			status_data["status"] = "error"
			status_data["code"] = "1001"
			status_data["doit"] = "NO_URL"
			status_data["callback"] = "INVALID_HOOK"
			w.Header().Set("Content-Type", "application/json")
			json_data, err := json.Marshal(status_data)
			if err != nil {
					fmt.Println(err)
					return
			}
			w.Write(json_data)
	}
}

// 定义一个辅助函数，将json格式的字符串转换为map类型
func jsonToMap(s string) map[string]interface{} {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
			fmt.Println(err)
			return nil
	}
	return m
}

func str_encode(query map[string]string) string {
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