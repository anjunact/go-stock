package models
import (
	_ "github.com/lib/pq"
	"github.com/go-ozzo/ozzo-dbx"
	//_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
)

func main() {
	type Stock struct {
		Id      int
		Name    string
		Code    string
		Price   float64
		Updated time.Time
	}
	db, err := dbx.Open("postgres", "host=192.168.56.1 user=stock dbname=stock password=stock sslmode=disable")
if err != nil{
	fmt.Println(err)
}
	// create a new query
	q := db.NewQuery("SELECT id, name,code,price,updated FROM stocks LIMIT {:id}")
	q.Bind(dbx.Params{"id": 3})
	// fetch all rows into a struct array
	var stocks []Stock
	q.All(&stocks)
fmt.Printf("%+v",stocks)
	// fetch a single row into a struct
	stock := Stock{}
	q.One(&stock)
	//fmt.Printf("%+v",stock)
	// fetch a single row into a string map
	data := dbx.NullStringMap{}
	q.One(data)
	fmt.Printf("%+v",data)
	// fetch row by row
	rows2, err := q.Rows()
	if err !=nil{
		fmt.Println(err)
	}
	for rows2.Next() {
		rows2.ScanStruct(&stock)
		// rows.ScanMap(data)
		// rows.Scan(&id, &name)
	}
}
