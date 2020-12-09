package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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

func startHttpServer2(group *errgroup.Group) *http.Server {
	service := &http.Server{Addr: ":8002"}

	http.HandleFunc("/s2", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello2\n")
	})

	group.Go(func() error {
		defer log.Println("Server2 exit")
		var err error
		if err = service.ListenAndServe(); err != http.ErrServerClosed {
			log.Panicf("ListenAndServe(): %v", err)
		}
		return err
	})

	return service
}

func startHttpServer3(group *errgroup.Group) *http.Server {
	service := &http.Server{Addr: ":8003"}

	http.HandleFunc("/s3", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello3\n")
	})

	group.Go(func() error {
		defer log.Println("Server3 exit")
		var err error
		if err = service.ListenAndServe(); err != http.ErrServerClosed {
			log.Panicf("ListenAndServe(): %v", err)
		}
		return err
	})

	return service
}
