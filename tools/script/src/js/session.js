

function erp() {
    const sess = {
        type: 1,
        platformType: 101,
        thirdCompId: 0,
        compId: 10001,
        filterCompIds: [10001],
        employeeId: 2006,
        employeeName: 'sale2',
        piId: "1234567"
    };
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'));
}



function crm() {
    const sess = {
        userId: 200,
        userName: '',
        userRealName: '',
        clientIp: '',
    };
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'));
}

function escape() {
    const token = "eyJjb21wSWQiOjEsImNvbXBOYW1lIjoi5rex5Zyz5biC5Y2O5a6H6K6v56eR5oqA5pyJ6ZmQ5YWs5Y+4IiwiZGVwdElkIjoxLCJkZXB0TmFtZSI6Iua3seWcs+W4guWNjuWuh+iur+enkeaKgOaciemZkOWFrOWPuCIsImVtcGxveWVlSWQiOjEsImVtcGxveWVlTmFtZSI6Iui2hee6p+euoeeQhuWRmCIsImVuZ2xpc2hOYW1lIjoiZXJwYWRtaW4iLCJmaWx0ZXJDb21wSWRzIjpbXSwiZmlsdGVyRGVwdElkcyI6W10sImZpbHRlcklkcyI6W10sInBob25lIjoiIiwicGxhdGZvcm1UeXBlIjoxMDgsInJvbGVDb2RlTGlzdCI6WyIxMiJdLCJ0eXBlIjozfQ=="
    console.log(JSON.parse(Buffer.from(token, 'base64').toString()));
}

function encrypt() {
    const data = {"userId":2028,"userName":"yi","deptId":923,"deptName":" ","compId":17324,"compName":" ","englishName":"yi","roleCodeList":[" "],"operatorId":10001,"platformType":1,"thirdCompId":23,"systemVersion":0,"isTrial":0}
    console.log(Buffer.from(JSON.stringify(data)).toString('base64'));
}

function openErp(){
    const sess = {
        piId:"10001",
        employeeId: 2,
        employeeName: "xxx",
        deptId: 1,
        deptName: "xxx",
        compId: 5,
        compName: "xxx",
        englishName: "xxx",
        phone: "xxxxx",
        roleCodeList:["20","21","22"],
        platformType:103,
        accountType: 4,
        thirdCompId:10010,
        thirdCompCode: "",
        source:0,
        systemVersion:0,
        isTrial:0,
        type: 1,
        filterIds: [1,2,3],
        filterDeptIds: [],
        filterCompIds: []
    };
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'));
}

const all  = "e30="

console.log(Buffer.from('{}').toString('base64'));
escape()