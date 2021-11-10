# hrpc

## 基本介绍
此框架当前主要用于特殊项目使用。

## 依赖
- 腾讯云的容器服务（及其Istio做的Mesh）
- Consul做配置管理与服务管理

## 基本架构
1. 腾讯云多节点构建的容器集群，使用istio的sidecar模式，加入Mesh依赖。
2. 部署consul于当前容器集群。得到集群内服务名。
3. 编写其他模块服务（引用此框架），开启grpc接口。（如果选项中开启了健康检查，那么镜像的6688端口也会被打开，用于http请求的健康检查）
4. 部署服务于容器集群（设置环境变量，表示consul地址、Token等信息）

## 使用方式
```
import (
    "github.com/hirpc/hrpc"
    "github.com/hirpc/hrpc/database/mysql"
	"github.com/hirpc/hrpc/database/redis"
    "github.com/hirpc/hrpc/option"
    pb "...."
)
func main() {
    // 需要携带配置
    server, err := hrpc.NewServer(
        option.WithConsul(option.Consul{
			Token: "xxxxxx",
		}),
        // 开启mysql和redis的支持
		option.WithDatabase(mysql.New(), redis.New()),
		option.WithServerName("服务名称"),/
        // 定义环境
		option.WithEnvironment(option.Development),
		option.WithHealthCheck(true),
    )
    if err != nil {
        panic(err)
    }
    pb.RegisterAuthenticationCenterServiceServer(server.Server(), &authenticationCenterServerImpl{})
    if err := server.Serve(); err != nil {
        panic(err)
    }
}
```

## 框架基本基本流程
1. 服务启动，先规整配置信息并且有些依赖组件需要初始化，比如`uniqueid`
2. 建立与consul的连接
3. 对pb文件注册
4. (可能)根据配置情况，建立于mysql或者redis的连接
   1. 相对应的，会注册close函数
5. 尝试宣告获取本地8888端口
6. 向consul注册服务
7. goroutine启动grpc接口
8. 注册一些析构函数
9.  阻塞监听退出等信号

## 有效环境变量
- `CONFIGS_TOKEN` 配置中心token
- `CONFIGS_ADDR` 配置中心请求地址
- `CONFIGS_DATACENTER` 配置中心的数据中心

## TODO: