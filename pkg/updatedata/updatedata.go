package updatedata

import (
	"database/sql"
	"log"
)

func UpdateRows(rows [][]string, db *sql.DB) {
	stmt, err := db.Prepare("UPDATE `nse`.`nsedata` SET `col1` = ?, `col2` = ?, `col3` = ?, `col4` = ?, `col5` = ?, `col6` = ?, `col7` = ?, `col8` = ?, `col9` = ?, `col10` = ? WHERE `col0` = ?;")
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows[1:] {
		pkey := row[1] + "_" + row[2]
		_, err := stmt.Exec(row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], pkey)
		if err != nil {
			log.Fatal(err)
		}
	}
}
