// 从文件读取
import Excel from "exceljs";
import fs from "fs";
import dayjs from "dayjs";

const workbook = new Excel.Workbook();
await workbook.xlsx.readFile("D:/work/sql/同步客户.xlsx");

const sheet = workbook.worksheets[0];
const rows = sheet.findRows(2,sheet.rowCount-2);

const file  = fs.createWriteStream("D:/work/sql/同步客户.sql",{flags:"w",mode:0o666,encoding:"utf-8"});

for(let row of rows) {
    const obj = {
        name:row.getCell("F").value,
        owner_id:row.getCell("B").value,
        customer_num:row.getCell("E").value,
        alias:row.getCell("H").value,
        level:row.getCell("I").value + 1,
        info_source:row.getCell("J").value,
        property:row.getCell("M").value,
        estimated_demand_time:row.getCell("N").value,
        registered_fund:row.getCell("O").value,
        employee_number:row.getCell("P").value,
        business_scope:row.getCell("Q").value,
        address:row.getCell("R").value,
        zipcode:row.getCell("S").value,
        telephone:row.getCell("T").value,
        category:row.getCell("Y").value,
        uniform_social_credit_code:row.getCell("V").value,
        created_at:row.getCell("AI").value,
        updated_at:row.getCell("AJ").value,
        created_by:row.getCell("AL").value,
        applicant_id:row.getCell("AL").value,
        req_region_range:row.getCell("AU").value,
        req_history_order:row.getCell("AV").value,
        req_demand_type:row.getCell("AW").value,
        req_product_type:row.getCell("AX").value,
        req_cooperation_type:row.getCell("BE").value,
        req_settlement_type:row.getCell("BG").value,
        req_remark:row.getCell("BL").value,
    }
    const customer = {
        name:obj.name?obj.name:"",
        owner_id:obj.owner_id?obj.owner_id:0,
        customer_num:obj.customer_num?obj.customer_num:"",
        alias:obj.alias?obj.alias:"",
        level:obj.level?obj.level:0,
        info_source:obj.info_source?obj.info_source:0,
        property:obj.property?obj.property:0,
        estimated_demand_time:obj.estimated_demand_time?dayjs(obj.estimated_demand_time).add(-8, 'hour').format("YYYY-MM-DD"):"0000-00-00 00:00:00",
        registered_fund:obj.registered_fund?obj.registered_fund:"0",
        employee_number:obj.employee_number?obj.employee_number:0,
        business_scope:obj.business_scope?obj.business_scope:"",
        address:obj.address?obj.address:"",
        zipcode:obj.zipcode?obj.zipcode:"",
        telephone:obj.telephone?obj.telephone:"",
        category:obj.category?obj.category:0,
        uniform_social_credit_code:obj.uniform_social_credit_code?obj.uniform_social_credit_code:"",
        created_at:obj.created_at?dayjs(obj.created_at).add(-8, 'hour').format("YYYY-MM-DD HH:mm:ss"):"0000-00-00 00:00:00",
        updated_at:obj.updated_at?dayjs(obj.updated_at).add(-8, 'hour').format("YYYY-MM-DD HH:mm:ss"):"0000-00-00 00:00:00",
        created_by:obj.created_by?obj.created_by:0,
        applicant_id:obj.applicant_id?obj.applicant_id:0,
        status:1,
        succeeded:1,
        stage:1,
        req_region_range:obj.req_region_range?obj.req_region_range:"",
        req_history_order:obj.req_history_order?obj.req_history_order:"",
        req_demand_type:obj.req_demand_type?obj.req_demand_type:0,
        req_product_type:obj.req_product_type?obj.req_product_type:"",
        req_cooperation_type:obj.req_cooperation_type?obj.req_cooperation_type:0,
        req_settlement_type:obj.req_settlement_type?obj.req_settlement_type:0,
        req_remark:obj.req_remark?obj.req_remark:"",
    }

    if (customer.employee_number<50) customer.employee_number = 1;
    else if (customer.employee_number>=50 && customer.employee_number<100) customer.employee_number = 2;
    else if (customer.employee_number>=100 && customer.employee_number<200) customer.employee_number = 3;
    else if (customer.employee_number>=200 && customer.employee_number<300) customer.employee_number = 4;
    else if (customer.employee_number>=300 && customer.employee_number<400) customer.employee_number = 5;
    else customer.employee_number = 6;

    file.write(`INSERT INTO customer_info(name,alias,level,category,property,registered_fund,employee_number,business_scope,address,zipcode,telephone,info_source,req_region_range,req_history_order,req_remark,req_demand_type,req_product_type,req_cooperation_type,req_settlement_type,owner_id,applicant_id,status,created_by,created_at,updated_at,uniform_social_credit_code,customer_num,succeeded,stage,estimated_demand_time) VALUES('${customer.name}','${customer.alias}',${customer.level},${customer.category},${customer.property},'${customer.registered_fund}',${customer.employee_number},'${customer.business_scope}','${customer.address}','${customer.zipcode}','${customer.telephone}',${customer.info_source},'${customer.req_region_range}','${customer.req_history_order}','${customer.req_remark}','${customer.req_demand_type}','${customer.req_product_type}',${customer.req_cooperation_type},${customer.req_settlement_type},${customer.owner_id},${customer.applicant_id},${customer.status},${customer.created_by},'${customer.created_at}','${customer.updated_at}','${customer.uniform_social_credit_code}','${customer.customer_num}',${customer.succeeded},${customer.stage},'${customer.estimated_demand_time}');\n`)
    file.write(`SET @pid:=LAST_INSERT_ID();\n`)
    file.write(`INSERT INTO customer_extra_info(customer_id,last_visit_time,last_deal_time) VALUES(@pid,'${customer.created_at}','${customer.estimated_demand_time}');\n`)
    file.write(`INSERT INTO customer_erptask(customer_id,task_created_at,created_by,created_at) VALUES(@pid,'${customer.created_at}',${customer.created_by},'${customer.created_at}');\n`)
}
file.end();