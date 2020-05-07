async function test():Promise<number>{
    return new Promise<number>(((resolve, reject) => {resolve(5)}));
}
