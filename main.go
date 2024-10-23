package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/matmazurk/acc2/backup"
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
	logger.Info().Msg("starting...")

	db, err := db.New(flags.dbFilename)
	if err != nil {
		logger.Fatal().Err(err).Str("filename", flags.dbFilename).Msg("could not open db")
	}
	logger.Info().Str("filename", flags.dbFilename).Msg("database opened")

	store, err := imagestore.NewStore(flags.storeDir, logger)
	if err != nil {
		logger.Fatal().Err(err).Str("dir", flags.storeDir).Msg("could not open imagestore")
	}
	logger.Info().Str("dir", flags.storeDir).Msg("imagestore opened")

	server := &http.Server{
		Addr:    flags.httpListenAddr,
		Handler: lhttp.NewMux(db, store, logger),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info().Str("listen_addr", flags.httpListenAddr).Msg("starting http server")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("error http server listen")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		const jobHour = 3
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		now := time.Now()
		nextBackup := time.Date(now.Year(), now.Month(), now.Day(), jobHour, 0, 0, 0, time.UTC)
		if now.After(nextBackup) {
			nextBackup = nextBackup.Add(24 * time.Hour)
		}

		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("stopping cron")
				return
			case <-ticker.C:
				now := time.Now()
				if now.Before(nextBackup) {
					continue
				}

				filename := fmt.Sprintf("acc-backup-%s.zip", now.Format("2006-01-02_15:04:05"))
				logger.Info().Msgf("starting backup job, filename %s", filename)

				f, err := os.Create(filename)
				if err != nil {
					logger.Error().Err(err).Msg("could not create new backup file")
					continue
				}

				err = backup.Backup(f, ".")
				if err != nil {
					logger.Error().Err(err).Msg("could not execute backup")
					f.Close()
					os.Remove(filename)
					continue
				}
				f.Close()

				logger.Info().Msg("backup job successfully finished")

				nextBackup = nextBackup.Add(24 * time.Hour)
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	fmt.Println()

	logger.Info().Msg("shutting down http server...")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("error shutting down http server")
	}
	cancel()

	logger.Info().Msg("http server gracefully shutdown")
	wg.Wait()
	logger.Info().Msg("all finished")
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
