# 鲲鹏

> 鲲鹏（kūn péng） Kun Peng   
> 中国古代神话传说中出现的神兽，是奇大无比的两种生物。  
> 取名为此是希望haetae项目能跑的更快


## 1. 项目背景
### 1.1 业务背景

在 [haetae](https://github.com/BetaSummer/haetae/) 项目中，数据查询跑的实在太慢了，急需Redis缓存，故建立此项目完成Redis的写任务。

### 1.2 技术栈&运行环境

```
golang
docker
canal(https://github.com/alibaba/canal)
```

---
## 2. 项目使用
### 2.1 项目启动
运行docker-compose.yml即可  
*请务必先阅读 **init.sql***

### 2.2 注意事项
现在使用的网络模式是host，如果实在无法读取到 `mysql`/`Redis` 服务, 请使用bridge连接或者直接在相应的配重中修改为内网ip进行访问

***现在还是实验阶段, Redis服务尚未编写***