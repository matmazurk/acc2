package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matmazurk/acc2/db"
	lhttp "github.com/matmazurk/acc2/http"
	"github.com/matmazurk/acc2/imagestore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const listenAddr = ":80"

func main() {
	flags := parseFlags()
	cleanup := setup(flags)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())

	log.Info().Msg("starting...")
	db, err := db.New("exps.db")
	if err != nil {
		panic(err)
	}
	store, err := imagestore.NewStore(".")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:    listenAddr,
		Handler: lhttp.NewMux(db, store),
	}

	go func() {
		log.Info().Str("listen_addr", listenAddr).Msg("starting http server")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("error http server listen")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println()

	log.Info().Msg("shutting down http server...")
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("error shutting down http server")
	}

	log.Info().Msg("http server gracefully shutdown")

	cancel()
}

type flags struct {
	printToStdout bool
}

func parseFlags() flags {
	f := flags{}
	flag.BoolVar(&f.printToStdout, "s", false, "print output to stdout")
	flag.Parse()

	return f
}

func setup(f flags) func() {
	var callbacks []func()

	if f.printToStdout {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		f, err := os.OpenFile("logs", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			panic(err)
		}
		log.Logger = log.Output(f)
		callbacks = append(callbacks, func() { f.Close() })
	}

	return func() {
		for _, c := range callbacks {
			c()
		}
	}
}
