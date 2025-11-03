import type { User } from "./user";
export interface VideoDataWithCommentsCount {
    avid: number
    bvid: string
    title: string
    topic: string
    description: string
    count: number
    author: User
}