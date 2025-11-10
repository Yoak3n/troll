export interface Setting {
    cookies: SettingItem[]
    proxies: SettingItem[]
    interval: RequestInterval
}

export interface SettingItem {
    type: string
    data: string
    id?: number
}

export interface RequestInterval {
    basic: number
    random: number
}