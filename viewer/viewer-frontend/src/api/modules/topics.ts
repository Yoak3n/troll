import request from "../../utils/request"
import type { TopicListResponse } from "../types"
import type {VideoDataWithCommentsCount} from "../../types"
const API = {
    TOPIC_URL: "/topics/list",
    TOPIC_VIDEOS_URL: "/topics/", // + topicName
} as const;

export const fetchTopics = ()=>{
    return request.get<any, TopicListResponse[]>(API.TOPIC_URL)
}
export const fetchVideosByTopic = (topicName: string)=>request.get<any, VideoDataWithCommentsCount[]>(API.TOPIC_VIDEOS_URL + encodeURIComponent(topicName))