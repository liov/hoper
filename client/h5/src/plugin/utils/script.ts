export const dynamicLoadJs = function(url, callback?) {
    const head = document.getElementsByTagName("head")[0];
    const script:HTMLScriptElement = document.createElement("script");
    script.type = "text/javascript";
    script.src = url;
    if (callback && typeof callback == "function") {
        script.onload = function() {
            if (!document.readyState || document.readyState === "complete" ) {
                callback();
                script.onload = null;
            }
        };
    }
    head.appendChild(script);
  }