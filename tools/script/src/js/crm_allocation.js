import Excel from "exceljs";
import axios from "axios";

async function customerNums() {
    // 从文件读取
    const workbook = new Excel.Workbook();
    await workbook.xlsx.readFile('E:/xxx.xlsx');

    const sheet = workbook.worksheets[0]
    const rows = sheet.findRows(2,sheet.rowCount-2)
    console.log(rows.map(row=>`'${row.getCell("C").value}'`).join())
}

await customerNums()