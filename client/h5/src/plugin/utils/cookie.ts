export default class CookieUtils {
  // 获取cookie
  static getCookie(key: string, cookie: string): string {
    key = key + "=";
    const decodedCookie = decodeURIComponent(cookie);
    const ca = decodedCookie.split(";");
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === "") {
        c = c.substring(1);
      }
      if (c.indexOf(key) === 0) {
        return c.substring(key.length, c.length);
      }
    }
    return "";
  }
}
