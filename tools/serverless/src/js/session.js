const skipValidation = "e30=";

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
    const token = "eyJwaUlkIjoiIiwiZW1wbG95ZWVJZCI6NywiZW1wbG95ZWVOYW1lIjoi6L+Q6JCl5ZWGIiwiZGVwdElkIjoyLCJkZXB0TmFtZSI6Iua3seWcs+W4guS4reWoseWuj+WbvuaWh+WMluaKlei1hOaciemZkOWFrOWPuCIsImNvbXBJZCI6MTAwMDEsImNvbXBOYW1lIjoi5Luj55CG5ZWGMS1jb21wIiwiZW5nbGlzaE5hbWUiOiJvcGVyYXRvciIsInBob25lIjoiMTUwMTg1MjE2NzMiLCJyb2xlQ29kZUxpc3QiOlsiMzIiLCI1NCJdLCJwbGF0Zm9ybVR5cGUiOjEwMSwiYWNjb3VudFR5cGUiOjEsInNvdXJjZSI6MSwidHlwZSI6MCwiZmlsdGVySWRzIjpudWxsLCJmaWx0ZXJEZXB0SWRzIjpudWxsLCJmaWx0ZXJDb21wSWRzIjpudWxsLCJvcGVyYXRvcklkTGlzdCI6WzEwMDAxXSwic3lzdGVtVmVyc2lvbiI6MCwiaXNUcmlhbCI6MH0="
    console.log(JSON.parse(Buffer.from(token, 'base64').toString()));
}

function encrypt() {
    const data = {"userId":2028,"userName":"yi","deptId":923,"deptName":" ","compId":17324,"compName":" ","englishName":"yi","roleCodeList":[" "],"operatorId":10001,"platformType":1,"thirdCompId":23,"systemVersion":0,"isTrial":0}
    console.log(Buffer.from(JSON.stringify(data)).toString('base64'));
}

crm()