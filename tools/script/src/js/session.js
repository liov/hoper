

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
    const token = "eyJhY2NvdW50VHlwZSI6MSwiY29tcElkIjoxNDcxOCwiY29tcE5hbWUiOiLnpLzmnI3pgJrnpLzlk4Hlhazlj7giLCJkZXB0SWQiOjUyMTgsImRlcHROYW1lIjoi56S85pyN6YCa56S85ZOB5YWs5Y+4IiwiZW1wbG95ZWVJZCI6MTM1NDEsImVtcGxveWVlTmFtZSI6IueuoeeQhuWRmCIsImVuZ2xpc2hOYW1lIjoibGlmdXRvbmd2MTc4IiwiZmlsdGVyQ29tcElkcyI6WzE0NzE4XSwiZmlsdGVyRGVwdElkcyI6W10sImZpbHRlcklkcyI6W10sImlzVHJpYWwiOjAsInBob25lIjoiMTUwMTg1MjE2NzMiLCJwbGF0Zm9ybVR5cGUiOjEwMSwicm9sZUNvZGVMaXN0IjpbIjI4IiwiMjAiXSwicm9sZUNvZGVzIjpbIm9wX2FkbWluIiwidGVzdF9yb2xlMiJdLCJyb2xlSWRzIjpbIjI4IiwiMjAiXSwicm9sZU5hbWVzIjpbIueuoeeQhuWRmCIsIumihOeUn+S6p+eJiOacrOmqjOaUtuinkuiJsiJdLCJzeXN0ZW1WZXJzaW9uIjozLCJ0aGlyZENvbXBJZCI6MCwidHlwZSI6M30="
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
console.log(btoa("TY1lPZZPBVuX8h4Rij8QIg=="))
console.log(atob("TY1lPZZPBVuX8h4Rij8QIg=="))