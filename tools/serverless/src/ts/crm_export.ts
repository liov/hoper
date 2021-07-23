import mysql, {RowDataPacket} from 'mysql2/promise';
import bluebird from 'bluebird';
import {connect} from '../utils/mysql.js';
import Excel from 'exceljs';

const REQ_PRODUCT_TYPE = [''];
const STAGE = ['目标客户','开发中客户','意向客户','成交客户','忠实客户','流失客户'];
const crm = await mysql.createConnection({
    host: 'crm_host',
    user: 'crm',
    password: '123456',
    database: 'crm',
    Promise: bluebird
});

const erp = await mysql.createConnection({
    host: 'erp_host',
    user: 'erp',
    password: '123456',
    database: 'erp',
    Promise: bluebird
});



await connect(crm)
await connect(erp)



async  function gen(rows: RowDataPacket[],path: string) {
    const userMap = await getUserInfoMap((rows as RowDataPacket[]).map((row)=>row.owner_id > 0 ? row.owner_id:row.applicant_id));
    const workbook = new Excel.Workbook();
    const sheet = workbook.addWorksheet('Sheet1');
    sheet.addRow(['战区','分部','销售人员','客户名称','客户编号','地址','单位电话','客户阶段','需求产品类型','联系人姓名','性别','职位','联系方式','电子邮箱',"是否在资源池"]);
    (rows as RowDataPacket[]).forEach((row)=>{
        const erpId = (row.owner_id > 0 ? row.owner_id:row.applicant_id) as number;
        const isRespool = row.owner_id > 0 ? "否":"是";
        const erpInfo = userMap.get(erpId) ? (userMap.get(erpId) as RowDataPacket) : { id: 0, name: '', area_name: '', dname: '' }
        const reqProductType = (row.req_product_type as string).split(',').map((v)=>REQ_PRODUCT_TYPE[parseInt(v)]).join(",");
        const stage = STAGE[(row.stage as number)];
        const rowValues = [erpInfo.area_name,erpInfo.dname,erpInfo.name,row.name,row.customer_num,row.address,row.telephone,stage,reqProductType,row.cname,row.gender==1?'男':'女',row.position,row.mobile,row.mail,isRespool];
        sheet.addRow(rowValues);
    })
    await workbook.xlsx.writeFile(path);

}
async function f1() {
    const [rows,] = await crm.query(`SELECT a.owner_id,a.applicant_id,a.name,a.customer_num,a.address,a.telephone,a.stage,a.req_product_type,b.name AS cname,b.gender,b.position,b.mobile,b.mail
FROM customer_info a LEFT JOIN customer_contacts b ON a.id = b.customer_id AND b.is_deleted = 0 WHERE a.status != 2  
AND b.mobile IN (SELECT b.mobile FROM customer_info a LEFT JOIN customer_contacts b ON a.id = b.customer_id AND b.is_deleted = 0 WHERE a.status != 2 AND b.mobile !=''  AND LENGTH(b.mobile) = 11 GROUP BY b.mobile HAVING COUNT(*) > 2)`);
    await gen(rows as RowDataPacket[],'./手机号码重复应用于多个客户的客户详情.xlsx')
}

async function f2() {
    const [rows,] = await erp.query(`SELECT phone FROM o_employee_info WHERE phone !='' `);
    const phones = (rows as RowDataPacket[]).map((row)=>row.phone);

    const [rows1,] = await crm.query(`SELECT a.owner_id,a.applicant_id,a.name,a.customer_num,a.address,a.telephone,a.stage,a.req_product_type,b.name AS cname,b.gender,b.position,b.mobile,b.mail 
FROM customer_info a LEFT JOIN customer_contacts b ON a.id = b.customer_id AND b.is_deleted = 0 WHERE a.status != 2  
AND b.mobile IN (${phones.join(',')})`);

    await gen(rows1 as RowDataPacket[],'./联系人手机号为公司员工客户的客户详情.xlsx')
}

async function f3() {
    const [rows,] = await erp.query(`SELECT phone FROM o_employee_info `);
    const phones = (rows as RowDataPacket[]).map((row)=>row.phone);

    const [rows1,] = await crm.query(`SELECT a.owner_id,a.applicant_id,a.name,a.customer_num,a.address,a.telephone,a.stage,a.req_product_type,b.name AS cname,b.gender,b.position,b.mobile,b.mail 
FROM customer_info a LEFT JOIN customer_contacts b ON a.id = b.customer_id AND b.is_deleted = 0 WHERE a.status != 2  AND  LENGTH(b.mobile) != 11`);

    await gen(rows1 as RowDataPacket[],'./联系人手机号缺位的.xlsx')
}
async function f4() {
    const [rows,] = await erp.query(`SELECT phone FROM o_employee_info `);
    const phones = (rows as RowDataPacket[]).map((row)=>row.phone);

    const [rows1,] = await crm.query(`SELECT a.owner_id,a.applicant_id,a.name,a.customer_num,a.address,a.telephone,a.stage,a.req_product_type,b.name AS cname,b.gender,b.position,b.mobile,b.mail 
FROM customer_info a LEFT JOIN customer_contacts b ON a.id = b.customer_id AND b.is_deleted = 0 WHERE a.status != 2  AND (a.name LIKE '%？%' OR a.name LIKE '%！%' OR a.name LIKE '%，%' OR a.name LIKE '%,%'OR a.name LIKE '%￥%' OR a.name LIKE '%$%' OR a.name LIKE '%#%' OR a.name LIKE '%&%' OR a.name LIKE '%?%' OR a.name LIKE '%.%')`);

    await gen(rows1 as RowDataPacket[],'./客户名称带符号.xlsx')
}
await f4()


crm.end()
erp.end()
async function getUserInfoMap(ids:number[]){
const [rows,] = await erp.query(`SELECT a.id,a.name,a.area_name,b.name AS dname 
FROM o_employee_info a LEFT JOIN o_department_info b ON a.dept_id = b.id AND b.status = 0 
WHERE a.id IN (${ids.join(',')})`);

let m = new Map<number,RowDataPacket>();
(rows as RowDataPacket[]).forEach((row)=>{
    m.set(row.id,row);
});

return m
}


