import axios from 'axios';
axios.defaults.headers["Authorization"] = `bearer 125cba5407b20222d37e8fd0aefa1382`
axios.defaults.headers["Cookie"] = `laravel_session=J9LSlILsO9PMWE0lKtonOKhaN3HFxVrs0b7IKM5p`
axios.defaults.headers["Host"] = `api.xiaoyoucaip2p.com`
axios.defaults.headers["user-agent"] = `96bao yu/1.0 (iPhone; iOS 14.6; Scale/3.00)`
axios.defaults.headers["accept-language"] = `zh-Hans-CN;q=1`


async function getUrl(platform=441){
    const map = new Map();
    let res = await axios.get(`https://api.xiaoyoucaip2p.com/api/live/plat/detail?id=${platform}&room_type=0&page=2&pageSize=20`)
    res.data.data.room_list.forEach((element)=>{
        map.set(element.room_id,element)
    })
    res = await axios.get(`https://api.xiaoyoucaip2p.com/api/live/plat/detail?id=${platform}&room_type=1&page=2&pageSize=20`)
    res.data.data.room_list.forEach((element)=>{
        map.set(element.room_id,element)
    })
    map.forEach(async (v,k)=>{
        let res = await axios.get(`https://api.xiaoyoucaip2p.com/api/live/room/detail?id=${k}&plat_id=449`)
        console.log(res.data.data.name,res.data.data.play_url)
    })
}

getUrl()