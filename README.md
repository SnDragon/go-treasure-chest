# go-treasure-chest
Go百宝箱,收集Go有趣项目

## Roadmap
* [短地址服务](https://github.com/SnDragon/go-treasure-chest/blob/master/docs/shorturl)

## 项目目录
参考[golang-standards/project-layout](https://github.com/golang-standards/project-layout)项目组织
```
.
├── LICENSE
├── Makefile
├── README.md
├── build
│   └── app1         // 编译后的二进制文件
├── cmd
│   └── app1         // 程序启动入口
├── configs
│   └── app1         // 配置文件
├── docs
│   └── app1         // 文档
├── pkg1             // 外部可引用包
├── go.mod
├── go.sum
└── internal         // 内部包
    ├── app1
    └── pkg1
```