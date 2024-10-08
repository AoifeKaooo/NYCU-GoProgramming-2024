package main

import "fmt"

func main() {
	var n int64

	fmt.Print("Enter a number: ")
	fmt.Scanln(&n)

	result := Sum(n)
	fmt.Println(result)
}

func Sum(n int64) string {
	var sum int64 = 0
	var expression string

	for i := int64(1); i <= n; i++ {
		if i%7 != 0 { // 只累加不是 7 的倍數的數字
			if sum > 0 {
				expression += "+" // 如果不是第一個數字，就在前面加上 "+"
			}
			expression += fmt.Sprintf("%d", i) // 將數字加入運算式
			sum += i
		}
	}

	expression += fmt.Sprintf("=%d", sum) // 最後加上等號和總和
	return expression
}
