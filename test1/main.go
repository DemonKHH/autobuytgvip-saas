package main

import (
	"fmt"
	"strconv"
)

func main() {
	usdtAmount := 15.000000
	amount, err := strconv.ParseFloat(fmt.Sprintf("%.2f", usdtAmount), 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(amount)
}
