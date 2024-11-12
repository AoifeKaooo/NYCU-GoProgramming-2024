package main

import (
	"math/big"
	"syscall/js"
)

// CheckPrime checks if the input number is prime and displays the result in the answer element.
func CheckPrime(this js.Value, args []js.Value) interface{} {
	// 取得前端輸入的數字字串
	input := js.Global().Get("document").Call("getElementById", "value").Get("value").String()

	// 將輸入字串轉換為大數，以處理大範圍的整數
	num := new(big.Int)
	num.SetString(input, 10) // 10 表示十進位

	// 使用 ProbablyPrime 方法檢查是否為質數 (0 表示完全檢查)
	isPrime := num.ProbablyPrime(0)

	// 顯示結果
	if isPrime {
		js.Global().Get("document").Call("getElementById", "answer").Set("innerText", "It's prime")
	} else {
		js.Global().Get("document").Call("getElementById", "answer").Set("innerText", "It's not prime")
	}

	return nil
}

func registerCallbacks() {
	js.Global().Set("CheckPrime", js.FuncOf(CheckPrime))
}

func main() {
	registerCallbacks()
	select {}
}