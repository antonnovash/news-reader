package mySql

import (
"bytes"
"database/sql"
"fmt"
"html/template"
"log"
"net/http"
)

var createTableStatements = `CREATE TABLE IF NOT EXISTS news (
		id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		title TEXT NULL,
		publishedDate VARCHAR(255) NULL
	)`
var insertIntoDatabase = `INSERT INTO news (title, publishedDate) VALUES (?, ?)`

var dropTable = `DROP TABLE news;`

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func newMySQLDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:1234567@/news")
	if err != nil {
		log.Fatal("news.go.Error with open database : ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db, nil
}

// createTable creates the table
func createTable(db *sql.DB) error {
	//for _, stmt := range createTableStatements {
	_, err := db.Exec(createTableStatements)
	if err != nil {
		log.Fatal("Error with create table: %v", err)
	}
	//}
	return nil
}

/*func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		id := r.FormValue("id")
		title := r.FormValue("title")
		publishedDate := r.FormValue("publishedDate")

		_, err = database.Exec(`INSERT INTO news (id, title, publishedDate) VALUES ('id', 'title', 'publishDate')`, id, title, publishedDate)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		//http.ServeFile(w, r, "templates/create.html")
	}
}*/
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("SELECT * FROM news")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	items := []Item{}

	for rows.Next() {
		p := Item{}
		err := rows.Scan(&p.ID, &p.Title, &p.PublishedDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}

	tmpl, err := template.ParseFiles("template/index.html")
	tmpl.Execute(w, items)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var b bytes.Buffer
	tmpl.Execute(&b, items)
	log.Println(b.String())
}

/*func scanNews(s rowScanner) (*Item, error) {
	var (
		id            int64
		title         sql.NullString
		publishedDate sql.NullString
		description   sql.NullString
	)
	if err := s.Scan(&id, &title, &publishedDate, &description); err != nil {
		return nil, err
	}
	news := &Item{
		ID:            id,
		Title:         title.String,
		PublishedDate: publishedDate.String,
		Description:   description.String,
	}
	return news, nil
}*/
