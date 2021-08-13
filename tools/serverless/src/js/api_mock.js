import  express from 'express'
const app = express()

let count = 0;
const result = {status:0,data:{total:0,list:null}}
const orderInfo = {}
const list =  Array.from(new Array(200), () => {
    return orderInfo;
});
app.post('/api/exportOrder/exportList/v1', function (req, res) {
    if (count < 50){
        result.data.list = list
        res.json(result)
        count++
    }else {
        result.data.list = null
        res.json(result)
    }

})

app.listen(3000)