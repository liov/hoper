type AsyncFunc<T> = () => Promise<T>;

export function doAsync<T>(func: AsyncFunc<T>): T {
    func().then(
        (res) => {
            return res;
        }
    ).catch((err) => {
        throw err
    })
    throw Error("执行失败")
}

export function makeAsync(){
    return new Promise<number>((resolve, reject) => {

    })
}