/** 动态加载微信 JS-SDK */
export function loadwxSDK() {
  if (document.querySelector('script[src*="jweixin"]')) return;
  const script = document.createElement("script");
  script.src = "https://res.wx.qq.com/open/js/jweixin-1.6.0.js";
  document.head.appendChild(script);
}

/** 判断是否在微信小程序内嵌 webview 中 */
export function IsWeappPlatform(): boolean {
  return /miniProgram/i.test(navigator.userAgent) ||
    (window as any).__wxjs_environment === "miniprogram";
}
