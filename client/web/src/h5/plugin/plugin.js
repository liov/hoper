import dayjs from "dayjs";
import "dayjs/locale/zh-cn";

import { upload } from "./utils/upload";

dayjs.locale("zh-cn");

export default {
  install: (app, options) => {
    app.directive("datefmt", function (value) {
      return dayjs(value).format("YYYY-MM-DD HH:mm:ss");
    });

    app.config.globalProperties.$s2date = (value) =>
      dayjs(value, "YYYY-MM-DD HH:mm:ss.SSS Z");
    app.config.globalProperties.$date2s = (value) =>
      dayjs(value).format("YYYY-MM-DD HH:mm:ss");
    app.config.globalProperties.$datefmt = (value, format) =>
      dayjs(value).format(format);

    app.config.globalProperties.$customUpload = async ({
      action,
      data,
      file,
      filename,
      headers,
      onError,
      onProgress,
      onSuccess,
      withCredentials,
      classify,
    }) => {
      const res = await upload(file);
      onSuccess({ data: res, status: 200 }, file);
      file.status = "done";
    };
  },
};
