import { dynamicLoadJs } from "@/plugin/utils/script";

let active = false;
// web-view下的页面内
function ready() {
  window.WeixinJSBridge.on("onPageStateChange", function (res) {
    console.log("res is active", res.active);
    active = res.active;
  });

  // 或者
  window.wx.miniProgram.getEnv(function (res) {
    console.log(res.miniprogram); // true
  });
}

const wxsdk = "https://res.wx.qq.com/open/js/jweixin-1.3.2.js";

function loadwxSDK() {
  dynamicLoadJs(wxsdk);
}

if (!window.WeixinJSBridge || !window.WeixinJSBridge.invoke) {
  document.addEventListener("WeixinJSBridgeReady", ready, false);
} else {
  ready();
}

function IsWeappPlatform(): boolean {
  return window.__wxjs_environment === "miniprogram";
}

export default {
  active,
  IsWeappPlatform,
  loadwxSDK,
};
