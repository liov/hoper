const schedule = require('node-schedule');
const axios = require('axios').default;

const scheduleCron = () => {
    //每分钟的第30秒定时执行一次:
    schedule.scheduleJob('*/20 * 0,8,9,10,13,18,19,20,21,22,23 * * *', request);
}

let last_id = 0;

const at = new Map([
    ["000002204", "xxx"],
    ["000002190", "xxx"],
])

async function request() {
    const res = await axios.post(`http://inmail.miz.so:1234/grid/att/CheckInOutGrid/`,
        "page=3&rp=100", {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded; param=value',
                'Cookie': 'sessionidadms=93e5c9bb2f79a3af699ef9464511c504',
            }
        })
    console.log(res.data)
/*    const ding = {
        msgtype: "text",
        text: {
            content: ''
        },
        at: {
            atMobiles: [] as string[],
            isAtAll: false
        }
    }
    for (let obj of res.data.rows) {
        if ((obj.name as string).length == 2) {
            obj.name = obj.name + "    "
        }
        if (obj.id > last_id && obj.DeptName == "xxx中心") {
            last_id = obj.id
            ding.text.content = ding.text.content + obj.name + ` : ` + obj.checktime + "\n"
            let mobile = at.get(obj.badgenumber)
            if (mobile) ding.at.atMobiles = ding.at.atMobiles.concat(ding.at.atMobiles, mobile)
        }
    }
    if (ding.text.content == "") return
    axios.post(`https://oapi.dingtalk.com/robot/send?access_token=xxx`,ding,{
        headers:{
            'Content-Type': 'application/json',
        }
    })
    console.log("请求钉钉")*/
}

request()