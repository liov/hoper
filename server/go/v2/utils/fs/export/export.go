package export

import (
	"log"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tealeg/xlsx/v3"
)

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
		cell.Value = v
	}
	return f, sheet
}
