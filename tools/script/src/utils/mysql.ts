import mysql from "mysql2/promise";

export async function connect(db:mysql.Connection){
    await db.connect().catch(err => console.error('error connecting: ' + err.stack));
    console.log('connected as id ' + db.threadId);
}
