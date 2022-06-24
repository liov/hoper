import express from 'express'
import proxy from 'express-http-proxy'

const app = express()


app.put('/webhdfs/v1/:file', proxy('hadoop-http.tools:9870', {
    userResDecorator: function (rsp, data, req, res) {
        if (req.query.op !== "CREATE"){
            return data
        }
        console.log(data.toString());
        const hostname = req.hostname.replace('hadoop', 'hdfs');
        let body = data.toString().replace("host:9864", hostname);
        if (hostname !== 'hdfs.d') {
            body.replace("http", "https")
        }
        return body
    },
}))

app.get('/webhdfs/v1/:file', proxy('hadoop-http.tools:9870', {
    userResHeaderDecorator: function (headers, req, res, preq,pres) {
        if (req.query.op !== "OPEN"){
            return headers
        }
        console.log(headers.location.toString());
        const hostname = req.hostname.replace('hadoop', 'hdfs');
        let location = headers.location.toString().toString().replace("host:9864", hostname);
        if (hostname !== 'hdfs.d') {
            location.replace("http", "https")
        }
        headers.location = location
        return headers
    },
}))

app.use('/webhdfs/v1/', proxy('hadoop-http.tools:9870', {
    proxyReqPathResolver:function (req) {
        const queryString = req.url.split('?')[1];
        return  '/webhdfs/v1'+ (queryString ? '?' + queryString : '');
    }
}))

app.listen(3000)