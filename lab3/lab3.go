package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", Calculator)
	http.ListenAndServe(":8080", nil)
}

func Calculator(w http.ResponseWriter, r *http.Request) {
	// 拿到URL的路径并移除前导的'/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 分割路径
	parts := strings.Split(path, "/")

	// 检查路径至少应该有三个部分（操作符和两个操作数）
	if len(parts) < 3 {
		fmt.Fprint(w, "Error!")
		return
	}

	// 拿到操作符和操作数
	op := parts[0]
	nums := parts[1:]

	// 确保有两个操作数
	if len(nums) != 2 {
		fmt.Fprint(w, "Error!")
		return
	}

	// 解析操作数
	a, err1 := strconv.Atoi(nums[0])
	b, err2 := strconv.Atoi(nums[1])

	// 检查操作数解析是否成功
	if err1 != nil || err2 != nil {
		fmt.Fprint(w, "Error!")
		return
	}

	// 执行相应的操作
	switch op {
	case "add":
		fmt.Fprintf(w, "%d + %d = %d", a, b, a+b)
	case "sub":
		fmt.Fprintf(w, "%d - %d = %d", a, b, a-b)
	case "mul":
		fmt.Fprintf(w, "%d * %d = %d", a, b, a*b)
	case "div":
		// 检查除法是否除以0
		if b == 0 {
			fmt.Fprint(w, "Error!")
		} else {
			fmt.Fprintf(w, "%d / %d = %d, reminder = %d", a, b, a/b, a%b)
		}
	default:
		fmt.Fprint(w, "Error!")
	}
}

