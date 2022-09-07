

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
    const token = "eyJ1c2VySWQiOjIxNzYsInVzZXJOYW1lIjoieW95b2d1IiwidXNlclJlYWxOYW1lIjoi6aG+5Zut5ZutIiwidXNlclJvbGUiOiIiLCJjbGllbnRJcCI6IiJ9"
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

escape()