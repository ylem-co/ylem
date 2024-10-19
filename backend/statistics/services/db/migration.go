package db

import (
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source"
	log "github.com/sirupsen/logrus"
)

func NewMigrator() (*migrate.Migrate, error) {
	gormDb, err := Instance()
	if err != nil {
		log.Debug("Migrator instance creation failed: " + err.Error())
		return nil, err
	}

	db, err := gormDb.DB()
	if err != nil {
		log.Debug("Migrator instance creation failed: " + err.Error())
		return nil, err
	}

	driver, err := clickhouse.WithInstance(db, &clickhouse.Config{})
	if err != nil {
		log.Debug("Migrator instance creation failed: " + err.Error())
		return nil, err
	}

	source.Register("golang", Driver{})

	m, err := migrate.NewWithDatabaseInstance(
		"golang://migrations",
		"clickhouse",
		driver,
	)

	if err != nil {
		log.Debug("Migrator instance creation failed: " + err.Error())
		return nil, err
	}

	return m, nil
}

/**
CREATE TABLE migrations (id String, PRIMARY KEY(id)) Engine=MergeTree;
*/

// type Migrator struct {
// 	migrators []*gormigrate.Gormigrate
// }

// func (m *Migrator) Migrate() error {
// 	for _, im := range m.migrators {
// 		err := im.Migrate()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func NewMigrator() (*Migrator, error) {
// 	db, err := Instance()
// 	if err != nil {
// 		log.Debug("Migrator instance creation failed: " + err.Error())
// 		return nil, err
// 	}
// 	opts := &gormigrate.Options{
// 		TableName:                 "migrations",
// 		IDColumnName:              "id",
// 		IDColumnSize:              255,
// 		UseTransaction:            true,
// 		ValidateUnknownMigrations: false,
// 	}

// 	sortedMigrations := sortedMigrations()
// 	migrators := make([]*gormigrate.Gormigrate, len(sortedMigrations))
// 	for k := range sortedMigrations {
// 		migrators[k] = gormigrate.New(db, opts, sortedMigrations[k:k+1])
// 	}

// 	m := &Migrator{
// 		migrators: migrators,
// 	}

// 	return m, nil
// }

type RawMigration struct {
	ID        uint
	UpQuery   string
	DownQuery string
}

type RawMigrations []*RawMigration

func (rm RawMigrations) findByID(ID uint) (int, *RawMigration) {
	for k, v := range rm {
		if v.ID == ID {
			return k, v
		}
	}

	return -1, nil
}

func (a RawMigrations) Len() int           { return len(a) }
func (a RawMigrations) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RawMigrations) Less(i, j int) bool { return a[i].ID < a[j].ID }

var migrations RawMigrations = make(RawMigrations, 0)

func AddMigration(ID uint, upQuery string, downQuery string) {
	migrations = append(migrations, &RawMigration{
		ID:        ID,
		UpQuery:   upQuery,
		DownQuery: downQuery,
	})
}

func sortedMigrations() {
	sort.Sort(migrations)
}

type Driver struct {
	migrations RawMigrations
}

func (d Driver) Open(url string) (source.Driver, error) {
	sortedMigrations()
	d.migrations = migrations
	return d, nil
}
func (d Driver) Close() error {
	return nil
}
func (d Driver) First() (version uint, err error) {
	return d.migrations[0].ID, nil
}
func (d Driver) Prev(version uint) (prevVersion uint, err error) {
	k, _ := d.migrations.findByID(version)
	if k == -1 {
		return 0, errors.New("previous version not found")
	}

	return d.migrations[k-1].ID, nil
}
func (d Driver) Next(version uint) (nextVersion uint, err error) {
	k, _ := d.migrations.findByID(version)
	if k == -1 {
		return 0, errors.New("next version not found")
	}

	if len(d.migrations) < k+2 {
		return 0, os.ErrNotExist
	}

	return d.migrations[k+1].ID, nil
}
func (d Driver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	k, v := d.migrations.findByID(version)
	if k == -1 {
		return nil, "", errors.New("version not found")
	}

	return io.NopCloser(strings.NewReader(v.UpQuery)), "", nil
}

func (d Driver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	k, v := d.migrations.findByID(version)
	if k == -1 {
		return nil, "", errors.New("version not found")
	}

	return io.NopCloser(strings.NewReader(v.UpQuery)), "", nil
}
