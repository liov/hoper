
function second(tmp) {
    let obj = {};
    let temp = {}
    for (let k in tmp){
        if (k.indexOf('.')>0){
            let list = k.split('.');
            for (let i in list){
                if(!obj[list[i]]){
                    obj[list[i]] ={}
                    temp = obj[list[i]]
                }
            }
            temp[list[list.length-1]] = tmp[k]
        }else {
            obj[k] = tmp[k]
        }
    }
    return obj;
}

let obj = { 'A': 1, 'B.A': 2, 'B.B': 3, 'CC.D.E': 4, 'CC.D.F': 5}

console.log(second(obj))
