package scheduler

import (
	"database/sql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rishijain300900/nse-data-scrapper/pkg/adddata"
	"github.com/rishijain300900/nse-data-scrapper/pkg/readcsv"
	"github.com/rishijain300900/nse-data-scrapper/pkg/updatedata"
)

var (
	str, date string
	fileno    int
)

const (
	nseUrl     = "https://www1.nseindia.com/archives/nsccl/var/C_VAR1_"
	nseFormat  = ".DAT"
	connString = "root:jaibabaki123@tcp(localhost:3306)/nse"
)

func init() {
	str = time.Now().Format("01-02-2006")
	date = str[3:5] + str[:2] + str[6:]
	fileno = 1

}

//sat, sun will not run
//err handling on line 38
func PerformScheduling() {
	for fileno <= 6 {
		link := nseUrl + date + "_" + strconv.Itoa(fileno) + nseFormat
		response, _ := http.Get(link)
		if response.StatusCode == 200 {
			name := storeData(response.Body)
			data := readcsv.ReadCsv(name)
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
	file, err := os.Create(filepath.Join("Data", filepath.Base(name)))
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
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer db.Close()
	if fileno == 1 {
		adddata.ClearAndInsertRows(data, db)
		log.Println("Data added to Database")
	} else {
		updatedata.UpdateRows(data, db)
		log.Println("Database Updated")
	}
}
