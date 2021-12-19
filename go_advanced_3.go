package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	group, errCtx := errgroup.WithContext(ctx)

	// 启动 HTTP Server
	srv := &http.Server{Addr: ":8080"}
	group.Go(func() error {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, World")
		})
		err := srv.ListenAndServe()
		if err != nil {
			cancel()
		}
		return err
	})

	group.Go(func() error {
		<-errCtx.Done()
		fmt.Println("http server stop")
		return srv.Shutdown(errCtx)
	})

	// 添加对信号的处理
	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-chanel:
				cancel()
			}
		}
		return nil
	})

	// Wait
	if err := group.Wait(); err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("Success")
}
