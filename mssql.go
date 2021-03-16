package mssqlDriver

import (
	"database/sql"
	"io"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	//  "strconv"
)

// FileCSV - file csv
type FileCSV struct {
	filename string
	file     FileReader
	types    []*DataTypeWithOption

	delimiter   rune
	comment     rune
	withHeaders bool
	sampleSize  int

	withOrderColumn bool
}

// PingServer uses a passed database handle to check if the database server works
func (inst *FileCSV) Open() error {

	//connect to the database
	db, err := sql.Open("sqlserver", "sqlserver://[username]:[password]@[host]?database=[database]&connection+timeout=30")
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}

	defer Close()

}

// Close - close file
func (inst *FileCSV) Close() {
	db.Close()
}

func GetName(db *sql.DB) (string, error) {}

// GetSelector - returns selector
// What is to be selected
func GetSelector() *Selector {}

func GetData() (chan *DataRow, error) {
	var (
		id   int
		name string
	)
	rows, err := db.Query("querry for data", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&_)
		if err != nil {
			log.Fatal(err)
		}
		log.Println()
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

//Copied from csv_file to have an understanding of what data to get
func (inst *FileCSV) GetData() (chan *DataRow, error) {
	dataStream := make(chan *DataRow)

	types, err := inst.GetTypes()
	if err != nil {
		return dataStream, err
	}
	reader := inst.reader(true)
	if inst.withOrderColumn {
		types = types[1:]
	}

	go func(stream chan *DataRow, types []*DataTypeWithOption, withHeaders bool) {
		hasMetFirst := false
		orderCounter := 0
		hasError := false
		for {
			dr := &DataRow{}

			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				dr.Error = err
				stream <- dr
				hasError = true
				break
			}
			if len(record) != len(types) {
				dr.Error = ErrorIncorrectTypesNumber
				stream <- dr
				hasError = true
				break
			}
			if !hasMetFirst {
				hasMetFirst = true
				if inst.withHeaders {
					continue
				}
			}
			if inst.withOrderColumn {
				dr.Row = append(dr.Row, orderCounter+1)
			}
			for i, rv := range record {
				v, err := inst.applyColumnType(rv, types[i].Type(), types[i].Options()...)
				if err != nil {
					dr.Error = err
					stream <- dr
					hasError = true
					break
				}
				dr.Row = append(dr.Row, v)
			}
			if hasError {
				break
			}
			stream <- dr
			orderCounter++
		}
		close(dataStream)
	}(dataStream, types, inst.withHeaders)

	return dataStream, nil
}
