import dayjs from "dayjs";

declare global {
  interface Date {
    format(fmt: string): string;
  }
}

Date.prototype.format = function (fmt) {
  //author: meizz
  const o = {
    "M+": this.getMonth() + 1, //月份
    "D+": this.getDate(), //日
    "H+": this.getHours(), //小时
    "m+": this.getMinutes(), //分
    "s+": this.getSeconds(), //秒
    "q+": Math.floor((this.getMonth() + 3) / 3), //季度
    S: this.getMilliseconds(), //毫秒
  };
  if (/(Y+)/.test(fmt))
    fmt = fmt.replace(
      RegExp.$1,
      (this.getFullYear() + "").substr(4 - RegExp.$1.length)
    );
  for (const k in o)
    if (new RegExp("(" + k + ")").test(fmt))
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length)
      );
  return fmt;
};

const dateTool = {
  parse(dateStr: string) {
    let formatStr = "YYYY-MM-DD HH:mm:ssZ";
    if (dateStr.indexOf(".") >= 0) {
      formatStr = "YYYY-MM-DD HH:mm:ss.SSSSSSSSSZ";
    }
    return dayjs(dateStr, formatStr);
  },

  zeroTime: "0001-01-01T08:00:00+08:00",
  format: "YYYY-MM-DD HH:mm:ss",
  formatYMD: "YYYY-MM-DD",
  formatYMD2: "YYYY年MM月DD日",
  formatYMDHM: "YYYY-MM-DD HH:mm",
  formatYMDHM2: "YYYY年MM月DD日 HH点mm分ss秒",

  getReplyTime(date: string) {
    const time = this.parse(date).valueOf();
    const currentT = new Date().getTime();
    const diff = (currentT - time) / 1000;
    if (diff < 60) {
      return "刚刚";
    } else if (diff < 60 * 60) {
      return `${Math.round(diff / 60)}分钟前`;
    } else if (diff < 24 * 60 * 60) {
      return `${Math.round(diff / 60 / 60)}小时前`;
    } else if (diff < 7 * 24 * 60 * 60) {
      return `${Math.round(diff / 24 / 60 / 60)}天前`;
    } else {
      return dayjs(time, this.formatYMD);
      // return `${parseInt(diff / 365 / 24 / 60 / 60)}年前`
    }
  },
};

export default dateTool;
