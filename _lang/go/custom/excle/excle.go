package main

import (
	"log"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type Foo struct {
	A, B, C string
	D       []string
}

var column = [...]string{"A", "B", "C", "D"}

func main() {
	datas := make([]Foo, 1000)
	for i := range datas {
		datas[i].A = "A"
		datas[i].B = "B"
		datas[i].C = "C"
		datas[i].D = make([]string, 3)
		for j := range datas[i].D {
			datas[i].D[j] = "D"
			datas[i].D[j] = "E"
			datas[i].D[j] = "F"
		}
	}
	inOneForLoop(datas)
	inTwoForLoop(datas)
}

var style = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
}

func inOneForLoop(datas []Foo) {
	defer func(t time.Time) { log.Println("inOneForLoop", time.Now().Sub(t)) }(time.Now())
	f := excelize.NewFile()
	style, err := f.NewStyle(&style)
	f.SetColStyle("Sheet1", "A:D", style)
	if err != nil {
		log.Println(err)
	}
	var headerRow = 0
	var headerRowStr string
	for _, data := range datas {
		for _, sub := range data.D {
			headerRow++
			f.SetRowHeight("Sheet1", headerRow, 30)
			headerRowStr = strconv.Itoa(headerRow)
			f.SetCellValue("Sheet1", column[3]+headerRowStr, sub)
		}
		//合并单元格
		if len(data.D) != 1 {
			for j := 0; j < 3; j++ {
				f.MergeCell("Sheet1", column[j]+strconv.Itoa(headerRow-len(data.D)+1), column[j]+headerRowStr)
			}
		}
		currentRow := headerRow - len(data.D) + 1
		currentRowStr := strconv.Itoa(currentRow)
		f.SetCellValue("Sheet1", column[0]+currentRowStr, data.A)
		f.SetCellValue("Sheet1", column[1]+currentRowStr, data.B)
		f.SetCellValue("Sheet1", column[2]+currentRowStr, data.C)
	}
	f.SaveAs("./test.xlsx")
}

func inTwoForLoop(datas []Foo) {
	defer func(t time.Time) { log.Println("inTwoForLoop", time.Now().Sub(t)) }(time.Now())
	f := excelize.NewFile()
	style, err := f.NewStyle(&style)
	f.SetColStyle("Sheet1", "A:D", style)
	if err != nil {
		log.Println(err)
	}
	var headerRow = 0
	var headerRowStr string
	for _, data := range datas {
		for _, sub := range data.D {
			headerRow++
			f.SetRowHeight("Sheet1", headerRow, 30)
			headerRowStr = strconv.Itoa(headerRow)
			f.SetCellValue("Sheet1", column[3]+headerRowStr, sub)
		}
		currentRow := headerRow - len(data.D) + 1
		currentRowStr := strconv.Itoa(currentRow)
		f.SetCellValue("Sheet1", column[0]+currentRowStr, data.A)
		f.SetCellValue("Sheet1", column[1]+currentRowStr, data.B)
		f.SetCellValue("Sheet1", column[2]+currentRowStr, data.C)
	}
	headerRow = 1
	for _, data := range datas {
		headerRowStr = strconv.Itoa(headerRow)
		mergerRow := strconv.Itoa(headerRow + len(data.D) - 1)
		if len(data.D) != 1 {
			for j := 0; j < 3; j++ {
				f.MergeCell("Sheet1", column[j]+mergerRow, column[j]+headerRowStr)
			}
		}
		headerRow += len(data.D)
	}
	f.SaveAs("./test1.xlsx")
}
