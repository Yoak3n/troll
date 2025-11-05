export const jumpToReply = (bvid: string, commentId: number)=>{
    const uri = `https://www.bilibili.com/video/${bvid}#reply${commentId}`
    window.open(uri,"_blank")
}
export const jumpToRelyWithAvid = (avid: number, commentId: number) => {
    // TODO 之后用到了再实现
    const bvid = avid.toString()
    jumpToReply(bvid,commentId)
}