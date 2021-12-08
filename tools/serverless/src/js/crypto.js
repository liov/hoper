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
console.log(CryptoJS.AES.decrypt("U2FsdGVkX1/GmjNjqqcYjiAyluGPwydpoyy6LCYg3UUkTdcWNtmBI+6NPbIEZqtXkCRJj4+QoWDv9PEy9pXSnKeXLPM2yGvfUFTP88WuI9p17GwW05prJra00+xJSYO295puMMYjACguS/mYRF+MqVYrV+0l38Q10TIMWeNWBa6l0CZP5D0RRpJ+J2RgarvZl+MWBAtjk+kici99IEukMPQXHxrfx50+TvI8YI5YkTlQ6b9QHdcELkTJkQZq36mnNVhvsIWss2fW/6LPt3UKRBp2f/LSKAz371J83Q9n4XeQHbTLhu63Jkm/bICr9ZoNttPDGDWK7hnRd7UFP3/ZziXrat5obpSBZr6hELXTouCtyG325Qj3FUgE00/KQVzviDFmbgWk1yMdsbrawnhwWQCqGs7cOOHCIlje0DqSramVu2jTj2vIEbHa3vGD1wF/y1ev+a+VGRaHu4KwHrepPmjTSREqI2EalhwezxKA0izwn/rqe1W34rEQO1DgEXLVzpUCU0jmdMgFytpR3UOTaDiij5ZDRgKFe4n1DmbeIIn9CdJwuORGAZFLGV62IVkNwFL4lNCYBNSa7VlfdcVACZeym/ADVNS9FqfJVdHPYaiZ194z4TV35vYr3coYk/cYfXcor7Y5mYgIipOoqx15BS0AZZkrGY7LEAF2s60TYggz0eOHV+qg21fqlt/iGnVXImNbeKiqMSwSxhUc3IteDERrzQzH9xzR3oyoLjT/RL1m9hz5rZqacl1SYmlAd55GetcJvFyxFXWFqa0VcV3vcVN1vC5PLdBE+vBQwBB+CI8FyuPB2vLT1/yRmpmNPY35+K3+spwc0eoXCCPN18i3i55uIJLP5yYAwKcbcBQl6IeIg0DFHENnmR1APtsW644fOpoCoJRxZkCeJ2xBvTVN7kk5rNZu3YYaEnT/peyhplhUzlYFSHlUt71T3Q40CU+wIrcc+xayXOSIzaI00Jwwqw==", 'ppvod').toString(CryptoJS.enc.Utf8))
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
