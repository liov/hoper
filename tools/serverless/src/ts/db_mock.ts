import mysql, {RowDataPacket} from 'mysql2/promise';
import bluebird from 'bluebird';
import {randomNum, lottery} from '../utils/random.js';


const connection = await mysql.createConnection({
    host: 'rm-bp1y6kuovzmbz1757.mysql.rds.aliyuncs.com',
    user: 'web',
    password: '123456',
    Promise: bluebird
});

const DB1 = `d_aura_jike.`
const DB2 = `d_scm_product_3rd.`

await connection.connect().catch(err => console.error('error connecting: ' + err.stack));
console.log('connected as id ' + connection.threadId);

const [rows1,] = await connection.query(`SELECT MAX(id) AS max_id
                                         FROM ${DB1}sku_info
                                         WHERE status = 0`)
console.log(rows1)
const max_id = (rows1 as RowDataPacket[])[0].max_id as number

console.log(max_id)
if (max_id === 0) throw Error("获取最大id失败")
const start_id = randomNum(11930243, max_id)
console.log(start_id)
const [rows2,] = await connection.query(`SELECT a.id, a.price
                                         FROM ${DB1}sku_info a
                                                  LEFT JOIN ${DB1}product_info b ON a.product_id = b.id
                                         WHERE a.id > ${start_id}
                                           AND a.status = 0
                                           AND b.status = 0
                                            AND b.product_source = 1
                                         LIMIT 500`);


// 异步join
Promise.all((rows2 as RowDataPacket[]).map((row)=>{
    // 这个sku配置了
    if (lottery(90)) return  insert(row.id, row.price)
})).then(()=> connection.end())


function insert(skuId: number, price: number) {
    const asyncArray = [];
    // 这个sku配置了两个平台
    if (lottery(90)) {
        asyncArray.push(comp(skuId, price, 1),comp(skuId, price, 2))
    } else {
        asyncArray.push(comp(skuId, price, lottery(50) ? 1 : 2))
    }
    return Promise.all(asyncArray)
}

async function comp(skuId: number, price: number, platform: number) {
    const count = (() => {
        const randNum = Math.floor(Math.random() * 100)
        switch (true) {
            case randNum < 90:
                return 1
            case randNum < 95:
                return 2
            case randNum < 100:
                return 3
            default:
                return 1
        }
    })()

    const diffRate = (() => {
        // 50的概率是1
        if (lottery(50)) {
            return 1
        } else {
            // 50概率正向或负向
            if (lottery(50)) return 1 + Math.random() * 2
            else return 1 - Math.random() * 2
        }
    })()

    const [comp_price, comp_seckill_price] = (() => {
        if (lottery(90)) {
            return [price * diffRate, 0]
        } else return [0, price * 0.8]
    })()
    const [rows1,] = await connection.query(`SELECT COUNT(*) AS count
                                             FROM ${DB2}sku_competition
                                             WHERE sku_id = ${skuId}
                                               AND platform = ${platform}
                                               AND status = 0`)
    let exists = (rows1 as RowDataPacket[])[0].count as number

    if (exists == 0) {
       await connection.execute(`INSERT INTO ${DB2}sku_competition (sku_id, platform, count, comp_price, comp_seckill_price)
                            VALUES (${skuId}, ${platform}, ${count}, ${comp_price}, ${comp_seckill_price})`)

    } else {
       await connection.execute(`UPDATE ${DB2}sku_competition
                            SET count             =${count},
                                comp_price        = ${comp_price},
                                comp_seckill_price= ${comp_seckill_price}
                            WHERE sku_id = ${skuId}
                              AND platform = ${platform}
                              AND status = 0`)
    }
}

