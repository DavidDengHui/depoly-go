package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
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

	http.HandleFunc("/api/bili", bili)
	http.HandleFunc("/get_api", getApiHandler)
	http.HandleFunc("/get_img", getImgHandler)
	listener := gateway.ListenAndServe
	portStr := "n/a"

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	log.Fatal(listener(portStr, nil))
}

func bili(w http.ResponseWriter, r *http.Request) {
	url := "https://api.bilibili.com/x/web-interface/popular?ps=1&pn=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(body)
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



	// status_data["callback"] = fmt.Sprintf("url:[%d] | type:[%d] | filename:[%d]", url, typ, filename)
	// w.Header().Set("Content-Type", "application/json")
	// json_data, err := json.Marshal(status_data)
	// if err != nil {
	// 		fmt.Println(err)
	// 		return
	// }
	// w.Write(json_data)
	// return

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