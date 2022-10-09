import Excel from 'exceljs';
import fs from 'fs';


async function set() {
    const writeStream = fs.createWriteStream("fix.sql", {flags: 'w+', encoding: 'utf-8', mode: 0o666});
    // 从文件读取
    const workbook = new Excel.Workbook();
    await workbook.xlsx.readFile("D:\\work\\dingding\\近期编辑过的商品0919(1).xlsx");

    const sheet = workbook.worksheets[0]
    const rows = sheet.findRows(2, sheet.rowCount)

    const data = fs.readFileSync("D:\\work\\商品id.txt", {encoding:'utf-8'});

    // split the contents by new line
   // const rows = data.split(/\r?\n/);

    //对数据进行处理
    rows.forEach(row=> {
        const productId = row.getCell("G").value;
        //const productId = row;
        writeStream.write(`UPDATE product_attr SET status = 0 WHERE product_id = ${productId} AND attr_id IN (SELECT attr_id FROM (SELECT attr_id FROM product_attr WHERE product_id = ${productId} AND updated_at > '2022-09-14' AND updated_at < '2022-09-19 16:41:00' GROUP BY attr_id HAVING COUNT(*) = 1) tmp) AND updated_at > '2022-09-14' AND updated_at < '2022-09-19 16:41:00' AND status = 1;
`)
   /*     const sortNumber = row.getCell("O").value;
        if (sortNumber && sortNumber !== '') {
            writeStream.write(`UPDATE product_info
                               SET sort_number = ${sortNumber}
                               WHERE id = ${productId};
            `)
        }*/
    })
    writeStream.end()
}


set().then()
