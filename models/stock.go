package models

import (
	"database/sql"
	"time"

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

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=stock dbname=stock password=stock sslmode=disable")
	if err != nil {
		panic(err)
	}
}
func Stocks(page int, pageSize int) (stocks []Stock, err error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	fmt.Printf("%v,%v", offset, page)
	rows, err := Db.Query("select id, name,code,price,updated from stocks offset $1  limit $2", offset, pageSize)
	if err != nil {
		return
	}
	for rows.Next() {
		stock := Stock{}
		err = rows.Scan(&stock.Id, &stock.Name, &stock.Code, &stock.Price, &stock.Updated)
		if err != nil {
			return
		}
		stocks = append(stocks, stock)
	}
	rows.Close()
	return
}

func GetStock(code string) (stock Stock, err error) {
	stock = Stock{}
	err = Db.QueryRow("select id, name, code,price,updated from stocks where code = $1", code).Scan(&stock.Id, &stock.Name, &stock.Code, &stock.Price, &stock.Updated)
	return
}
func (stock *Stock) Create() (err error) {
	statement := "insert into stocks ( name,code,price) values ($1, $2,$3) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(stock.Name, stock.Code, stock.Price).Scan(&stock.Id)
	return
}

func (stock *Stock) Update() (err error) {
	_, err = Db.Exec("update stocks set price = $2,name=$3,updated=$4 where id = $1", stock.Id, stock.Price, stock.Name, time.Now())
	return
}
func (stock *Stock) Save() (err error) {
	s, _ := stock.Get(stock.Code)
	if s.Code == stock.Code {
		stock.Update()
	} else {
		stock.Create()
	}
	return
}
func (stock *Stock) Count() (count int) {
	statement := "select count(id) from stocks"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	stmt.QueryRow().Scan(&count)
	return
}
func (stock *Stock) Get(code string) (rs *Stock, err error) {
	statement := "select id, name,code,price,updated from stocks where code=$1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	stock1 := Stock{}
	stmt.QueryRow(code).Scan(&stock1.Id, &stock1.Name, &stock1.Code, &stock1.Price, &stock1.Updated)
	rs = &stock1
	return
}
func (stock *Stock) Delete() (err error) {
	_, err = Db.Exec("delete from stocks where id = $1", stock.Id)
	return
}

func DeleteAll() (err error) {
	_, err = Db.Exec("TRUNCATE TABLE stocks")
	return
}
