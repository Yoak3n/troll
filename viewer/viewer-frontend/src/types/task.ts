export interface TaskData {
    id: string
    type : string
    data: string[]
    topic: string
    page: number
}

export type TaskForm = Omit<TaskData, 'id'>