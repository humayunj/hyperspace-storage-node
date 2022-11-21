package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func longOpr(ctx context.Context, ch chan<- string) {
	select {
	case <-time.After(time.Second * 3):
		ch <- "Success"
	case <-ctx.Done():
		close(ch)
	}
}
func epHandler(w http.ResponseWriter, r *http.Request) {
	notif, ok := w.(http.CloseNotifier)
	if !ok {
		panic("Error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan string)
	go longOpr(ctx, ch)
	select {
	case result := <-ch:
		fmt.Fprint(w, result)
		cancel()
		return
	case <-time.After(time.Second * 5):
		fmt.Fprint(w, "Server busy")
	case <-notif.CloseNotify():
		fmt.Println("Client discon")

	}
	cancel()
	<-ch
}

func createServer() {
	http.HandleFunc("/handle", epHandler)
}
