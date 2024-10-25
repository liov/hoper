import axios from "axios";
import SparkMD5 from "spark-md5";

const exist = async function (md5, size) {
  const res = await axios.get(`/api/v1/exists/${md5}/${size}`);
  if (res.data.code === 0) {
    return res.data.data;
  } else {
    return null;
  }
};

const getBase64 = function (img, callback) {
  const reader = new FileReader();
  reader.addEventListener("load", () => callback(reader.result));
  reader.readAsDataURL(img);
};

const getMD5 = function (file) {
  return new Promise(function (resolve, reject) {
    const blobSlice =
      File.prototype.slice ||
      File.prototype.mozSlice ||
      File.prototype.webkitSlice;

    const chunkSize = 2097152; // Read in chunks of 2MB
    const chunks = Math.ceil(file.size / chunkSize);
    let currentChunk = 0;
    const spark = new SparkMD5();
    const fileReader = new FileReader();

    fileReader.onload = function (e) {
      // console.log('read chunk nr', currentChunk + 1, 'of', chunks)
      spark.append(e.target.result); // Append array buffer
      currentChunk++;

      if (currentChunk < chunks) {
        loadNext();
      } else {
        // console.log('finished loading')
        resolve(spark.end()); // Compute hash
      }
    };
    fileReader.onerror = function () {
      reject(new Error("oops, something went wrong."));
    };
    function loadNext() {
      const start = currentChunk * chunkSize;
      const end =
        start + chunkSize >= file.size ? file.size : start + chunkSize;

      fileReader.readAsArrayBuffer(blobSlice.call(file, start, end));
    }
    loadNext();
  });
};

const upload = async function ($file) {
  const md5 = await getMD5($file);
  const existFile = await exist(md5, $file.size);

  if (existFile !== null) {
    return existFile;
  }

  // 第一步.将图片上传到服务器.
  const formdata = new FormData();
  formdata.append("file", $file);
  const res = await axios({
    url: `/api/v1/upload/${md5}`,
    method: "post",
    data: formdata,
    headers: { "Content-Type": "multipart/form-data" },
  });
  if (res.data.code === 0) {
    return res.data.data;
  } else {
    return null;
  }
};

const isExist = async function (file) {
  const md5 = await getMD5(file);
  const existUrl = await exist(md5, file.size);
  if (existUrl) {
    return existUrl;
  } else {
    return null;
  }
};

export { upload, getBase64, getMD5, isExist };
