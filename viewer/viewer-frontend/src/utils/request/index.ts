import  axios   from "axios"
// 第一步：利用axios的create方法，去创建axios实例（其他的配置：基础路径、超时限制）
let request = axios.create({
    baseURL:import.meta.env.VITE_APP_BASE_API,
    timeout:5000,
});

request.interceptors.response.use((response)=>{
    // 简化输出
    return response.data
},(error)=>{
    // 失败的回调，处理http网络错误
    let message = ''
    let statusCode = error.response.data.code;
    switch(statusCode){
        case 401:
            message = "TOKEN过期"
            break;
        case 403:
            message = "无授权"
            break;
        case 404:
            message = "请求地址错误"
            break;
        case 500:
            message = "服务器出现问题"
            break;
        default:
            message = "网络出现问题"
            break;
    }
    // window.$message.error(message,{ duration: 2500 })
    return Promise.reject(message)
})
export default request;