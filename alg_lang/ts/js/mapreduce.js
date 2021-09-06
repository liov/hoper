const s = new Set();

[2, 3, 5, 4, 5, 2, 2].forEach(x => s.add(x));

for (let i of s) {
    console.log(i);
}

const m = new Map();
const o = {p: 'Hello World'};

m.set(o, 'content')
m.get(o) // "content"

m.has(o) // true
m.delete(o) // true
m.has(o) // false


console.log([1, 2, 3].map(function (n) {
    return n + 1;
}))
console.log([1, 2, 3, 4, 5].filter(function (elem) {
    return (elem > 3);
}))
console.log([1, 2, 3, 4, 5].reduce(function (a, b) {
    return a + b;
}))
