import CryptoJS from 'crypto-js';

const key = ""

const object = ['','电影','蛋糕','礼包/灵通券','图书','体检','团建','福利券（商城通）','超级卡','锦福利','定制券'];
const result = CryptoJS.AES.encrypt(JSON.stringify(object),key).toString();
console.log(result);
console.log(JSON.parse(CryptoJS.AES.decrypt(result, key).toString(CryptoJS.enc.Utf8)));

console.log(Buffer.from(JSON.stringify(object)).toString('base64'));