import mysql from 'mysql';

const connection = mysql.createConnection({
    host: 'example.org',
    user: 'bob',
    password: 'secret'
});

connection.connect(function(err) {
    if (err) {
        console.error('error connecting: ' + err.stack);
        return;
    }

    console.log('connected as id ' + connection.threadId);
});

connection.query('SELECT 1', function (error, results, fields) {
    if (error) throw error;
    // connected!
});