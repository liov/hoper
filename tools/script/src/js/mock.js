let settleInfos = [
    {
        "ticketType": 0,
        "feeType": 0,
        "ticketNum": 1,
        "amount": 0
    }
];

for (let i = 1; i < 4; i++) {
    settleInfos.push( {
        "ticketType": 0,
        "feeType": i,
        "ticketNum": 1,
        "amount": 0
    })
}

for (let i = 5; i < 7; i++) {
    settleInfos.push( {
        "ticketType": 1,
        "feeType": i,
        "ticketNum": 1,
        "amount": 0
    })
}

console.log(settleInfos)


function Classification(array){
    return [array.filter(item=>item.ticketType === 0), array.filter(item=>item.ticketType === 1)]
}

console.log(Classification(settleInfos))

function Classification2(array){
    let dianzipiao = [],dingzuopiao =[]
    array.map((item)=>{
        if (item.ticketType === 0) dianzipiao.push(item)
        if (item.ticketType === 1) dingzuopiao.push(item)
    })
    return [dianzipiao,dingzuopiao]
}

console.log(Classification2(settleInfos))