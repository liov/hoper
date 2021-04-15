import dayjs from "dayjs";

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
