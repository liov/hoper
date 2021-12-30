
function isArray(obj) {
   return  Object.prototype.toString.call(obj).slice(8, -1) === 'Array'
}

function isArray2(obj) {
   return  Array.isArray(obj)
}