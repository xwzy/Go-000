package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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

func main() {

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	// 用于捕捉SIGTERM的goroutine
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
			done()
		case <-gctx.Done():
			fmt.Printf("closing signal goroutine\n")
			return gctx.Err()
		}

		return nil
	})

	service1 := startHttpServer1(g)
	service2 := startHttpServer2(g)
	service3 := startHttpServer3(g)

	// 进行
	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				fmt.Printf("Receiving done\n")
				service1.Shutdown(context.Background())
				service2.Shutdown(context.Background())
				service3.Shutdown(context.Background())
				return gctx.Err()
			}
		}
	})

	// 等待所有GoRoutine终止
	err := g.Wait()

	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Println("Exit all http server gracefully!")
		} else {
			fmt.Printf("Exit with error: %v\n", err)
		}
	} else {
		fmt.Println("Exit!")
	}
}
