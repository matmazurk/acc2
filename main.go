package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
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
)

func main() {
	flags := parseFlags()
	cleanup, err := setup(flags)
	if err != nil {
		slog.Error("could not setup", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	slog.Info("staring...")

	db, err := db.New(flags.dbFilename)
	if err != nil {
		slog.Error("could not setup db", "error", err)
		os.Exit(1)
	}
	slog.Info("database opened", "filename", flags.dbFilename)

	store, err := imagestore.NewStore(flags.storeDir)
	if err != nil {
		slog.Error("could not open imagestore", slog.String("dir", flags.storeDir), "error", err)
	}
	slog.Info("imagestore opened", slog.String("dir", flags.storeDir))

	server := &http.Server{
		Addr:    flags.httpListenAddr,
		Handler: lhttp.NewMux(db, store),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		slog.Info("starting http server", slog.String("listen_addr", flags.httpListenAddr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error http server listen", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	fmt.Println()

	slog.Info("shutting down http server...")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("error shutting down http server", "error", err)
		os.Exit(1)
	}
	cancel()

	slog.Info("http server gracefully shutdown")
	wg.Wait()
	slog.Info("all finished")
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

func setup(f flags) (func(), error) {
	var callbacks []func()

	if f.printToStdout {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})))
	} else {
		f, err := os.OpenFile("logs", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			return nil, fmt.Errorf("could not open logs file: %w", err)
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{})))
		callbacks = append(callbacks, func() { f.Close() })
	}

	return func() {
		for _, c := range callbacks {
			c()
		}
	}, nil
}

func obsoleteCron(ctx context.Context) {
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
			// logger.Info().Msg("stopping cron")
			return
		case <-ticker.C:
			now := time.Now()
			if now.Before(nextBackup) {
				continue
			}

			filename := fmt.Sprintf("acc-backup-%s.zip", now.Format("2006-01-02_15:04:05"))
			// logger.Info().Msgf("starting backup job, filename %s", filename)

			f, err := os.Create(filename)
			if err != nil {
				// logger.Error().Err(err).Msg("could not create new backup file")
				continue
			}

			err = backup.Backup(f, ".")
			if err != nil {
				// logger.Error().Err(err).Msg("could not execute backup")
				f.Close()
				os.Remove(filename)
				continue
			}
			f.Close()

			// logger.Info().Msg("backup job successfully finished")

			nextBackup = nextBackup.Add(24 * time.Hour)
		}
	}
}
