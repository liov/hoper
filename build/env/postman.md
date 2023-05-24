# 一级目录
```javascript
pm.collectionVariables.set("proxyHost", "http://127.0.0.1:8001/api/v1/namespaces");


let erpHeaderBase64 = "e30="

// 可以设置不同环境判断
let raw  =  true
const erpHeader = {
    compId:10001,
    systemVersion:2,
    platformType:101,
    type:3,
    filterCompIds:[10001]
};


if (!raw) {
    erpHeaderBase64 = btoa(JSON.stringify(erpHeader))
}

pm.request.addHeader({
    key: "header",
    value: erpHeaderBase64,
})



switch (pm.environment.name) {
case "local" :{
    pm.collectionVariables.set("baseUrl", "localhost:8080");
    break;
}
}


```

## 二级目录
`pm.collectionVariables.set("namespace", "xxx");`

### 三级目录
```javascript
const serviceName = "xxx-center";
const namespace = pm.collectionVariables.get("namespace");
const proxyHost = pm.collectionVariables.get("proxyHost");

switch (pm.environment.name) {
case "dev" :{
    pm.collectionVariables.set("baseUrl", `http://${serviceName}.${namespace}`);
    break;
}
case "proxy" :{
    pm.collectionVariables.set("baseUrl", `${proxyHost}/${namespace}/services/${serviceName}:80/proxy`);
    break;
}
}

// 如果有特殊请求头
const scmHeader = {
  account:{
      supplierId:35
  }
};

pm.request.addHeader({
    key: "header",
    value: btoa(JSON.stringify(scmHeader)),
},)
```