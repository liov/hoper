package excel

import (
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func Excel() {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "test")
	f.SaveAs("./test.xlsx")
	os.Remove("./test.xlsx")
}
