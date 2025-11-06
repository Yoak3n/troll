export interface WsMessage {
    action: string
    data: any
}

export interface LogItem {
    time: string
    content: string
}

export interface TaskProcess {
    id: string
    label: string
    total: number
    current: number
    completed: boolean
}