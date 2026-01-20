import type { TaskForm } from "../../types"
import request from "../../utils/request"

export const refreshVideoDataTask = (taskData: TaskForm) => {
    return request.post("/task/video/refresh", taskData)
}