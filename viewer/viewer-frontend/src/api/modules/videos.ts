import request from "../../utils/request"
import type { CommentViewWithVideo, TopicUpdateRequest } from "../../types";
const API = {
    COMMENT_LSIT: "/videos/", // + topicName + /comments
    TOPIC_UPDATE: "/videos/topic/", // + {avid + topic}
    DELETE_VIDEOS: "/videos/", // + {avidList}

} as const;

export const fetchCommentsByVideo = (avid: number)=> request.get<any, CommentViewWithVideo>(API.COMMENT_LSIT + avid + "/comments")
export const updateTopicOfVideos = (data: TopicUpdateRequest)=>request.post<any, {message: string}>(API.TOPIC_UPDATE, data)
export const deleteVideos = (avidList: number[])=>request.delete<any, {message: string}>(API.DELETE_VIDEOS, {data: {avidList}})
