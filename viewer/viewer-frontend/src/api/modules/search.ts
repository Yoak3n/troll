import request from "../../utils/request";
import type { SearchOptionRequest,SearchOptionResponse } from "../../types";
const API = {
    SEARCH_RANGE: "/search/options"
} as const;

export const fetchDataRange = (option:SearchOptionRequest) => {
    const uri = API.SEARCH_RANGE + '?' + toQueryString(option)
    console.log(uri);
    return request.get<any,SearchOptionResponse>(uri)
}

function toQueryString(params: SearchOptionRequest): string {
    const searchParams = new URLSearchParams();
    
    if (params.uid !== undefined) {
        searchParams.append('uid', params.uid.toString());
    }
    if (params.name !== undefined) {
        searchParams.append('name', params.name);
    }
    searchParams.append('rangeType', params.rangeType);
    
    return searchParams.toString();
}