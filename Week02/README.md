
## Week02 作业题目：
1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

代码见`homework`目录

## 解决方式
对于业务相关的“错误”，如：用户不存在，不需要提供stacktrace，可以记录日志返回404或者降级处理。其他db驱动返回的error可以wrap后向上抛出。

## 测试用例
### http://127.0.0.1:8080/user/aaa
> 错误的参数类型，返回400
```
2020/11/29 21:08:09 [Info] service.go:28 strconv.Atoi: parsing "aaa": invalid syntax
invalid user id
main.(*Service).handleGetUserNameByID
        C:/Users/wzy/GolandProjects/Go-000/Week02/homework/service.go:27
github.com/gin-gonic/gin.(*Context).Next
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/context.go:161
github.com/gin-gonic/gin.RecoveryWithWriter.func1
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/recovery.go:83
github.com/gin-gonic/gin.(*Context).Next
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/context.go:161
github.com/gin-gonic/gin.LoggerWithConfig.func1
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/logger.go:241
github.com/gin-gonic/gin.(*Context).Next
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/context.go:161
github.com/gin-gonic/gin.(*Engine).handleHTTPRequest
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/gin.go:409
github.com/gin-gonic/gin.(*Engine).ServeHTTP
        C:/Users/wzy/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/gin.go:367
net/http.serverHandler.ServeHTTP
        C:/Go/src/net/http/server.go:2843
net/http.(*conn).serve
        C:/Go/src/net/http/server.go:1925
runtime.goexit
        C:/Go/src/runtime/asm_amd64.s:1374
[GIN] 2020/11/29 - 21:08:09 |?[90;43m 400 ?[0m|      2.9799ms |       127.0.0.1 |?[97;44m GET     ?[0m "/user/aaa"
```

### http://127.0.0.1:8080/user/12345
> 不存在的用户id，返回404，并记录日志。
```
2020/11/29 21:09:28 [Info] service.go:37 sql: no rows in result set
no user with id 12345
[GIN] 2020/11/29 - 21:09:28 |?[90;43m 404 ?[0m|      6.9984ms |       127.0.0.1 |?[97;44m GET     ?[0m "/user/12345"
```

### http://127.0.0.1:8080/user/1
> 正确的请求
```
[GIN] 2020/11/29 - 21:09:51 |?[97;42m 200 ?[0m|      2.0036ms |       127.0.0.1 |?[97;44m GET     ?[0m "/user/1"
```


## Reference：
* https://github.com/gin-gonic/gin
* https://golang.org/src/database/sql/example_test.go

## 学习笔记

### 错误处理思想
> 只处理一次错误

### 错误处理方式
1. 可重用的底层包，只需要返回基本错误信息。
2. 调用基础库的时候，不处理错误，wrap之后向上抛出。
3. 服务降级的时候，可以吞掉错误，返回一个默认结果。

### 什么库是基础库？
1. 调用Go语言库
2. 调用第三方库
3. 调用自己的kit库

### 何时直接传递？
1. 自己的业务中相互调用


