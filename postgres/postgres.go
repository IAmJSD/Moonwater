package postgres

import (
	"database/sql"
	"fastsql/driver"
	_ "github.com/lib/pq"
)

type pgDriver struct {
	db *sql.DB
}

func (pgDriver) NewSelectQueryBuilder() driver.SelectQueryBuilder {
	return driver.NewSQLSelectQueryGenerator()
}

func (p pgDriver) Translate(stmt interface{}) (interface{}, error) {
	return p.db.Prepare(stmt.(string))
}

func (p pgDriver) Exec(translatedStmt interface{}, args []interface{}) (affectedRows int64, err error) {
	res, err := translatedStmt.(*sql.Stmt).Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (p pgDriver) Query(translatedStmt interface{}, args []interface{}) (res driver.QueryResult, err error) {
	return translatedStmt.(*sql.Stmt).Query(args...)
}

// New is used to create a new instance of the postgres driver.
func New(ConnectionString string) (driver.DBDriver, error) {
	db, err := sql.Open("postgres", ConnectionString)
	if err != nil {
		return nil, err
	}
	return &pgDriver{db: db}, nil
}
