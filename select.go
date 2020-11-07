package fastsql

import (
	"fastsql/driver"
	"reflect"
)

// Used to define a select SQL filter.
type selectFilter struct {
	column string
	op driver.WhereOperator
	res interface{}
	limit int
}

// Limit is used to create a select filter which limits the number of results.
// Note this is ignored if the destination is not an array.
func Limit(n uint) *selectFilter {
	return &selectFilter{limit: int(n)}
}

// Equal is used to create a select filter based on a column being equal to the interface specified.
func Equal(ColumnName string, EqualTo interface{}) *selectFilter {
	return &selectFilter{op: driver.WhereEqual, column: ColumnName, res: EqualTo, limit: -1}
}

// NotEqual is used to create a select filter based on a column being not equal to the interface specified.
func NotEqual(ColumnName string, NotEqualTo interface{}) *selectFilter {
	return &selectFilter{op: driver.WhereNotEqual, column: ColumnName, res: NotEqualTo, limit: -1}
}

// LessThanOrEqual is used to create a select filter based on a column being less than or equal to the interface specified.
func LessThanOrEqual(ColumnName string, LeTo interface{}) *selectFilter {
	return &selectFilter{op: driver.WhereLessThanOrEqual, column: ColumnName, res: LeTo, limit: -1}
}

// GreaterThanOrEqual is used to create a select filter based on a column being greater than or equal to the interface specified.
func GreaterThanOrEqual(ColumnName string, GeTo interface{}) *selectFilter {
	return &selectFilter{op: driver.WhereGreaterThanOrEqual, column: ColumnName, res: GeTo, limit: -1}
}

// GreaterThan is used to create a select filter based on a column being greater than the interface specified.
func GreaterThan(ColumnName string, Gt interface{}) *selectFilter {
	return &selectFilter{op: driver.WhereGreaterThan, column: ColumnName, res: Gt, limit: -1}
}

// GreaterThan is used to create a select filter based on a column being less than the interface specified.
func LessThan(ColumnName string, Lt interface{}) *selectFilter {
	return &selectFilter{op: driver.WhereLessThan, column: ColumnName, res: Lt, limit: -1}
}

// Select is used to get an item/array of items from the database which match the criteria specified.
func (q *QueryManager) Select(dest interface{}, filters ...*selectFilter) error {
	// Get the select query builder.
	builder := q.d.NewSelectQueryBuilder()

	// Run reflect on the destination.
	r := reflect.ValueOf(dest)

	// Get the kind of the reflect value.
	kind := r.Kind()

	// Checks if this is a slice.
	isArray := kind == reflect.Slice

	// Build a key array.
	//keys := make([]string, 0)
	itemType := r.Type()
	if isArray {
		// This is an array. Get the type.
		itemType = r.Type().Elem()
	}
	structKind := itemType.Kind()
	if structKind != reflect.Ptr  {

	}

	//builder.From(tableName)

	// Handle setting the filters.
	for _, v := range filters {
		if v.limit == -1 {
			// This is a where filter.
			builder.Where(v.column, v.op, v.res)
			continue
		}
		if isArray {
			// This is an array. We can limit this.
			builder.Limit(v.limit)
		}
	}

	// Compile the query and get the statement.
	stmt, args := builder.Compile()
	translatedStmt, err := q.handleTranslationCache(stmt)
	if err != nil {
		return err
	}

	// Run the query.
	res, err := q.d.Query(translatedStmt, args)
	if err != nil {
		return err
	}
	// TODO
	return res.Close()
}
