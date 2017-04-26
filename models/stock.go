package models

import (
	"time"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"fmt"
)

type Stock struct {
	Id      int
	Name    string
	Code    string
	Price   float64
	Updated time.Time
}

var db *dbx.DB

func init() {
	var err error
	db, err = dbx.Open("postgres", "host=192.168.56.1 user=stock dbname=stock password=stock sslmode=disable")
	if err != nil {
		panic(err)
	}
}
func Stocks(page int, pageSize int,params  map[string]interface{}) (stocks []Stock, err error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	fmt.Printf("%v,%v", offset, page)
	contitions := "where 1 "
	var bindParams dbx.Params = dbx.Params{"offset": offset,"pageSize":pageSize}
	if params["code"] !=nil{
		contitions += " and code=:code"
		bindParams["code"] = params["code"]
	}
	q := db.NewQuery("select id, name,code,price,updated from stocks "+contitions+" offset {:offset}  limit {:pageSize}")

	defer q.Close()
	q.Bind(bindParams)
	q.All(&stocks)
	return
}

func GetStock(code string) (stock Stock, err error) {
	stock = Stock{}
	q := db.NewQuery("select id, name, code,price,updated from stocks where code = {:code}", )
	q.Bind(dbx.Params{"code":code})
	q.One(&stock)
	return
}
func (stock *Stock) Create() (err error) {
	parms := dbx.Params{ "name":stock.Name,"code":stock.Code,"price":stock.Price }
	db.Insert("stocks",parms).Execute()
	return
}

func (stock *Stock) Update() (err error) {
	 db.Update("stocks",  dbx.Params{ "name":stock.Name,"code":stock.Code,"price":stock.Price },dbx.HashExp{"id":stock.Id}).Execute()
	return
}
func (stock *Stock) Save() (err error) {
	s := stock.Get(stock.Code)
	if s !=nil {
		stock.Update()
		fmt.Printf("db:%+v",stock)
	} else {
		stock.Create()
	}
	return
}
func (stock *Stock) Count() (count int) {
	q := db.NewQuery("select count(id) from stocks")
	defer  q.Close()
	q.Row(&count)
	return
}
func (stock *Stock) Get(code string) (rs *Stock) {
	q := db.NewQuery("select id, name,code,price,updated from stocks where code={:code}")
	defer  q.Close()
	q.Bind(dbx.Params{"code":code})
	q.One(&rs)
	return
}
func (stock *Stock) Delete() (err error) {
	db.Delete("stocks",dbx.HashExp{"id": stock.Id}).Execute()
	return
}

func DeleteAll() (err error) {
	//_, err = db.Exec("TRUNCATE TABLE stocks")
	return
}

