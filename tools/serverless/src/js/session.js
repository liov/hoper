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
    const token = "eyJhY2NvdW50VHlwZSI6MSwiY29tcElkIjoxMDAwMSwiY29tcE5hbWUiOiLku6PnkIbllYYxLWNvbXAiLCJkZXB0SWQiOjIsImRlcHROYW1lIjoi5rex5Zyz5biC5Lit5aix5a6P5Zu+5paH5YyW5oqV6LWE5pyJ6ZmQ5YWs5Y+4IiwiZW1wbG95ZWVJZCI6NywiZW1wbG95ZWVOYW1lIjoi6L+Q6JCl5ZWGIiwiZW5nbGlzaE5hbWUiOiJvcGVyYXRvciIsImlzVHJpYWwiOjAsInBob25lIjoiMTUwMTg1MjE2NzMiLCJwbGF0Zm9ybVR5cGUiOjEwMSwicm9sZUNvZGVMaXN0IjpbInJvbGVPcGVyYXRvcjAwMSIsIjI3MjYiXSwic291cmNlIjoxLCJzeXN0ZW1WZXJzaW9uIjowLCJ0aGlyZENvbXBDb2RlIjoiIiwidGhpcmRDb21wSWQiOjB9"
    console.log(JSON.parse(Buffer.from(token, 'base64').toString()));
}

escape()