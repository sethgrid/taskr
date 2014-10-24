package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB_NAME string

var NEW_LABEL string
var LABEL string

var MESSAGE string

var SHOW bool

var DB *sql.DB

func init() {
	flag.StringVar(&DB_NAME, "db_name", "taskr.db", "defaults to taskr.db")
	flag.StringVar(&NEW_LABEL, "label", "", "create a new label to use for messages")
	flag.StringVar(&LABEL, "l", "default", "set a list of labels allowing you to group your message, comma delimited")
	flag.StringVar(&MESSAGE, "m", "", "a message entry to log")
	flag.BoolVar(&SHOW, "show", false, "set to true to show all entries")
}

func main() {
	flag.Parse()

	if !dbExists(DB_NAME) {
		log.Println("No DB Found. Creating")
		err := createDB(DB_NAME)
		if err != nil {
			log.Fatal(err)
		}
	}
	var err error

	DB, err = sql.Open("sqlite3", fmt.Sprintf("./%s", DB_NAME))
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if NEW_LABEL != "" {
		err := createLabel(NEW_LABEL)
		if err != nil {
			log.Fatal(err)
		}
	}

	labels := labelMapper(LABEL)

	if SHOW {
		show()
		return
	}

	err = insert(MESSAGE, labels)
	if err != nil {
		log.Fatal(err)
	}
}

func show() {
	r, err := DB.Query(`select id, created, message from entries`)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	for r.Next() {
		var messageID int
		var created int64
		var message string

		err := r.Scan(&messageID, &created, &message)
		if err != nil {
			log.Fatal(err)
		}

		r2, err := DB.Query(`select l.label from labels as l join entry_labels as el on l.id=el.id where el.entry=?`, messageID)
		if err != nil {
			log.Fatal(err)
		}
		labels := make([]string, 0)
		for r2.Next() {
			var label string
			err := r2.Scan(&label)
			if err != nil {
				log.Fatal(err)
			}
			labels = append(labels, label)
		}
		r2.Close()

		fmt.Printf("%s - %s %v\n", date(created), message, labels)
	}
}

func date(unix int64) string {
	t := time.Unix(unix, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%s %d %d", month.String(), day, year)

}

func insert(message string, labels map[string]int) error {
	if message == "" {
		return nil
	}

	r, err := DB.Exec(`insert into entries (message, created) values (?, ?)`, message, time.Now().Unix())
	if err != nil {
		return err
	}
	messageID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	for _, labelID := range labels {
		log.Println("inserting label: %d", labelID)
		_, err := DB.Exec(`insert into entry_labels (entry, label) values (?, ?)`, messageID, labelID)
		if err != nil {
			return err
		}
	}
	return nil
}

func createLabel(name string) error {
	_, err := DB.Exec(`insert into labels (label) values (?);`, name)
	return err
}

func labelMapper(list string) map[string]int {
	kv := make(map[string]int)
	labels := strings.Split(list, ",")

	for _, label := range labels {
		r, err := DB.Query("select id from labels where label=? limit 1", label)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		for r.Next() {
			var labelID int
			err = r.Scan(&labelID)
			if err != nil {
				log.Fatal(err)
			}
			kv[label] = labelID
		}
	}

	return kv
}

func dbExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func createDB(name string) error {
	err := os.Remove(fmt.Sprintf("./%s", name))
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s", name))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table entries (id integer not null primary key, message text, created integer);
	create table labels (id integer not null primary key, label text);
	create table entry_labels (id integer not null primary key, entry integer, label integer);

	insert into labels (label) values ("default");
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}
