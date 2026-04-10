import axios from "axios";

class ArticleClient {
  static async list(pageNo: number, pageSize: number) {
    const { data } = await axios.get(
      `/api/article?pageNo=${pageNo}&pageSize=${pageSize}`
    );
    return data;
  }

  static async info(id) {
    const { data } = await axios.get(`/api/article/` + id);
    return data.data;
  }
}

export default ArticleClient;
