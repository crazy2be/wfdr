package parsetime

import (
	//"fmt"
	"strings"
	"strconv"
	"time"
	"os"
)

func Parse(sDate string) (d time.Time, e os.Error) {
	//fmt.Printf("Parsing date: %s\n", sDate)
	// Splits into time and date parts
	var dateTime = strings.Split(sDate, "T", -1)
	// Splits off the timezone
	//var timeZone = strings.Split(dateTime[1], "Z", -1)
	// Parse the date
	var date = strings.Split(dateTime[0], "-", -1)
	d.Year, _ = strconv.Atoi64(date[0])
	d.Month, _ = strconv.Atoi(date[1])
	d.Day, _ = strconv.Atoi(date[2])
	// Parse the time
	var time = strings.Split(dateTime[1], ":", -1)
	d.Hour, _ = strconv.Atoi(time[0])
	d.Minute, _ = strconv.Atoi(time[1])
	d.Second, _ = strconv.Atoi(time[2])
	//fmt.Printf("Parsed date into %#v\n", d)
	return
}
