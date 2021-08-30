import {NacosConfigClient} from 'nacos';   // ts
//const NacosConfigClient = require('nacos').NacosConfigClient; // js

const acm = {
    "endpoint": "acm.aliyun.com",
    "namespace": "xxx",
    "accessKey": "xxx",
    "secretKey": "xxx",
    dataId:'xxx',
    group:'xx'
}
// for find address mode
const configClient = new NacosConfigClient({
    endpoint:  acm.endpoint,
    namespace: acm.namespace,
    accessKey: acm.accessKey,
    secretKey: acm.secretKey,
    requestTimeout: 6000,
});

/*
// for direct mode
const configClient = new NacosConfigClient({
    serverAddr: '127.0.0.1:8848',
});
*/

// get config once
const content = await configClient.getConfig(acm.dataId, acm.group);
console.log('getConfig = ', content);
/*

// listen data changed
configClient.subscribe({
    dataId: 'test',
    group: 'DEFAULT_GROUP',
}, content => {
    console.log(content);
});

// publish config
const content= await configClient.publishSingle('test', 'DEFAULT_GROUP', '测试');
console.log('getConfig = ',content);

// remove config
await configClient.remove('test', 'DEFAULT_GROUP');*/
