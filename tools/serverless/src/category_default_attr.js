const Excel = require('exceljs');
const axios = require('axios');

const cardCategory = new Map()

async function init() {
    const res = await axios.get('/api/sc-attr/value/list?attrId=7')
    if (res.status !== 200 || res.data.status !== 0) {
        throw new Error("获取卡主品类出错")
    }
    const list = res.data.data.list
    for (let obj of list) {
        cardCategory.set( obj.attrValueDisplay,obj.attrValue)
    }
    console.log(cardCategory)
}

async function set() {
    // 从文件读取
    const workbook = new Excel.Workbook();
    await workbook.xlsx.readFile('F:/类目新增默认属性配置4.2.xlsx');

    const sheet = workbook.worksheets[0]
    const rows = sheet.findRows(2,sheet.rowCount)
    console.log(rows.length)
    //对数据进行处理
    for (let row of rows) {
        const param = {};
        param.categoryId = row.getCell("D").value
        param.cardCategory = cardCategory.get(row.getCell("F").value.toString())
        param.expressRule = "0"
        param.supportCard = row.getCell("G").value === "支持" ? "0" : "1"
        param.supportBalance = row.getCell("I").value === "支持" ? "0" : "1"
        param.supportCoupons = row.getCell("H").value === "支持" ? "0" : "1"
        param.scope = "1"
        param.serviceAttr = ((v) => {
            switch (v) {
                case "无理由退换货":
                    return "2"
                case "不可退换":
                    return "1"
                case "质量问题退换":
                    return "3"
                case "优先赔":
                    return "5"
            }
        })(row.getCell("K").value)
        param.refundAmount = row.getCell("J").value === "自动计算" ? "0" : "1"
        param.serviceDeadline = row.getCell("M").value.toString()
        param.logisticDeadline = row.getCell("N").value.toString()
        param.deliveryDay = row.getCell("O").value.toString()
        param.buyNote = row.getCell("P").value.toString()
        param.shieldSearch = row.getCell("Q").value === "不屏蔽" ? "0" : "1"
        param.shieldRecommend = row.getCell("R").value === "不推荐" ? "0" : "1"
        param.deliveryScopeId =[310000,530000,110000,220000,510000,120000,340000,370000,140000,440000,450000,320000,360000,130000,410000,330000,460000,420000,430000,350000,520000,210000,500000,610000,230000].map(v=>v.toString())
        console.log(param)
        const res = await axios.post('/api/erp/category/setDefaultAttr',param).catch(e=>console.log(e))

        console.log(res.data)
        await sleep(200)
    }

}

let sleep = function (delay) {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            try {
                resolve(1)
            } catch (e) {
                reject(0)
            }
        }, delay);
    })
}

axios.defaults.baseURL = 'http://dev-scm.sdyxmall.com/scm-api';
axios.defaults.headers.common["X-Mng-Token"] = '94bcaaf9-d082-48e2-94ca-207d618f6bd0';
init().then()
set().then()
