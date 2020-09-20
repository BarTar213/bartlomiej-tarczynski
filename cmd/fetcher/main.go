package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BarTar213/bartlomiej-tarczynski/api"
	"github.com/BarTar213/bartlomiej-tarczynski/config"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
)

func main() {
	configFile := flag.String("fetcher-config", "fetcher.yml", "name of yml file with fetcher config")
	conf := config.NewConfig(configFile)
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("%+v\n", conf)

	postgres, err := storage.NewPostgres(&conf.Postgres)
	if err != nil {
		logger.Fatalf("Postgres error: %s", err)
	}

	a := api.NewApi(
		api.WithConfig(conf),
		api.WithLogger(logger),
		api.WithStorage(postgres),
		api.WithWorker(),
	)

	go a.Run()
	logger.Print("started app")

	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal
	a.Worker.Stop()
	logger.Print("exited from app")
}
