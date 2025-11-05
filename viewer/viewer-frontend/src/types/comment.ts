import type { User } from './user'
import type { VideoData } from './video'



export interface CommentView {
    id: number
    content: string
    owner: User
    children: CommentView[]
}

export interface CommentViewWithVideo extends VideoData {
    comments: CommentView[]
}