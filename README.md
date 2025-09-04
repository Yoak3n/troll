网上冲浪，看到各种神人既是乐趣，也是速效升压药

对于二者，在需要的时候，搜寻并品鉴他们是如此刻不容缓，那就赶快行动起来吧！

## Introduction

一个使用go语言编写的cli工具,使用简单的命令收集评论数据并找出其中的巨魔



## Usage

### Install
```bash
go install github.com/Yoak3n/troll/troll@latest
```
安装cli工具成功后，帮助信息如下
```bash
troll -h

>> OUTPUT:
NAME:
   troll - search trolls from bilibili

USAGE:
   troll [global options] [command [command options]]

VERSION:
   0.0.1

COMMANDS:
   fetch    fetch comments of a topic from bilibili
   query    query something from troll
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cache string, -C string  cache path (default: "data/cache")
   --title string, -T string  specify title as directory
   --help, -h                 show help
   --version, -V              print only the version (default: false)
```
两个全局参数`cache`和`title`:

cache即存放评论信息数据的根目录，默认当前目录下的`data/cache`

title即存放一批视频评论信息数据的目录，在cache目录下，在获取和查询时都需要有相应的值以限定范围


### Prepare
首先参考`config.example.yaml`文件，准备一个`config.yaml`文件放在当前目录（使用cli工具时所在目录必须要有该文件）

获取bilibili的cookie填写好cookie字段（必填）， 如果要使用代理可以填好系统代理的地址

### Fetch
使用`fetch`子命令先获取巨魔活动的范围环境信息——单独视频或一个话题下多个视频的评论区
```bash
troll fetch -h

>> OUTPUT:
NAME:
   troll fetch - fetch comments of a topic from bilibili

USAGE:
   troll fetch [options]

OPTIONS:
   --help, -h                 show help
   --avid int, -a int         specify a video by avid (default: -1)
   --bvid string, -b string   specify a video by bvid
   --topic string, -t string  specify many video by topic name

GLOBAL OPTIONS:
   --cache string, -C string  cache path (default: "data/cache")
   --title string, -T string  specify title as directory
   --version, -V              print only the version (default: false)
```
avid,bvid和topic三个参数任选其一

当指定topic时，title未指定则与topic相同，当指定avid或bvid时，title未指定则报错

### Query
使用`query`子命令以从现有数据中筛选突出的用户或评论
```bash
troll query -h

>> OUTPUT:
NAME:
   troll query - query something from troll

USAGE:
   troll query [options]

OPTIONS:
   --help, -h               show help
   --top string, -t string  show top users or comments
   --count int              limit the count of top users or comments or words  (default: 10)
   --user string            show the comments of a user

GLOBAL OPTIONS:
   --cache string, -C string  cache path (default: "data/cache")
   --title string, -T string  specify title as directory
   --version, -V              print only the version (default: false)
```
`top`与`count`两个参数来获取出现次数最多的user或comment或word，`--top word`需要通过tfidf算法计算有价值的词条，性能开销较大，目前版本慎用
`user`获取指定用户名的所有评论

