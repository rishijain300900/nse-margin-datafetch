package scheduler

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	add_to_db "github.com/rishijain300900/nse-margin-datafetch/internal/db/add_to_db"
	update_db "github.com/rishijain300900/nse-margin-datafetch/internal/db/update_db"
	read "github.com/rishijain300900/nse-margin-datafetch/internal/read"
)

var (
	str, date, ConnString string
	fileno                int
)

const (
	nseUrl    = "https://www1.nseindia.com/archives/nsccl/var/C_VAR1_"
	nseFormat = ".DAT"
)

type jsonread struct {
	Username string `json:"username"`
	Paswword string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Server   string `json:"server"`
}

func init() {
	jsonread := jsonread{}
	file, err := ioutil.ReadFile("input.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = json.Unmarshal([]byte(file), &jsonread)
	if err != nil {
		log.Fatal(err)
		return
	}
	ConnString = jsonread.Username + ":" + jsonread.Paswword + "@tcp(" + jsonread.Ip + ":" + jsonread.Port + ")/" + jsonread.Server
	str = time.Now().Format("01-02-2006")
	date = str[3:5] + str[:2] + str[6:]
	fileno = 1

}

func PerformScheduling() {
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		log.Println("Its a weekend")
		return
	}
	for fileno <= 6 {
		link := nseUrl + date + "_" + strconv.Itoa(fileno) + nseFormat
		response, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode == 200 {
			name := storeData(response.Body)
			data := read.ReadCsv(name)
			sqladd(data)
			fileno++
		} else {
			time.Sleep(2 * time.Minute)
		}
	}
	log.Println("All files downloaded for today")
}

func storeData(resBody io.ReadCloser) string {
	name := " C_VAR1_" + date + "_" + strconv.Itoa(fileno) + ".csv"
	file, err := os.Create(filepath.Join("data", filepath.Base(name)))
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resBody)
	if err != nil {
		log.Fatal(err)
	}
	file.Write(bytes)
	log.Println(name, "file downloaded")

	return name
}

func sqladd(data [][]string) {
	db, err := sql.Open("mysql", ConnString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer db.Close()
	if fileno == 1 {
		add_to_db.ClearAndInsertRows(data, db)
		log.Println("Data added to Database")
	} else {
		update_db.UpdateRows(data, db)
		log.Println("Database Updated")
	}
}
