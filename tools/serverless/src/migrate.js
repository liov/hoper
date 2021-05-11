import Excel from 'exceljs'
import axios from'axios'
import { sleep } from "../dist/utils/sleep.js"



async function set() {
    // 从文件读取
    const workbook = new Excel.Workbook();
    await workbook.xlsx.readFile('F:/xxx.xlsx');

    const sheet = workbook.worksheets[0]
    const rows = sheet.findRows(3,sheet.rowCount-2)
    console.log(rows.length)
    //对数据进行处理
    for (let row of rows) {
        const param = {};
        param.productId = parseInt(row.getCell("A").value)
        param.categoryId = row.getCell("B").value
        param.newType = 10
        console.log(param)
        const res = await axios.post('/v1',param).catch(e=>console.log(e))
        console.log(res.data)
        await sleep(200)
    }

}

async function defer() {
    const res = await axios.get('/v2')
    if (res.status !== 200 || res.data.status !== 0) {
        throw new Error("补充草稿出错")
    }
}



axios.defaults.baseURL = 'http://';
axios.defaults.headers.common["Token"] = '';

//set().then()
defer().then()