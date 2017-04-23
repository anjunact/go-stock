package main

import (
	"github.com/sclevine/agouti"
	"log"
	"strconv"
	"github.com/anjunact/go-stock/models"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"fmt"
	"time"
)

func main() {
	//models.DeleteAll()
	getHtml()
	//s,_ := stock.Get("600000")
	//s.Price = 11.22
	//s.Save()
}
func getHtml()  {
	//driver := agouti.PhantomJS()
	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()
	//page, err := driver.NewPage(agouti.Browser("phantomjs"))
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}
	if err := page.Navigate("http://finance.sina.com.cn/data/#stock"); err != nil {
		log.Fatalf("Failed to navigate:%v", err)
	}
	for true {
		html,_ := page.HTML()
		save(html)
		page.FindByClass("nextPageSpan").Click()

		currentPageSpan,err2:= page.FindByClass("currentPageSpan").Text()
		totalPageSpan,err3 := page.FindByClass("totalPageSpan").Text()
		if( err3 !=nil ){
			fmt.Println(err2)
			break
		}
		if(err3 != nil){
			fmt.Println(err3)
			break
		}
		if(currentPageSpan == totalPageSpan){
			fmt.Println(currentPageSpan)
			break
		}
	}
	return
}

func  save(html string)  {
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("#block_1 tbody tr").Each(func(i int, s *goquery.Selection) {
		var row []string
		s.Find("td").Each(func(j int, s1 *goquery.Selection){
			row = append(row,s1.Text())
		})
		var stock   models.Stock
		stock.Name = row[1]
		stock.Code = row[0][2:]
		price,err := strconv.ParseFloat(row[2],64)
		if err !=nil{
			fmt.Println(err.Error())
		}
		stock.Price = price
		stock.Updated = time.Now()
		stock.Save()
		fmt.Println(stock)
	})
}