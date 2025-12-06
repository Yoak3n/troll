import request from "../../utils/request"
import type { TopicListResponse } from "../types"
import type {VideoDataWithCommentsCount} from "../../types"
const API = {
    TOPIC_URL: "/topics/list",
    TOPIC_VIDEOS_URL: "/topics/", // + topicName
    TOPIC_UPDATE_URL: "/topics/update",
    TOPIC_DELETE_URL: "/topics/", // + topicName
} as const;

export const fetchTopics = ()=>{
    return request.get<any, TopicListResponse[]>(API.TOPIC_URL)
}
export const fetchVideosByTopic = (topicName: string)=>request.get<any, VideoDataWithCommentsCount[]>(API.TOPIC_VIDEOS_URL + encodeURIComponent(topicName) + "/videos")
export const updateTopic = (topic: string, newTopic: string)=>request.post<any, string>(API.TOPIC_UPDATE_URL, {topic, new_topic: newTopic})
export const deleteTopic = (topicName: string)=>request.delete<any, string>(API.TOPIC_DELETE_URL + encodeURIComponent(topicName))