import axios from "axios";
import { showFailToast } from "vant";


export async function momentList(
  pageNo: number,
  pageSize: number
){
  try {
    const { response, status } = await axios.get(`/api/moment?pageNo=${pageNo}&pageSize=${pageSize}`);
    console.log(status);
    return response;
  } catch (e) {
    showFailToast(decodeURI(e.message));
    return Promise.reject(e);
  }
}
