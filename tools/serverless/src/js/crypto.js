import crypto from 'crypto';
import CryptoJS from 'crypto-js';
import Base64  from 'crypto-js/enc-base64.js';

const tmp = ["WPE6bdcCIZYIGdCuy8rAWeCI0MLeytXrUGRgQ+qS7u"];
const message =`{"code":200,"msg":"success","time":16343130`
const key = CryptoJS.enc.Utf8.parse('123456789');
console.log(CryptoJS.AES.encrypt(message,key,{
    mode: CryptoJS.mode.ECB,
    padding: CryptoJS.pad.Pkcs7
}).toString());
//console.log(CryptoJS.AES.decrypt(tmp[1],key).toString());
console.log(Buffer.from(message).toString('base64'));
//console.log(Encrypt(message,key,""))
/**
 * 解密
 * @param dataStr {string}
 * @param key {string}
 * @param iv {string}
 * @return {string}
 */
function Decrypt(dataStr, key, iv) {
    let cipherChunks = [];
    let decipher = crypto.createDecipheriv('aes-128-cbc', key, iv);
    decipher.setAutoPadding(true);
    cipherChunks.push(decipher.update(dataStr, 'base64', 'utf8'));
    cipherChunks.push(decipher.final('utf8'));
    return cipherChunks.join('');
}

/**
 * 加密
 * @param dataStr {string}
 * @param key {string}
 * @param iv {string}
 * @return {string}
 */
function Encrypt(dataStr, key, iv) {
    let cipherChunks = [];
    let cipher = crypto.createCipheriv('aes-128-cbc', key, iv);
    cipher.setAutoPadding(true);
    cipherChunks.push(cipher.update(dataStr, 'utf8', 'base64'));
    cipherChunks.push(cipher.final('base64'));
    return cipherChunks.join('');
}
