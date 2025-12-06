export interface Step {
    title: string
    content: string
}

export const steps: Step[] = [
    {
        title: "在设置页面填写Cookie",
        content: "爬取评论需要填写哔哩哔哩Cookie，Cookie可以在浏览器的开发者工具中获取。"
    },
    {
        title: "在控制台页面创建任务",
        content: "在控制台创建任务需要填写话题名称或视频BV号。"
    },
    {
        title: "在控制台页面查看任务状态",
        content: "在控制台页面可以查看任务的状态。"
    },
    {
        title: "在面板页面查看任务结果",
        content: "任务完成时，在面板页面可以查看视频的评论数据。"
    }
]