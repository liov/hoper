
const x = {
    a:1,
    b:"1",
    c: {
        d:"e"
    }
}

console.log({...x,a:3})
console.log({a:3,...x})