const skipValidation = "e30=";

function erp() {
    const sess = {
        employeeId: 1,
        type: 1,
        platformType: 101,
        thirdCompId: 0,
        compId: 10001,
        filterCompIds: [10001]
    };
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'));
}

erp();

function crm() {
    const sess = {
        userId: 310,
        userName: '',
        userRealName: '',
        clientIp: '',
    };
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'));
}

crm()