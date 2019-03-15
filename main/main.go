package main

import (
"database/sql"
"encoding/xml"
"fmt"
_ "github.com/go-sql-driver/mysql"
"log"
"net/http"
)

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	ID            int64
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Guid          string    `xml:"guid"`
	Enclosure     Enclosure `xml:"enclosure"`
	PublishedDate string    `xml:"pubDate"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

var database *sql.DB

func main() {
	resp, err := http.Get("https://news.tut.by/rss/sport/football.rss")
	if err != nil {
		log.Fatal("main.Error with http.Get: ", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal("main.Error with resp.Body.Close: ", err)
		}
	}()
	rss := Rss{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)

	database, err = newMySQLDb() //TODO
	_ = createTable(database)    //TODO
	/*stm, err := database.Prepare(dropTable);
	if err != nil {
		panic(err)
	}
	defer stm.Close()
	_, err = stm.Exec()
	if err != nil {
		panic(err)
	}*/
	//http.HandleFunc("/create", CreateHandler)
	fmt.Printf("Title: %v\n", rss.Channel.Title)
	fmt.Printf("Description: %v\n", rss.Channel.Desc)
	for i, item := range rss.Channel.Items {
		stmt, err := database.Prepare(insertIntoDatabase)
		log.Println(err)
		stmt.Exec(item.Title, item.PublishedDate)
		fmt.Printf("%d.\t Title: %v\n\t PubDate: %v\n", i+1, item.Title, item.PublishedDate)
		//i+1, item.Title, item.PublishedDate
	}
	http.HandleFunc("/", IndexHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":3006", http.DefaultServeMux)
}