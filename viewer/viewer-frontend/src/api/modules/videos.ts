import request from "../../utils/request"
import type { CommentViewWithVideo } from "../../types";
const API = {
    COMMENT_LSIT: "/videos/", // + topicName + /comments
} as const;


export const fetchCommentsByVideo = (avid: number)=>{
    return request.get<any, CommentViewWithVideo>(API.COMMENT_LSIT + avid + "/comments")
}