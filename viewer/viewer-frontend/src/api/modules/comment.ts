import type { CommentView } from "../../types";
import request from "../../utils/request";

const API = {
    SEARCH_RANGE: "/comments/search" // ? keyword
} as const;

export const fetchCommentsByKeyword = (keyword: string) => {
    const uri = API.SEARCH_RANGE + '?keyword=' + encodeURIComponent(keyword);
    console.log(uri);
    return request.get<any, CommentView>(uri);
}