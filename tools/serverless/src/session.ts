interface ErpUserPara {
    piId?: string
    employeeId?: number
    employeeName?: string
    platformType?: number
    accountType?: number
    thirdCompId?: number
    deptId?: number
    deptName?: string
    compId?: number
    compName?: string
    englishName?: string
    phone?: string
    roleCodeList?: string[]
    type: number
    filterIds?: number[]
    filterDeptIds?: number[]
    filterCompIds?: number[]
}

interface Session {
    userID: number
    userName: string
    uerRealName: string
    clientIp: string
}

function generate() {
    const sess: ErpUserPara = {
        employeeId: 1,
        type: 3,
        compId:1,
        filterCompIds:[10001]
    }
    console.log(Buffer.from(JSON.stringify(sess)).toString('base64'))
}

generate()