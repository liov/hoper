import axios from 'axios'
import mpAdapter from 'axios-miniprogram-adapter'
axios.defaults.adapter = mpAdapter

export default axios;
