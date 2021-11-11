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
        userId: 310,
        userName: '',
        userRealName: '',
        clientIp: '',
    };
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'));
}

function escape() {
    const token = "eyJhY2NvdW50VHlwZSI6MCwiY29tcElkIjoxLCJjb21wTmFtZSI6Iua3seWcs+W4guWNjuWuh+iur+enkeaKgOaciemZkOWFrOWPuCIsImRlcHRJZCI6MTkzLCJkZXB0TmFtZSI6IkVSUOiZmuaLn+mDqOmXqCIsImVtcGxveWVlSWQiOjYsImVtcGxveWVlTmFtZSI6IuW5s+WPsOaWuSIsImVuZ2xpc2hOYW1lIjoicGxhdGZvcm0iLCJpc1RyaWFsIjowLCJwaG9uZSI6IjE2Njc1NTIzMjAyIiwicGxhdGZvcm1UeXBlIjoxMDAsInJvbGVDb2RlTGlzdCI6WyIxMSIsIjI3MjYiXSwic291cmNlIjoxLCJzeXN0ZW1WZXJzaW9uIjowLCJ0aGlyZENvbXBDb2RlIjoiIiwidGhpcmRDb21wSWQiOjB9"
    console.log(JSON.parse(Buffer.from(token, 'base64').toString()));
}

escape()