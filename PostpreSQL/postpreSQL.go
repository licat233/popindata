package PostpreSQL

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
)

type SQL struct {
	Host string `json:"pg_host"`
	Port string `json:"pg_port"`
	User string `json:"pg_account"`
	Password string `json:"pg_password"`
	DBname string `json:"pg_dbname"`
	TBname string `json:"pg_tbname"`
}

var PG SQL
var db *sql.DB
var err error

func init()  {
	var file *os.File
	var configByte []byte
	if file, err = os.Open("config.json");err != nil {
		log.Fatal(err)
	}
	if configByte,err = ioutil.ReadAll(file);err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(configByte,&PG);err != nil {
		log.Fatal(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable", PG.Host, PG.Port, PG.User, PG.Password, PG.DBname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected!")
}