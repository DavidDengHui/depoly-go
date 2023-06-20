package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)


func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
	var status_data = make(map[string]interface{})
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

	// 打印 HTTP 状态码
	fmt.Println("Status:", resp.Status)

	// 打印 HTTP 头部
	fmt.Println("Headers:", resp.Header)

	// 读取 HTTP 主体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// 打印 HTTP 主体的长度
	fmt.Println("Body length:", len(body))

	return

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