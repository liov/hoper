const VERSION = "v1";

interface IConfiguration {
  name: string;
}

let configuration: IConfiguration;

const filesToCache = ["../static/template/liov.html"];

// 缓存
const initialize = (service: ServiceWorkerGlobalScope): void => {
  service.addEventListener("install", function (event) {
    event.waitUntil(
      caches.open(VERSION).then(function (cache) {
        return cache.addAll(filesToCache);
      })
    );
  });

  // 缓存更新
  service.addEventListener("activate", function (event) {
    event.waitUntil(
      caches.keys().then(function (cacheNames) {
        return Promise.all(
          cacheNames.map(function (cacheName) {
            // 如果当前版本和缓存版本不一致
            if (cacheName !== VERSION) {
              return caches.delete(cacheName);
            }
          })
        );
      })
    );
  });

  // 捕获请求并返回缓存数据
  service.addEventListener("fetch", function (event) {
    event.respondWith(
      caches
        .match(event.request)
        .catch(function () {
          console.log(event.request);
          return fetch(event.request)
            .then(function (response) {
              caches.open(VERSION).then(function (cache) {
                cache.put(event.request, response!);
              });
              return response!.clone();
            })
            .then(function (response) {
              return response.text();
            })
            .then(function (body) {
              console.log(body);
            })
            .catch(function () {
              return caches.match("../static/template/liov.html");
            });
        })
        .then(function (response) {
          caches.open(VERSION).then(function (cache) {
            cache.put(event.request, response!);
          });
          return response!.clone();
        })
    );
  });
};

initialize(self as any);

export {};
