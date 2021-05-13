package export

import (
	"log"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tealeg/xlsx/v3"
)

// 360 excel库文档宣称性能很强，实际性能很差，尤其合并单元格
func NewFile(sheet string, header []string) *excelize.File {
	endColumn := ColumnLetter[len(header)-1]
	f := excelize.NewFile()
	//单元格样式
	style, err := f.NewStyle(Style)
	f.SetColStyle(sheet, "A:"+endColumn, style)
	headerStyle, _ := f.NewStyle(HeaderStyle)
	f.SetCellStyle(sheet, "A1", endColumn+"1", headerStyle)
	if err != nil {
		log.Println(err)
	}
	for i := range header {
		/*		axis, _ := excelize.CoordinatesToCellName(i+1, 1)
				f.SetCellValue(sheet, axis, header[i])*/
		f.SetCellValue(sheet, ColumnLetter[i]+"1", header[i])
	}
	f.SetRowHeight(sheet, 1, 30)
	return f
}

func NewSheet(f *excelize.File, sheet string, header []string) {
	endColumn := ColumnLetter[len(header)-1]
	//单元格样式
	style, _ := f.NewStyle(Style)
	headerStyle, _ := f.NewStyle(HeaderStyle)
	f.NewSheet(sheet)
	f.SetColStyle(sheet, "A:"+endColumn, style)
	f.SetCellStyle(sheet, "A1", endColumn+"1", headerStyle)

	for i := range header {
		f.SetCellValue(sheet, ColumnLetter[i]+"1", header[i])
	}
	f.SetRowHeight(sheet, 1, 30)
}

func NewXlsxFile(sheetName string, header []string) (*xlsx.File, *xlsx.Sheet) {
	f := xlsx.NewFile()
	sheet, _ := f.AddSheet(sheetName)
	row := sheet.AddRow()
	for _, v := range header {
		cell := row.AddCell()
		cell.SetString(v)
	}
	return f, sheet
}

func XlsxWriteValue(row *xlsx.Row, values []interface{}) {
	for _, v := range values {
		cell := row.AddCell()
		cell.SetValue(v)
	}
}
