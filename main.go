package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

type Vehicle struct {
	vAct   bool
	vVin   string
	vYear  string
	vMake  string
	vModel string
}

func main() {
	baseURL := "https://www.qualitynissansc.com/used-inventory/index.htm"

	db, err := sql.Open("mysql", "bradTest:Newpassword1!@tcp(192.168.2.2:3306)/bradScrape_Developement")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	stmtIns, err := db.Prepare("INSERT INTO `Vehicles` (`VIN`, `YEAR`, `Make`, `Model`) VALUES (?,?,?,?)") //?? is placeholder for inserts
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	Vehicles := []Vehicle{}
	c := colly.NewCollector()

	c.OnHTML("li.inv-type-used", func(e *colly.HTMLElement) {
		temp := Vehicle{}
		temp.vVin = e.ChildAttr("div[data-type=used]", "data-vin")
		temp.vYear = e.ChildAttr("div[data-type=used]", "data-year")
		temp.vMake = e.ChildAttr("div[data-type=used]", "data-make")
		temp.vModel = e.ChildAttr("div[data-type=used]", "data-bodystyle")
		Vehicles = append(Vehicles, temp)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("div.pull-right", func(h *colly.HTMLElement) {
		t := h.ChildAttr("a[rel=next]", "data-href")
		nextString := baseURL + t
		c.Visit(nextString)
	})

	c.Visit(baseURL)

	rows, err := db.Query("SELECT * FROM Vehicles")
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	count := 0
	//get rows form sql db
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string

		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
				count++
			}
			fmt.Println("Found ", count, "in database", col[i], value)

		}
	}

	for i := 0; i < len(Vehicles); i++ {
		//fmt.Println(i+1, "/", len(Vehicles), "|\t", Vehicles[i].vVin)
		_, err := stmtIns.Exec(Vehicles[i].vVin, Vehicles[i].vYear, Vehicles[i].vMake, Vehicles[i].vModel)
		if err != nil {
			panic(err.Error())
		}
	}
}
