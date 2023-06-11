package main

import (
	"context"
	"fmt"
	"github.com/RomaBiliak/Minesweeper/internal/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Addr: ":3000",
	}

	http.Handle("/", http.HandlerFunc(handlers.Home))
	http.Handle("/newGame", http.HandlerFunc(handlers.StartGame))
	http.Handle("/openCell", http.HandlerFunc(handlers.OpenCell))
	http.Handle("/disabledEnabledCell", http.HandlerFunc(handlers.DisabledEnabledCell))

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		fmt.Print("Server is listening...")
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
