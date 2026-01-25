网上冲浪，看到各种神人既是乐趣，也是速效升压药

对于二者，在需要的时候，搜寻并品鉴他们是如此刻不容缓，那就赶快行动起来吧！

## Introduction

一个使用go语言编写的cli工具,使用简单的命令收集评论数据并找出其中的巨魔



## Usage

### Install
下载源码后使用`scripts`目录下的编译脚本（内嵌了网页文件）获得二进制程序，放在Path环境变量目录下

`build.ps1`->Windows | `build.sh`->Linux

或从[Release](https://github.com/Yoak3n/troll/releases)直接下载

安装cli工具成功后，帮助信息如下
```bash
troll -h

>> OUTPUT:
NAME:
   troll - search trolls from bilibili

USAGE:
   troll [global options] [command [command options]]

VERSION:
   0.3.0

COMMANDS:
   fetch    fetch comments of a topic from bilibili
   query    query something from troll
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cache string, -C string  cache path(Deprecated) (default: "%UserConfigDir%/troll/data/cache")
   --title string, -T string  specify title as directory
   --help, -h                 show help
   --version, -V              print only the version (default: false)
```
两个全局参数`cache`和`title`:

cache即存放评论信息数据的根目录，默认当前目录下的`data/cache`

title即存放一批视频评论信息数据的目录，在cache目录下，在获取和查询时都需要有相应的值以限定范围


### Prepare
获取bilibili的cookie，使用【troll config --cookie】填写好cookie字段（必填），由于cookie字符串过长，请使用双引号包裹
如果要使用代理可以填好系统代理的地址

### View
最低学习成本的用法，打开webui直接进行设置和获取数据
```bash
troll view
```

### Fetch
使用`fetch`子命令先获取巨魔活动的范围环境信息——单独视频或一个话题下多个视频的评论区
```bash
troll fetch -h

>> OUTPUT:
NAME:
   troll - search trolls from bilibili

USAGE:
   troll [global options] [command [command options]]

VERSION:
   0.3.0

COMMANDS:
   fetch    fetch comments of a topic from bilibili
   query    query something from troll
   config   Set config
   view     open data view
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cache string, -C string  cache path(Deprecated) (default: "%UserConfigDir%/troll/data/cache")
   --title string, -T string  specify title as directory
   --help, -h                 show help
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
   --cache string, -C string  cache path(Deprecated) (default: "%UserConfigDir%/troll/data/cache")
   --title string, -T string  specify title as directory
   --version, -V              print only the version (default: false)
```
`top`与`count`两个参数来获取出现次数最多的user或comment或word，`--top word`需要通过tfidf算法计算有价值的词条，性能开销较大，目前版本慎用
`user`获取指定用户名的所有评论

### Config
为程序设置Cookie和代理地址

```bash
NAME:
   troll config - Set config

USAGE:
   troll config [options]

OPTIONS:
   --cookie string, -c string  Set cookie
   --proxy string, -p string   Set proxy
   --help, -h                  show help

GLOBAL OPTIONS:
   --cache string, -C string  cache path(Deprecated) (default: "%UserConfigDir%/troll/data/cache")
   --title string, -T string  specify title as directory
   --version, -V              print only the version (default: false)
```
