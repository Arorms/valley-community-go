package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/config"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/migrations"
)

func main() {
	log.SetOutput(os.Stdout)

	direction := flag.String("dir", "up", "action: up, down, version or force")
	steps := flag.Int("steps", 0, "number of steps for down (0 = all)")
	force := flag.Int("force", -1, "target version for force")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("load embedded migrations: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, databaseURL(cfg.DB))
	if err != nil {
		log.Fatalf("init migrate: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("close source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("close database: %v", dbErr)
		}
	}()

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate up: %v", err)
		}
	case "down":
		var err error
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate down: %v", err)
		}
	case "version":
		printVersion(m)
		return
	case "force":
		if *force < 0 {
			log.Fatal("force requires -force <version>")
		}
		if err := m.Force(*force); err != nil {
			log.Fatalf("migrate force: %v", err)
		}
	default:
		log.Fatalf("unknown -dir %q (want up, down, version or force)", *direction)
	}

	printVersion(m)
}

func databaseURL(cfg config.DBConfig) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.Database,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

func printVersion(m *migrate.Migrate) {
	v, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Println("version: none")
			return
		}
		log.Fatalf("read version: %v", err)
	}
	log.Printf("version=%d dirty=%t", v, dirty)
}
