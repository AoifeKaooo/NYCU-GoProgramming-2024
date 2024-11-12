package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	Expression string
	Result     string
	Error      string
}

// Calculator function to handle calculations
func Calculator(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	op := r.URL.Query().Get("op")
	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")

	num1, err1 := strconv.Atoi(num1Str)
	num2, err2 := strconv.Atoi(num2Str)

	if err1 != nil || err2 != nil || op == "" {
		http.ServeFile(w, r, "error.html")
		return
	}

	var data PageData
	tmpl := template.Must(template.ParseFiles("index.html"))
	switch op {
	case "add":
		data = PageData{Expression: fmt.Sprintf("%d + %d", num1, num2), Result: fmt.Sprintf("%d", num1+num2)}
	case "sub":
		data = PageData{Expression: fmt.Sprintf("%d - %d", num1, num2), Result: fmt.Sprintf("%d", num1-num2)}
	case "mul":
		data = PageData{Expression: fmt.Sprintf("%d * %d", num1, num2), Result: fmt.Sprintf("%d", num1*num2)}
	case "div":
		if num2 == 0 {
			http.ServeFile(w, r, "error.html")
			return
		}
		data = PageData{Expression: fmt.Sprintf("%d / %d", num1, num2), Result: fmt.Sprintf("%d", num1/num2)}
	case "gcd":
		data = PageData{Expression: fmt.Sprintf("GCD(%d, %d)", num1, num2), Result: fmt.Sprintf("%d", gcd(num1, num2))}
	case "lcm":
		data = PageData{Expression: fmt.Sprintf("LCM(%d, %d)", num1, num2), Result: fmt.Sprintf("%d", lcm(num1, num2))}
	default:
		http.ServeFile(w, r, "error.html")
		return
	}

	tmpl.Execute(w, data)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return (a * b) / gcd(a, b)
}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}

