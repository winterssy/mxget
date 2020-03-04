# mxget

通过命令行在线搜索你喜欢的音乐，下载并试听。

[![Actions Status](https://img.shields.io/github/workflow/status/winterssy/mxget/Build/master?logo=appveyor)](https://github.com/winterssy/mxget/actions)

## 支持的音乐平台

|                音乐平台                 |          平台标识           |
| :-------------------------------------: | :-------------------------: |
| **[网易云音乐](https://music.163.com)** |      `netease` / `nc`       |
|     **[QQ音乐](https://y.qq.com)**      |      `tencent` / `qq`       |
| **[咪咕音乐](http://music.migu.cn/v3)** |        `migu` / `mg`        |
|  **[酷狗音乐](http://www.kugou.com)**   |       `kugou` / `kg`        |
|   **[酷我音乐](http://www.kuwo.cn/)**   |        `kuwo` / `kw`        |
|  **[虾米音乐](https://www.xiami.com)**  |       `xiami` / `xm`        |
| **[千千音乐](http://music.taihe.com)**  | `qianqian` / `baidu` / `bd` |

## 下载安装

```sh
go get -u github.com/winterssy/mxget
```

## 使用帮助

```
 _____ ______      ___    ___ ________  _______  _________   
|\   _ \  _   \   |\  \  /  /|\   ____\|\  ___ \|\___   ___\ 
\ \  \\\__\ \  \  \ \  \/  / | \  \___|\ \   __/\|___ \  \_| 
 \ \  \\|__| \  \  \ \    / / \ \  \  __\ \  \_|/__  \ \  \  
  \ \  \    \ \  \  /     \/   \ \  \|\  \ \  \_|\ \  \ \  \ 
   \ \__\    \ \__\/  /\   \    \ \_______\ \_______\  \ \__\
    \|__|     \|__/__/ /\ __\    \|_______|\|_______|   \|__|
                  |__|/ \|__|                                

A simple tool that help you search and download your favorite music,
please visit https://github.com/winterssy/mxget for more detail.

Usage:
  mxget [command]

Available Commands:
  album       Fetch and download album's songs via its id
  artist      Fetch and download artist's hot songs via its id
  config      Specify the default behavior of mxget
  help        Help about any command
  playlist    Fetch and download playlist's songs via its id
  search      Search songs from the specified music platform
  serve       Run mxget as an API server
  song        Fetch and download single song via its id

Flags:
  -h, --help      help for mxget
      --version   version for mxget

Use "mxget [command] --help" for more information about a command.
```

- 搜索歌曲

```sh
$ mxget search --from nc -k Faded
```

- 下载歌曲

```sh
$ mxget song --from nc --id 36990266
```

- 下载专辑

```sh
$ mxget album --from nc --id 3406843
```

- 下载歌单

```sh
$ mxget playlist --from nc --id 156934569
```

- 下载歌手热门歌曲

```sh
$ mxget artist --from nc --id 1045123
```

- 自动更新音乐标签/下载歌词

如果你希望 `mxget` 为你自动更新音乐标签，可使用 `--tag` 指令，如：

```sh
$ mxget song --from nc --id 36990266 --tag
```

当使用 `--tag` 指令时，`mxget` 会同时将歌词内嵌到音乐文件中，一般而言你无须再额外下载歌词。如果你确实需要 `.lrc` 格式的歌词文件，可使用 `--lyric` 指令，如：

```sh
$ mxget song --from nc --id 36990266 --lyric
```

- 设置默认下载目录

默认情况下，`mxget` 会下载音乐到当前目录下的 `downloads` 文件夹，如果你想要更改此行为，可以这样做：

```sh
$ mxget config --dir <directory>
```

>  `directory` 必须为绝对路径。

- 设置默认音乐平台

`mxget` 默认使用的音乐平台为网易云音乐，你可以通过以下命令更改：

```sh
$ mxget config --from qq
```

这样，如果你不通过 `--from` 指令指定音乐平台，`mxget` 便会使用默认值。

在上述命令中，你会经常用到 `--from` 以及 `--id` 这两个指令，它们分别表示音乐平台标识和音乐id。

> 音乐id为音乐平台为对应资源分配的唯一id，当使用 `mxget` 进行搜索时，歌曲id会显示在每条结果的后面。你也可以通过各大音乐平台的网页版在线搜索相关资源，然后从结果详情页的URL中获取其音乐id。值得注意的是，酷狗音乐对应的歌曲id即为文件哈希 `hash` 。

- 多任务下载

`mxget` 支持多任务快速并发下载，你可以通过 `--limit` 参数指定同时下载的任务数，如不指定默认为CPU核心数。

```sh
$ mxget playlist --from nc --id 156934569 --limit 16
```

> `mxget` 允许设置的最高并发数是32，但使用时建议不要超过16。

## 免责声明

- 本项目仅供学习研究使用。
- 本项目使用的接口如无特别说明均为官方接口，音乐版权归源音乐平台所有，侵删。

## License

GPLv3。
