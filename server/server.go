package server

import (
	"fmt"
	"popindata/mysqlDB"
	"popindata/popin"
	"time"
)

//t 更新时间，单位：秒
var t time.Duration = 300

//TimingMoney 每隔t秒更新一次所有活动的消耗
func TimingMoney() {
	for {
		//开始更新
		yestM := mysqlDB.Mysql.ReadYesterdayCharge()
		nowM := popin.Popin.GetAllMoney()
		charge := nowM - yestM
		mysqlDB.Mysql.ReplaceTodayCharge(nowM, charge)
		fmt.Printf("yestM:%d,nowM:%d", yestM, nowM)
		<-time.NewTimer(time.Second * t).C
	}
}
