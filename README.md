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
*刚启动时go-kenpeng出现多次exit with code 1属于正常现象，因为canal启动需要时间，client连不上canal就会自动退出重启*

### 2.2 说明事项
现在canal client使用的是轮询方式，具体的间隔在config相应的项中修改  
如将来对延迟有需求，可考虑修改为基于kafka/RocketMQ的实现

### 2.3 注意事项
现在使用的网络模式是host，如果实在无法读取到 `mysql`/`Redis` 服务, 请使用bridge连接或者直接在相应的配重中修改为内网ip进行访问  
canal有两个release文件: deployer是后端连接mysql用；admin是web UI配置用，直接启动`startup.sh`即可使用  
***docker host网络模式仅在Linux下有效，请勿在Windows/MacOS下测试***