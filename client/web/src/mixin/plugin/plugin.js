import dayjs from "dayjs";
import "dayjs/locale/zh-cn";

import { upload } from "../../utils/upload";

dayjs.locale("zh-cn");

export default {
  install: (app, options) => {
    app.directive("datefmt", function (value) {
      return dayjs(value).format("YYYY-MM-DD HH:mm:ss");
    });

    app.config.globalProperties.$toDayjs = (value) =>
      dayjs(value, "YYYY-MM-DD HH:mm:ss.SSSZ");
    app.config.globalProperties.$dateFmtDateTime = (value) =>
      dayjs(value).format("YYYY-MM-DD HH:mm:ss");
    app.config.globalProperties.$dateFmt = (value, format) =>
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
