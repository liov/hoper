export function startWorker(url: string): Worker {
  if (typeof Worker !== "undefined") {
    return new Worker(url);
  } else {
    throw new Error("抱歉，你的浏览器不支持 Web Workers...");
  }
}
