package main

import (
	"fmt"
	"popindata/popin"
	"popindata/server"
)

func main() {
	server.TimingMoney()
	fmt.Println("start ...")
	money := popin.Popin.GetAllMoneyPost("20200619")
	fmt.Println(money)
}
