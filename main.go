package main

import (
	"fmt"
	//"time"

	"github.com/gocolly/colly"
)

type vehicleDescription struct {
	vdMilage    string
	vdMPG       string
	vdEngine    string
	vdTrans     string
	vdDriveLine string
	vdExterior  string
	vdInterior  string
	vdStock     string
}

type vehicle struct {
	vVin   string
	vYear  string
	vMake  string
	vModel string
	vTrim  string
	vStock string
	vDesc  vehicleDescription
	//vCrawl time.Time
}

func main() {
	vehicles := []vehicle{}
	c := colly.NewCollector()

	baseURL := "https://www.qualitynissansc.com/used-inventory/index.htm"

	c.OnHTML("li.inv-type-used", func(e *colly.HTMLElement) {
		temp := vehicle{}
		temp.vVin = e.ChildAttr("div[data-type=used]", "data-vin")
		temp.vYear = e.ChildAttr("div[data-type=used]", "data-year")
		temp.vMake = e.ChildAttr("div[data-type=used]", "data-make")
		temp.vModel = e.ChildAttr("div[data-type=used]", "data-bodystyle")
		temp.vTrim = e.ChildAttr("div[data-type=used]", "data-trim")
		temp.vDesc.vdExterior = e.ChildText("Exterior Color")
		//temp.vCrawl = time.Now()
		vehicles = append(vehicles, temp)
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

	for i := 0; i < len(vehicles); i++ {
		fmt.Println(i+1, "/", len(vehicles), "|\t", vehicles[i].vVin, vehicles[i].vDesc.vdExterior)
	}
}
