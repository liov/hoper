import { showToast } from "@tarojs/taro";

class Toast {
  static text(msg:string) {
    showToast({
      title: msg,
      icon: 'success',
      duration: 2000
    })
  }
  static fail(msg:string) {
    showToast({
      title: msg,
      icon: 'error',
      duration: 2000
    })
  }

}

export {Toast}
