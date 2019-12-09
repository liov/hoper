package main

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "test")
	f.SaveAs("./test.xlsx")
}
