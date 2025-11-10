import request from "../../utils/request";
import type{ Setting } from "../../types";
const API = {
    CONFIG_API: "/setting/"
} as const;

export const fetchSetting = ()=> request.get<any,Setting>(API.CONFIG_API)
export const updateSetting = (item:Setting) => request.post(API.CONFIG_API,item)

