package mysqlDB

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"popindata/function"
	"strconv"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

type SQL struct {
	Addr string `json:"mysql_addr"`
	Account string `json:"mysql_account"`
	Password string `json:"mysql_password"`
	DBname string `json:"mysql_dbname"`
	TBname string `json:"mysql_tbname"`
	Conn *sql.DB
}

var Mysql SQL
var err error

func init(){
	var file *os.File
	var configByte []byte
	if file, err = os.Open("config.json");err != nil {
		log.Fatal(err)
	}
	if configByte,err = ioutil.ReadAll(file);err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(configByte,&Mysql);err != nil {
		log.Fatal(err)
	}
	Mysql.Conn, err = sql.Open("mysql", Mysql.Account+":"+Mysql.Password+"@tcp("+Mysql.Addr+")/"+Mysql.DBname)
	if err != nil {
		panic(err.Error())
	}

	stmtIns,err := Mysql.Conn.Query("CREATE TABLE IF NOT EXISTS `"+Mysql.DBname+"`.`"+Mysql.TBname+"` ( `id` INT NOT NULL AUTO_INCREMENT ,`total_charge` INT(64) NOT NULL , `today_charge` INT(64) NOT NULL, `date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP , PRIMARY KEY (`id`)) ENGINE = InnoDB;")
	if err !=nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
}

//InsertTodayCharge 插入当天消耗
func (t *SQL)InsertTodayCharge(totalCharge,todayCharge int){
	sqlml := "INSERT INTO `"+t.TBname+"` (`id`, `total_charge`, `today_charge`) VALUES ('"+time.Now().Format("20060102")+"', '"+strconv.Itoa(totalCharge)+"', '"+strconv.Itoa(todayCharge)+"');"
	//fmt.Println(sqlml)
	stmtOut, err := t.Conn.Query(sqlml)
	if err != nil {
		fmt.Println(err)
	}
	defer stmtOut.Close()
}

//UpdataTodayCharge 更新当天消耗
func (t *SQL)UpdataTodayCharge(totalCharge,todayCharge int){
	sqlml := "UPDATE `"+t.TBname+"` SET `total_charge` = '"+strconv.Itoa(totalCharge)+"', `today_charge` = '"+strconv.Itoa(todayCharge)+"' WHERE `id` = '"+time.Now().Format("20060102")+"';"
	//fmt.Println(sqlml)
	stmtOut, err := t.Conn.Query(sqlml)
	if err != nil {
		fmt.Println(err)
	}
	defer stmtOut.Close()
}

//ReplaceTodayCharge 更新当天数据，不存在则创建
func (t *SQL)ReplaceTodayCharge(totalCharge,todayCharge int){
	datestr := function.GetInsertDate()
	sqlml := "REPLACE INTO "+t.TBname+" (id,total_charge,today_charge) VALUES('"+datestr+"','"+strconv.Itoa(totalCharge)+"', '"+strconv.Itoa(todayCharge)+"');"
	fmt.Println(sqlml)
	stmtOut, err := t.Conn.Query(sqlml)
	if err != nil {
		fmt.Println(err)
	}
	defer stmtOut.Close()
}

//ReadYesterdayCharge 获取前一天的消耗
func (t *SQL)ReadYesterdayCharge() int {
	sqlml := "SELECT total_charge FROM `"+t.TBname+"` WHERE `id` = '"+function.GetYesterdaydate()+"' LIMIT 1"
	//fmt.Println(sqlml)
	rows, err := t.Conn.Query(sqlml)
	if err != nil {
		fmt.Println(err)
	}

	//defer outres.Close()
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var value string
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// Now do something with the data.
		// Here we just print each column as a string.

		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "0"
			} else {
				value = string(col)
			}
			//fmt.Println(columns[i], ": ", value)
		}

	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	//fmt.Println(value)
	charge,_:=strconv.Atoi(value)
	return charge
}