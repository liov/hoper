function gen() {
    let header = {
        "piId": "",
        "employeeId": 1,
        "employeeName": "",
        "deptId": 0,
        "deptName": "",
        "compId": 11103,
        "compName": "",
        "englishName": "",
        "phone": "",
        "roleCodeList": null,
        "type": 0,
        "filterIds": null,
        "filterDeptIds": null,
        "filterCompIds": [10002]
    }
    console.log(Buffer.from(JSON.stringify(header)).toString('base64'))
}

gen()