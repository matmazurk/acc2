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

func main() {
	flags := parseFlags()
	logger, cleanup := setup(flags)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Info().Msg("starting...")

	db, err := db.New(flags.dbFilename, logger)
	if err != nil {
		log.Fatal().Err(err).Str("filename", flags.dbFilename).Msg("could not open db")
	}
	log.Info().Str("filename", flags.dbFilename).Msg("database opened")

	store, err := imagestore.NewStore(flags.storeDir, logger)
	if err != nil {
		log.Fatal().Err(err).Str("dir", flags.storeDir).Msg("could not open imagestore")
	}
	log.Info().Str("dir", flags.storeDir).Msg("imagestore opened")

	server := &http.Server{
		Addr:    flags.httpListenAddr,
		Handler: lhttp.NewMux(db, store, logger),
	}

	go func() {
		log.Info().Str("listen_addr", flags.httpListenAddr).Msg("starting http server")
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
	ctx, scancel := context.WithTimeout(ctx, 5*time.Second)
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("error shutting down http server")
	}
	scancel()

	log.Info().Msg("http server gracefully shutdown")
}

type flags struct {
	printToStdout  bool
	httpListenAddr string
	dbFilename     string
	storeDir       string
}

func parseFlags() flags {
	f := flags{}

	flag.BoolVar(&f.printToStdout, "s", false, "print output to stdout")
	flag.StringVar(&f.dbFilename, "db", "exps.db", "expenses database filename")
	flag.StringVar(&f.httpListenAddr, "httpaddr", ":80", "http server listen address")
	flag.StringVar(&f.storeDir, "store", ".", "imagestore directory")

	flag.Parse()

	return f
}

func setup(f flags) (zerolog.Logger, func()) {
	var callbacks []func()
	var logger zerolog.Logger

	if f.printToStdout {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		f, err := os.OpenFile("logs", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			log.Fatal().Err(err).Msg("could not open 'logs' file")
		}
		logger = log.Output(f)
		callbacks = append(callbacks, func() { f.Close() })
	}

	return logger, func() {
		for _, c := range callbacks {
			c()
		}
	}
}
