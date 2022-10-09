import Taro from "@tarojs/taro";

class Router {
  push(url:string) {
    Taro.navigateTo({
      url: url
    })

  }
}

const router = new Router()

export default router;
