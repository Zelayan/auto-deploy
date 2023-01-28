## 需求一

从git或svn上拉取最新的代码，并将其编译成go的二进制可运行文件
从docker仓库中拉需要的镜像
在镜像的基础上创建容器，包括配置容器的一些参数等等
启动容器
这样当我们需要发布一个项目的新版本时直接运行这个程序就能做到一键发布。一个容器运行时，就像一个操作系统运行一样，也有崩溃的时候，此时我们需要一个监听docker容器的健康状况来以防一些意外
## 需求二

监听docker容器运行时的相关参数
针对获取到的参数做出相应的处理，如mem使用打到80%时发送邮件通知小组的开发人员
在docker容器崩溃时能重新启动该容器
假设需求

### 现在我就上面介绍的两个需求简单综合一下，以完成一个自己的需求

假设本地已有我们需要的docker image

检查docker container中是否已存在目标容器

若有，则跳转到第5步

若没有，创建一个从container

启动该容器并按时检查该container的状态

若该container已崩溃，那么该程序能自动重启container

附：我们所期望的container内部还挂在了一个宿主机的目录

以上就是本篇文章将要实现的功能
