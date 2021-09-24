const skipValidation = "e30=";

function erp() {
    const sess = {
        employeeId: 7,
        type: 1,
        platformType: 101,
        thirdCompId: 0,
        compId: 10001,
        filterCompIds: [10001]
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
    const token = "eyJhY2NvdW50VHlwZSI6MSwiY29tcElkIjoxMDAwMSwiY29tcE5hbWUiOiLku6PnkIbllYYxLWNvbXAiLCJkZXB0SWQiOjE1MjIsImRlcHROYW1lIjoi5Zu96ZmF6ZSA5ZSu5LqM6YOoIiwiZW1wbG95ZWVJZCI6MjAwNiwiZW1wbG95ZWVOYW1lIjoic2FsZTIiLCJlbmdsaXNoTmFtZSI6InNhbGUyIiwiZmlsdGVyQ29tcElkcyI6WzEwMDAxXSwiZmlsdGVyRGVwdElkcyI6W10sImZpbHRlcklkcyI6WzIwMDZdLCJpc1RyaWFsIjowLCJwaG9uZSI6IjEzNDIxMzIxNTkyIiwicGxhdGZvcm1UeXBlIjoxMDEsInJvbGVDb2RlTGlzdCI6WyIzNSJdLCJzeXN0ZW1WZXJzaW9uIjowLCJ0aGlyZENvbXBJZCI6MCwidHlwZSI6M30="
    console.log(JSON.parse(Buffer.from(token, 'base64').toString()));
}

escape()