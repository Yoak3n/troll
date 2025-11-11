import type { User } from "./user";
export interface VideoDataWithCommentsCount extends VideoData {
    count: number
    update_at: string
}
export interface VideoData {
    avid: number
    bvid: string
    title: string
    topic: string
    description: string
    author: User
}

export interface VideoDataGroupyByTopic{
    topic: string
    videos: VideoData[]
}