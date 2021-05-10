import  express from 'express'
const app = express()

let count = 0;
app.post('/api/exportOrder/exportList/v1', function (req, res) {
    if (count < 100){
        res.json({})
        count++
    }else res.json({"code":0,"data":{"total":0,"list":null}})

})

app.listen(3000)