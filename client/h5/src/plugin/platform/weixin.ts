const WeixinJSBridge = window.WeixinJSBridge;
const wx = window.wx;
let active = false;
// web-view下的页面内
function ready() {
  WeixinJSBridge.on("onPageStateChange", function (res) {
    console.log("res is active", res.active);
    active = res.active;
  });

  // 或者
  wx.miniProgram.getEnv(function (res) {
    console.log(res.miniprogram); // true
  });
}

if (!WeixinJSBridge || !WeixinJSBridge.invoke) {
  document.addEventListener("WeixinJSBridgeReady", ready, false);
} else {
  ready();
}

function IsWeappPlatform(): boolean {
  return window.__wxjs_environment === "miniprogram";
}

export default {
  wx,
  WeixinJSBridge,
  active,
  IsWeappPlatform,
};
