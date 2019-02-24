package main

import (
	"encoding/json"
	"github.com/metakeule/fmtdate"
)

type Person struct {
	Name     string
	BirthDay fmtdate.TimeDate
}

func main() {
	bday, _ := fmtdate.NewTimeDate("YYYY-MM-DD", "2000-12-04")
	// do error handling
	paul := &Person{"Paul", bday}

	data, _ := json.Marshal(paul)
	// do error handling

	println(string(data))
}
