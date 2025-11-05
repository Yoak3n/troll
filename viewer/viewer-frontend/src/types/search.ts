export interface SearchOptionRequest {
    uid?: number,
    name?: string,
    rangeType: string 
}

export interface SearchFilterRequest extends SearchOptionRequest {
    rangeData?: string 
} 

export interface SearchOptionResponse {
    type: string,
    options: RangeDataType[]
}

export interface RangeDataType {
    label: string,
    value: string
}