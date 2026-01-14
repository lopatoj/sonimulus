package main

import (
	"flag"
	"log/slog"

	"lopa.to/sonimulus/config"
	"lopa.to/sonimulus/repository"
	"lopa.to/sonimulus/scraper"
)

func main() {
	rootHandle := flag.String("handle", "dxmfromcvs", "root user to perform bfs from")
	depth := flag.Int("depth", 0, "max bfs depth")
	flag.Parse()

	// Load config struct from environment variables and program arguments
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to initialize config", "error", err)
		return
	}

	scraper := scraper.NewScraper(*depth, cfg)
	db, err := repository.NewDB(cfg)
	if err != nil {
		slog.Error("failed to initialize database connection", "error", err)
		return
	}
	repo := repository.NewPeopleRepository(db)

	user := *rootHandle

	onPerson := func(handle, name, imageUrl string, verified bool, plan repository.Plan, trackCount int64) (id int64) {
		person, err := repo.CreatePerson(handle, name, imageUrl, verified, plan, trackCount)
		if err != nil {
			slog.Error("error creating person in people table", "error", err)
			return -1
		}
		if person == nil {
			slog.Error("person empty")
			return -1
		}
		return person.Id
	}

	onFollows := func(followerId int64, followeeHandles []string) {
		err := repo.CreateFollows(followerId, followeeHandles)
		if err != nil {
			slog.Error("error creating follows", "error", err)
		}
	}

	scraper.ScrapePeopleConcurrent(10, user, onPerson, onFollows)
}
