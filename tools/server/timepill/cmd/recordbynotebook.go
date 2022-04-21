package main

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	fmt.Println(timepill.Token)
	timepill.RecordByOrderNoteBook()
	//timepill.RecordByNoteBook(873)
}
