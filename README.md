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

## 加载方式

```
import (
    "github.com/hirpc/hrpc"
    "github.com/hirpc/hrpc/database/mysql"
    "github.com/hirpc/hrpc/database/redis"
    // "github.com/hirpc/hrpc/life"
    "github.com/hirpc/hrpc/option"
    
    pb "...."
)
func main() {
    // 需要携带配置
    s, err := hrpc.NewServer(
    	// 在TKE环境下的容器，可以忽略这个Consul部分的配置
        option.WithConsul(option.Consul{
            // Token 如果环境变量未配置，需要在此指定，否则连不上consul
	    Token: "xxxxxx",
            // Address 如果环境变量未配置，需要在此指定，否则连不上consul
            Address: "xxx",
            // DataCenter 默认数据中心是dc1，如果不一样，可以配置环境变量或者在此指定
            DataCenter: "",
	}),
        // 开启mysql和redis的支持
	option.WithDatabases(mysql.New(), redis.New()),
        // 定义该服务的名称
	option.WithServerName("服务名称"),
        // 定义环境，当前分为option.Development开发环境与option.Production生产环境
	option.WithEnvironment(option.Development),
        // 如果开启健康检查，则会启动http服务监听本地6688端口用于consul的http请求探测
	option.WithHealthCheck(),
    )
    if err != nil {
        panic(err)
    }
    pb.RegisterAuthenticationCenterServiceServer(s.Server(), &authenticationCenterServerImpl{})

    // 如果有任何需要的析构函数，可以在`s.Serve()`执行前，进行注册
    // life.WhenExit(xxxx.Destory)

    if err := s.Serve(); err != nil {
        panic(err)
    }
}
```

## 服务间调用

```
import (
    "github.com/hirpc/hrpc/service"
)

// Example 会调用consul获取对应健康状态下的服务
// 如果有多个实例，则会基于weight排序，返回weight最大的那个服务
// 注意：
//   返回的服务包含IP、名称、权重等等信息，不推荐在TKE环境中，使用Endpoint属性(POD IP)作为连接方式，因为POD IP重启变化
//   可以使用Target()属性，返回`服务名:端口`的字符串，利用TKE的DNS解析出服务IP，由TKE转发至POD中
func Example() error {
    // 默认是基于当前环境的设置（在NewServer阶段传入的option.Development或者option.Production），获取当前环境下的服务
    // 但支持传入tag标识，获取其他环境下的服务
    // 假设当前在开发环境，如果你的consul的Token允许的话，那么可以`service.Get("服务名称", service.Tag(option.Production))`获取生产环境下的服务
    p, err := service.Get("服务名称")
	if err != nil {
		return err
	}
	c, err := grpc.Dial(p.Target(), grpc.WithInsecure())
	if err != nil {
		return err
	}
    ...
    return nil
}

```

## 数据库访问

前提需要在main.go的`hrpc.NewServer`中，传入`option.WithDatabase(mysql.New(), redis.New()),`选项。
之后在其他子包中，可直接采用如下方式访问
```
import (
    "github.com/hirpc/hrpc/database/mysql"
    "github.com/hirpc/hrpc/database/redis"
)

func UserInfo(uid string) {
    // MySQL的访问
    var v string
	if err := mysql.Client().QueryRow(ctx, `
	SELECT
		uid
	FROM
		users_info
	WHERE
		uid = ? AND is_admin = 1 AND status = 0`, []interface{}{&v}, uid,
	); err != nil {
		if err != sql.ErrNoRows 
		}
		return err
	

    // Redis的访问
    redis.Get().Set(xxx)
}
```
因为，在`hrpc.NewServer`阶段，会基于配置，维护MySQL或者Redis的连接，可以用对应包的Get()方法获取链接。
- 如果与MySQL或者Redis建立连接失败，会在启动阶段panic
- 如果中途连接断掉，会自动尝试重连；如果对应服务挂掉，则会直接error

## 框架基本流程

1. 服务启动，先规整配置信息并且有些依赖组件需要初始化，比如`uniqueid`
2. 建立与consul的连接
3. 对pb文件注册
4. (可能)根据配置情况，建立于mysql或者redis的连接
   1. 相对应的，会注册close函数
5. 尝试宣告获取本地8888端口
6. 向consul注册服务
7. goroutine启动grpc接口
8. 注册一些析构函数
9. 阻塞监听退出等信号

## 有效环境变量

- `CONFIGS_TOKEN` 配置中心token
- `CONFIGS_ADDR` 配置中心请求地址
- `CONFIGS_DATACENTER` 配置中心的数据中心

## 插件列表-

- `hrpc-location` 用于实现常用时区信息
- `hrpc-configs` 用于对接consul

## TODO:
