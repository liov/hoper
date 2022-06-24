//生成从minNum到maxNum的随机数
export function randomNum(min?:number,max?:number):number{
    switch(arguments.length){
        case 1:
            return  Math.floor(Math.random()*min!+1);
        case 2:
            return Math.floor(Math.random()*(max!-min!+1)+min!);
        default:
            return 0;
    }
}

export function lottery(probability:number):boolean{
    if(probability === 100) return true
    const odds = Math.floor(Math.random() * 100);
    return odds < probability;
};
