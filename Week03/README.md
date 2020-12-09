# Week03

## 思路
使用`errgroup`创建多个HTTP server，并将其指针传出。

使用`signal.Notify`注册需要捕获的系统信号，接收到时，关闭所有HTTP server。

## 核心代码
```go
func main() {
	ctx, done := context.WithCancel(context.Background())
	g, _ := errgroup.WithContext(ctx)

	service1 := startHttpServer1(g)
	service2 := startHttpServer2(g)
	service3 := startHttpServer3(g)

	log.Println("Start services!")
	// 用于捕捉SIGTERM的goroutine
	g.Go(func() error {
		// 注册linux (kill -15) 信号
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			log.Printf("Received signal: %s\n", sig)
			if err := service1.Shutdown(context.Background()); err != nil {
				log.Panicf("Service1 exit fail: %v\n", err)
			}
			if err := service2.Shutdown(context.Background()); err != nil {
				log.Panicf("Service2 exit fail: %v\n", err)
			}
			if err := service3.Shutdown(context.Background()); err != nil {
				log.Panicf("Service3 exit fail: %v\n", err)
			}
			done()
		}
		return nil
	})

	// 等待所有GoRoutine终止
	err := g.Wait()

	if err != nil {
		log.Printf("Exit with : %v\n", err)
	} else {
		log.Println("Exit!")
	}
}

func startHttpServer1(group *errgroup.Group) *http.Server {
	service := &http.Server{Addr: ":8001"}

	http.HandleFunc("/s1", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello1\n")
	})

	group.Go(func() error {
		defer log.Println("Server1 exit")
		var err error
		if err = service.ListenAndServe(); err != http.ErrServerClosed {
			log.Panicf("ListenAndServe(): %v", err)
		}
		return err
	})

	return service
}
```

## 运行结果
```bash
2020/12/09 23:31:56 Start services!
2020/12/09 23:32:31 Received signal: interrupt
2020/12/09 23:32:31 Server1 exit
2020/12/09 23:32:31 Server2 exit
2020/12/09 23:32:32 Server3 exit
2020/12/09 23:32:32 Exit with : http: Server closed
```