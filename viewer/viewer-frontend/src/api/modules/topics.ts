import request from "../../utils/request"
import type { TopicListResponse } from "../types"
const API = {
    TOPIC_URL: "/topics/",
} as const;

export const fetchTopics = async()=>{
    return await request.get<any, TopicListResponse[]>(API.TOPIC_URL)
}