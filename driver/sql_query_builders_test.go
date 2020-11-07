package driver

import "testing"

func TestSQLSelectQueryGenerator_Compile(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	q, a := builder.Compile()
	if len(a) != 0 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_SelectRaw(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	q, a := builder.Compile()
	if len(a) != 0 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_From(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("test")
	q, a := builder.Compile()
	if len(a) != 1 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $1" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_Where(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("test")
	builder.Where("a", WhereEqual, "b")
	q, a := builder.Compile()
	if len(a) != 3 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $1 WHERE $2 = $3" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_Limit(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.Limit(10)
	q, a := builder.Compile()
	if len(a) != 0 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 LIMIT 10" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_FullJoin(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("testing")
	builder.FullJoin("test", "test", "test", "test")
	q, a := builder.Compile()
	if len(a) != 5 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $5 FULL OUTER JOIN $3 ON $1.$2 = $3.$4" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_InnerJoin(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("testing")
	builder.InnerJoin("test", "test", "test", "test")
	q, a := builder.Compile()
	if len(a) != 5 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $5 INNER JOIN $3 ON $1.$2 = $3.$4" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_LeftJoin(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("testing")
	builder.LeftJoin("test", "test", "test", "test")
	q, a := builder.Compile()
	if len(a) != 5 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $5 LEFT JOIN $3 ON $1.$2 = $3.$4" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_RightJoin(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("testing")
	builder.RightJoin("test", "test", "test", "test")
	q, a := builder.Compile()
	if len(a) != 5 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $5 RIGHT JOIN $3 ON $1.$2 = $3.$4" {
		t.Fatal("weird query:", q)
	}
}

func TestSQLSelectQueryGenerator_SelectColumn(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("testing")
	builder.RightJoin("test", "test", "test", "test")
	builder.Where("testing", WhereNotEqual, "a")
	q, a := builder.Compile()
	if len(a) != 7 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM $5 RIGHT JOIN $3 ON $1.$2 = $3.$4 WHERE $6 <> $7" {
		t.Fatal("weird query:", q)
	}
}

func TestNewSQLSelectQueryGenerator_MultiJoin(t *testing.T) {
	builder := NewSQLSelectQueryGenerator()
	builder.SelectRaw("1")
	builder.From("testing")
	builder.LeftJoin("test", "test", "test", "test")
	builder.FullJoin("test", "test", "test", "test")
	builder.RightJoin("test", "test", "test", "test")
	builder.Limit(10)
	q, a := builder.Compile()
	if len(a) != 13 {
		t.Fatal("weird array:", a)
	}
	if q != "SELECT 1 FROM (($5 LEFT JOIN $3 ON $1.$2 = $3.$4) FULL OUTER JOIN $8 ON $6.$7 = $8.$9) RIGHT JOIN $12 ON $10.$11 = $12.$13 LIMIT 10" {
		t.Fatal("weird query:", q)
	}
}
