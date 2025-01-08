import { dynamicLoadJs } from "@hopeio/utils/browser/script";

let active = false;

const wxsdk = "https://res.wx.qq.com/open/js/jweixin-1.3.2.js";

function loadwxSDK() {
  dynamicLoadJs(wxsdk, () =>
    // 或者
    window.wx.miniProgram.getEnv(function (res) {
      console.log(res.miniprogram); // true
    })
  );
  weBrowser();
}

function weBrowser() {
  if (!window.WeixinJSBridge || !window.WeixinJSBridge.invoke) {
    document.addEventListener("WeixinJSBridgeReady", ready, false);
  } else {
    ready();
  }
}

// web-view下的页面内
function ready() {
  window.WeixinJSBridge.on("onPageStateChange", function (res) {
    console.log("res is active", res.active);
    active = res.active;
  });
}

function IsWeappPlatform(): boolean {
  return window.__wxjs_environment === "miniprogram";
}

export default {
  active,
  IsWeappPlatform,
  loadwxSDK,
};
