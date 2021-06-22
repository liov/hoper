import  express from 'express'
const app = express()

let count = 0;
const result = {status:0,data:{total:0,list:null}}
const orderInfo = {"orderInfo":{"orderId":"2021050600325920349","payOrderId":"2021050600325862040","operatorId":10109,"operatorName":"八闽福","orderType":0,"payTime":1620232382,"orderStatus":13,"orderStatusName":"交易完成","serviceStatus":0,"userId":2995862,"name":"刘莹","phone":"18030126715","address":"福建省厦门市同安区西柯镇美峰1里保利·叁仟栋1期4号楼","remark":"","totalPrice":11010,"expressAmount":0,"epayAmount":11010,"externalPayAmount":0,"thirdOrderId":"168381535737","supplierRemark":"","deliveryType":1,"packageType":1,"supplierId":11,"supplierName":"苏打直供","supplierModel":1,"storeId":0,"storeName":"","userExpectDeliveryTime":"","deliveryTime":1620272539,"deliveryEndTime":1620272539,"createdAt":1620232382,"updatedAt":1620608400},"skuList":[{"skuId":161620039,"skuName":"爪心包夜280*9片","productId":161595306,"productName":"高洁丝 Kotex 澳洲进口纯棉  夜用迷你卫生巾 280mm* 9片","count":1,"price":1500,"costPrice":896,"operatorCostPrice":941,"skuCode":"100004050119","barCode":"100004050119","productSetInfoList":null},{"skuId":161659655,"skuName":"日用180mm*16片","productId":161629717,"productName":"七度空间（SPACE7） 少女极薄系列迷你卫生巾 ","count":1,"price":1500,"costPrice":846,"operatorCostPrice":889,"skuCode":"100007447401","barCode":"100007447401","productSetInfoList":null},{"skuId":161659180,"skuName":"超薄338mm*12片","productId":161629461,"productName":"七度空间（SPACE7） 少女超薄纯棉 加量超长甜睡夜用卫生巾","count":1,"price":1560,"costPrice":983,"operatorCostPrice":1033,"skuCode":"100003140579","barCode":"100003140579","productSetInfoList":null},{"skuId":121146926,"skuName":"3层*130抽*6包 ","productId":121132455,"productName":"五月花妇婴柔抽纸 3层*130抽*6包","count":1,"price":2390,"costPrice":1280,"operatorCostPrice":1344,"skuCode":"2327493","barCode":"2327493","productSetInfoList":null},{"skuId":161664188,"skuName":"紫色","productId":161632549,"productName":"惠寻 遇水开花晴雨伞","count":1,"price":1990,"costPrice":1280,"operatorCostPrice":1344,"skuCode":"100009625423","barCode":"100009625423","productSetInfoList":null},{"skuId":161623811,"skuName":"甄选虾皮58g","productId":161597607,"productName":"禾煜 甄选虾皮58g ","count":1,"price":2070,"costPrice":1400,"operatorCostPrice":1470,"skuCode":"100001554183","barCode":"100001554183","productSetInfoList":null}],"packageList":[{"packageId":3295716,"expressCompanyId":19,"expressCompanyName":"京东","expressId":"JD0042087742639","signTime":1620297921,"deliveryContactPerson":"","deliveryPhone":"","packageSkuList":[{"skuId":161659180,"masterSkuId":0,"count":1}],"createdAt":1620272539,"updatedAt":1620297921},{"packageId":3295717,"expressCompanyId":19,"expressCompanyName":"京东","expressId":"JD0042087736598","signTime":1620297919,"deliveryContactPerson":"","deliveryPhone":"","packageSkuList":[{"skuId":161659655,"masterSkuId":0,"count":1},{"skuId":161664188,"masterSkuId":0,"count":1},{"skuId":121146926,"masterSkuId":0,"count":1},{"skuId":161620039,"masterSkuId":0,"count":1}],"createdAt":1620272539,"updatedAt":1620297919},{"packageId":3295718,"expressCompanyId":2,"expressCompanyName":"圆通","expressId":"YT3155428214388","signTime":1620607822,"deliveryContactPerson":"","deliveryPhone":"","packageSkuList":[{"skuId":161623811,"masterSkuId":0,"count":1}],"createdAt":1620272539,"updatedAt":1620607822}],"orderDeliveryExtInfoList":null}
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