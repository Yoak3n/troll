import type { CommentViewWithVideo,SearchFilterRequest } from '../../types';
import request from '../../utils/request'

const API = {
    USER_COMMENTS: "/users/", // + uid + '/comments'
    USER_COMMENTS_FILTER: "/users/filter/coments" //?uid ?name ?rangeType ? rangeData 
} as const;

export const fetchCommentsByUserAndTopic = (uid:number,topic?:string) => {
    let uri = `${API.USER_COMMENTS}${uid}` + '/comments'
    return request.get<any,CommentViewWithVideo[]>(topic? uri+`?topicName=${topic}`: uri)
}


export const fetchCommentsBySearchFilter = (data :SearchFilterRequest) => {
    let uri = `${API.USER_COMMENTS_FILTER}?${searchFilterRequestToQuery(data)}`
    console.log(uri);
    
    return request.get<any,CommentViewWithVideo[]>(uri)
}

function searchFilterRequestToQuery(data :SearchFilterRequest):string{
    const searchParams = new URLSearchParams();
    if (data.uid !== undefined) {
        searchParams.append('uid', data.uid.toString());
    }
    if (data.name !== undefined) {
        searchParams.append('name', data.name);
    }
    if (data.rangeData !== undefined) {
        searchParams.append('rangeData', data.rangeData);
    }
    searchParams.append('rangeType', data.rangeType);
    return searchParams.toString()
}