package driver

import "strconv"

// Defines a join.
type join struct {
	// 0 = inner, 1 = left, 2 = right, 3 = full
	joinType uint8

	// Defines all the data required for the join.
	JoiningTable, FromTable, JoiningTableKey, FromTableKey string
}

// Defines a select query.
type selectQuery struct {
	Column   string
	Operator WhereOperator
	Value    interface{}
}

// Defines a select from object.
type selectFrom struct {
	Table, Column string
}

// SQLSelectQueryGenerator is used to generate SQL queries.
type SQLSelectQueryGenerator struct {
	joins []*join
	selection []interface{} // can be string for raw select or a selectFrom object.
	from  string
	where []*selectQuery
	limit int // -1 here means unlimited.
}

// NewSQLSelectQueryGenerator is used to create a new query generator.
func NewSQLSelectQueryGenerator() *SQLSelectQueryGenerator {
	return &SQLSelectQueryGenerator{
		joins: make([]*join, 0),
		where: make([]*selectQuery, 0),
		selection: make([]interface{}, 0),
		limit: -1,
	}
}

func (s *SQLSelectQueryGenerator) InnerJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string) {
	s.joins = append(s.joins, &join{
		joinType:        0,
		JoiningTable:    JoiningTable,
		FromTable:       FromTable,
		JoiningTableKey: JoiningTableKey,
		FromTableKey:    FromTableKey,
	})
}

func (s *SQLSelectQueryGenerator) LeftJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string) {
	s.joins = append(s.joins, &join{
		joinType:        1,
		JoiningTable:    JoiningTable,
		FromTable:       FromTable,
		JoiningTableKey: JoiningTableKey,
		FromTableKey:    FromTableKey,
	})
}

func (s *SQLSelectQueryGenerator) RightJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string) {
	s.joins = append(s.joins, &join{
		joinType:        2,
		JoiningTable:    JoiningTable,
		FromTable:       FromTable,
		JoiningTableKey: JoiningTableKey,
		FromTableKey:    FromTableKey,
	})
}

func (s *SQLSelectQueryGenerator) FullJoin(JoiningTable, FromTable, JoiningTableKey, FromTableKey string) {
	s.joins = append(s.joins, &join{
		joinType:        3,
		JoiningTable:    JoiningTable,
		FromTable:       FromTable,
		JoiningTableKey: JoiningTableKey,
		FromTableKey:    FromTableKey,
	})
}

func (s *SQLSelectQueryGenerator) SelectColumn(Table, Column string) {
	s.selection = append(s.selection, &selectFrom{
		Table:  Table,
		Column: Column,
	})
}

func (s *SQLSelectQueryGenerator) SelectRaw(Query string) {
	s.selection = append(s.selection, Query)
}

func (s *SQLSelectQueryGenerator) From(Table string) {
	s.from = Table
}

func (s *SQLSelectQueryGenerator) Limit(n int) {
	s.limit = n
}

func (s *SQLSelectQueryGenerator) Where(Column string, Operator WhereOperator, Value interface{}) {
	s.where = append(s.where, &selectQuery{
		Column:   Column,
		Operator: Operator,
		Value:    Value,
	})
}

func (s *SQLSelectQueryGenerator) Compile() (interface{}, []interface{}) {
	query := "SELECT "
	args := make([]interface{}, 0, (len(s.where) * 2) + len(s.where))
	argNum := 1
	if len(s.selection) != 0 {
		// Generate the "SELECT ..." part of the query.
		for _, item := range s.selection {
			raw, ok := item.(string)
			if ok {
				query += raw + ", "
				continue
			}
			sf := item.(*selectFrom)
			arg1 := strconv.Itoa(argNum)
			argNum++
			arg2 := strconv.Itoa(argNum)
			argNum++
			query += "$" + arg1 + "." + "$" + arg2 + ", "
			args = append(args, sf.Table, sf.Column)
		}
		query = query[:len(query)-2]
	} else {
		// Trim white space.
		query = query[:len(query)-1]
	}
	if s.from != "" {
		// Handle the "FROM ..." part of the query.
		query += " FROM "
		if len(s.joins) == 0 {
			// This is a simple one table query.
			query += "$" + strconv.Itoa(argNum)
			argNum++
			args = append(args, s.from)
		} else {
			// Handle joins with this query.
			joinText := ""
			addJoin := func(fragment string) {
				if joinText == "" {
					// Handle a join start.
					args = append(args, s.from)
					joinText = "$" + strconv.Itoa(argNum) + " " + fragment
					argNum++
				} else {
					// Handle a join continuation.
					joinText = "(" + joinText + ") " + fragment
				}
			}
			for _, v := range s.joins {
				fragment := ""
				switch v.joinType {
				case 0:
					fragment = "INNER JOIN "
				case 1:
					fragment = "LEFT JOIN "
				case 2:
					fragment = "RIGHT JOIN "
				case 3:
					fragment = "FULL OUTER JOIN "
				}
				fragment += "$" + strconv.Itoa(argNum + 2) + " ON $" + strconv.Itoa(argNum) + ".$" + strconv.Itoa(argNum + 1) + " = $" + strconv.Itoa(argNum + 2) + ".$" + strconv.Itoa(argNum + 3)
				args = append(args, v.FromTable, v.FromTableKey, v.JoiningTable, v.JoiningTableKey)
				argNum += 4
				addJoin(fragment)
			}
			query += joinText
		}
	}
	if len(s.where) != 0 {
		// Set the "WHERE ..." clause.
		query += " WHERE "
		for i, v := range s.where {
			if i != 0 {
				query += " AND "
			}
			query += "$" + strconv.Itoa(argNum) + " " + (string)(v.Operator) + " $" + strconv.Itoa(argNum + 1)
			args = append(args, v.Column, v.Value)
			argNum += 3
		}
	}
	if s.limit > 0 {
		// Set the query limit.
		query += " LIMIT " + strconv.Itoa(s.limit)
	}
	return query, args
}
