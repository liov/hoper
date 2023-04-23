// history

export const getQueryByNameHistory = (name) => {
  return new URL(location.href).searchParams.get(name);
  // 或
  // return new URLSearchParams(location.search).get(name)
};

// hash
export const getQueryByNameHash = (name) => {
  return new URLSearchParams(location.hash.split("?")[1]).get(name);
};

export const isHashMode = location.hash !== "";

export const getQueryByName = (name) => {
  const queryList = location.href.split("?")[1]?.split("&") || [];
  const curQuery = queryList.find((item) => item.includes(name)); // 或用 filter 方法

  return curQuery?.split("=")[1] || "";
};

export const parseQueryString = function (): Record<string, string> {
  const str = location.search;
  const objURL = {};

  str.replace(
    new RegExp("([^?=&]+)(=([^&]*))?", "g"),
    ($0, $1, $2, $3) => (objURL[$1] = $3)
  );
  return objURL;
};
