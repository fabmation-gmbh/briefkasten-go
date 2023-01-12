package models

import (
	"context"
	"crypto/tls"
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/fabmation-gmbh/briefkasten-go/internal/config"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
)

var db *bun.DB

// GetDB returns the internal db variable.
// NOTE: Please use this variable only if you know what you are doing AND if it necessary!
func GetDB() *bun.DB {
	return db
}

// Connect tries to open the SQLite DB and if it does not exist, it will create a new one.
func Connect() {
	dsn := os.Getenv("DB_CONNECTION")
	if dsn == "" {
		dsn = config.C.DB.URI
	}

	opts := []pgdriver.Option{
		pgdriver.WithDSN(dsn),
	}

	if strings.HasSuffix(dsn, "require") {
		opts = append(opts, pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		opts...,
	))

	db = bun.NewDB(sqlDB, pgdialect.New())

	if config.C.Debug.EnableSQLDebug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	if config.C.Debug.EnableTracing {
		db.AddQueryHook(bunotel.NewQueryHook(
			bunotel.WithDBName("nextcloud-vfs"),
			bunotel.WithFormattedQueries(true),
		))
	}
}

// newNullInt64 returns a valid and initialized sql.NullInt64,
func newNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{
		Valid: true,
		Int64: n,
	}
}

// StartTx starts a transaction with the given options.
func StartTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return bun.Tx{}, errors.Wrap(err, "unable to start transaction")
	}

	return tx, nil
}

// getDB returns the base DB connection or the given
// tx.
func getDB(dbs ...bun.IDB) bun.IDB {
	if len(dbs) == 0 {
		return db
	}

	return dbs[0]
}

func getBoolEnv(name string) bool {
	str := os.Getenv(name)

	ret, err := strconv.ParseBool(str)

	return err != nil && ret
}

// IsNoRows returns true if the given error represents
// the error, returned if no rows in the result set where found.
func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// IsUniqueViolationErr return true if the given error represents
// a unique constraint violation.
//
// See https://www.postgresql.org/docs/current/errcodes-appendix.html
func IsUniqueViolationErr(err error) bool {
	e, ok := getPsqlError(err)
	if !ok {
		return false
	}

	return e.Field('C') == "23505"
}

func getPsqlError(err error) (pgdriver.Error, bool) {
	pg, ok := err.(pgdriver.Error)
	if !ok {
		return pgdriver.Error{}, false
	}

	return pg, true
}
