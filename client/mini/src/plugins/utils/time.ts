import dayjs from "dayjs";

export const date2s = (value) => dayjs(value).format("YYYY-MM-DD HH:mm:ss");
export const s2date = (value) => dayjs(value, "YYYY-MM-DD HH:mm:ss.SSS Z");
export const datefmt = (value, format) => dayjs(value).format(format);
