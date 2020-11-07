package driver

// WhereOperator is used to define the operator in a where statement.
type WhereOperator string

// Defines the where operators.
const (
	WhereEqual              WhereOperator = "="
	WhereNotEqual           WhereOperator = "<>"
	WhereLessThanOrEqual    WhereOperator = "<="
	WhereGreaterThanOrEqual WhereOperator = ">="
	WhereGreaterThan        WhereOperator = ">"
	WhereLessThan           WhereOperator = "<"
)

// SelectQueryBuilder is the interface used for building select SQL queries.
type SelectQueryBuilder interface {
	// InnerJoin is used to perform a inner join on the table.
	InnerJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string)

	// LeftJoin is used to perform a left join on the table.
	LeftJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string)

	// RightJoin is used to perform a right join on the table.
	RightJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string)

	// FullJoin is used to perform a full join on the table.
	FullJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string)

	// SelectColumn used to add a column to select.
	SelectColumn(Table, Column string)

	// SelectRaw is used to add a raw selection.
	SelectRaw(Query string)

	// From is used to set the "FROM" statement.
	From(Table string)

	// Limit is used to specify the query limit.
	Limit(n int)

	// Where is used to handle selecting from the table.
	Where(Column string, Operator WhereOperator, Value interface{})

	// Compile is used to compile the query. The first argument is passed to Translate (whose result is then passed to Run),
	// the second is passed to Run with the Translate result.
	Compile() (interface{}, []interface{})
}

// QueryResult is used to handle a query result.
type QueryResult interface {
	// Next prepares the next result row for reading with the Scan method. Like the database/sql method, this returns true on success.
	Next() bool

	// Close is used to close the query.
	Close() error

	// Scan is used to scan a database row. The pointers will be in the same order as the items in the select clause.
	Scan(ptrs ...interface{}) error
}

// DBDriver is used to define a driver for a specific database.
type DBDriver interface {
	// NewSelectQueryBuilder is used to create a new select query builder.
	NewSelectQueryBuilder() SelectQueryBuilder

	// Translate is used to turn the SQL query string to a interface. This is called if the statement is not cached.
	Translate(stmt interface{}) (interface{}, error)

	// Exec is used to execute the SQL query. Note the statement is compiled earlier due to caching, hence why it is passed back through here.
	Exec(translatedStmt interface{}, args []interface{}) (rows int64, err error)

	// Query is used to do a database query. Note the statement is compiled earlier due to caching, hence why it is passed back through here.
	Query(translatedStmt interface{}, args []interface{}) (result QueryResult, err error)
}

type done struct {}
func (done) Next() bool { return false }
func (done) Close() error { return nil }
// TODO: custom error
func (done) Scan(...interface{}) error { return nil }

// DoneQueryResult is used to generate a query result which is done.
func DoneQueryResult() QueryResult {
	return &done{}
}
