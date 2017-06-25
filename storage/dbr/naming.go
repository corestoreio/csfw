// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

const quote string = "`"
const quoteRune = '`'

// Quoter at the quoter to use for quoting text; use Mysql quoting by default.
var Quoter = MysqlQuoter{
	replacer: strings.NewReplacer(quote, ""),
}

type alias struct {
	// Derived Tables (Subqueries in the FROM Clause). A derived table is a
	// subquery in a SELECT statement FROM clause. Derived tables can return a
	// scalar, column, row, or table. Ignored in any other case.
	DerivedTable *Select
	// Name can be any kind of SQL expression or a valid identifier. It gets
	// quoted when `IsExpression` is false.
	Name string
	// Alias must be a valid identifier allowed for alias usage.
	Alias string
	// IsExpression if true the field `Name` will be treated as an expression and
	// won't get quoted when generating the SQL.
	IsExpression bool
	// Sort applies only to GROUP BY and ORDER BY clauses. 'd'=descending,
	// 0=default or nothing; 'a'=ascending.
	Sort byte
}

const (
	sortDescending byte = 'd'
	sortAscending  byte = 'a'
)

// MakeNameAlias creates a new quoted name with an optional alias `a`, which can be
// empty.
func MakeNameAlias(name, a string) alias {
	return alias{
		Name:  name,
		Alias: a,
	}
}

// MakeExpressionAlias creates a new unquoted expression with an optional alias
// `a`, which can be empty.
func MakeExpressionAlias(expression, a string) alias {
	return alias{
		IsExpression: true,
		Name:         expression,
		Alias:        a,
	}
}

func (a alias) isEmpty() bool { return a.Name == "" && a.DerivedTable == nil }

// String returns the correct stringyfied statement.
func (a alias) String() string {
	if a.IsExpression {
		return Quoter.exprAlias(a.Name, a.Alias)
	}
	return a.QuoteAs()
}

// NameAlias always quuotes the name and the alias
func (a alias) QuoteAs() string {
	return Quoter.NameAlias(a.Name, a.Alias)
}

// appendArgs assembles the arguments and appends them to `args`
func (a alias) appendArgs(args Arguments) (_ Arguments, err error) {
	if a.DerivedTable != nil {
		args, err = a.DerivedTable.appendArgs(args)
	}
	return args, errors.WithStack(err)
}

// WriteQuoted writes the quoted table and its maybe alias into w.
func (a alias) WriteQuoted(w queryWriter) error {
	if a.DerivedTable != nil {
		w.WriteByte('(')
		if err := a.DerivedTable.toSQL(w); err != nil {
			return errors.WithStack(err)
		}
		w.WriteByte(')')
		w.WriteString(" AS ")
		Quoter.writeName(w, a.Alias)
		return nil
	}

	qf := Quoter.WriteNameAlias
	if a.IsExpression {
		qf = Quoter.WriteExpressionAlias
	}
	qf(w, a.Name, a.Alias)

	if a.Sort == sortAscending {
		w.WriteString(" ASC")
	}
	if a.Sort == sortDescending {
		w.WriteString(" DESC")
	}
	return nil
}

// TODO(CyS) if we need to distinguish between table name and the column or even need
// a sub select in the column list, then we can implement type aliases and replace
// all []string with type aliases. This costs some allocs but for modifying queries
// in dispatched events, it's getting easier.
type aliases []alias

// WriteQuoted writes all aliases comma separated and quoted into w.
func (as aliases) WriteQuoted(w queryWriter) error {
	for i, a := range as {
		if i > 0 {
			w.WriteString(", ")
		}
		if err := a.WriteQuoted(w); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (as aliases) appendArgs(args Arguments) (Arguments, error) {
	for _, a := range as {
		var err error
		args, err = a.appendArgs(args)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

// setSort applies to last n items the sort order `sort` in reverse iteration.
// Usuallay `lastNindexes` is len(object) because we decrement 1 from
// `lastNindexes`. This function panics when lastNindexes does not match the
// length of `aliases`.
func (as aliases) applySort(lastNindexes int, sort byte) aliases {
	to := len(as) - lastNindexes
	for i := len(as) - 1; i >= to; i-- {
		as[i].Sort = sort
	}
	return as
}

// AddColumns adds more columns to the aliases. Columns get quoted.
func (as aliases) AddColumns(columns ...string) aliases {
	return as.appendColumns(columns, false)
}

func (as aliases) appendColumns(columns []string, isExpression bool) aliases {
	if cap(as) == 0 {
		as = make(aliases, 0, len(columns)*2)
	}
	for _, c := range columns {
		as = append(as, alias{Name: c, IsExpression: isExpression})
	}
	return as
}

// columns must be balanced slice. i=column name, i+1=alias name
func (as aliases) appendColumnsAliases(columns []string, isExpression bool) aliases {
	if cap(as) == 0 {
		as = make(aliases, 0, len(columns)/2)
	}
	for i := 0; i < len(columns); i = i + 2 {
		as = append(as, alias{Name: columns[i], Alias: columns[i+1], IsExpression: isExpression})
	}
	return as
}

// MysqlQuoter implements Mysql-specific quoting
type MysqlQuoter struct {
	replacer *strings.Replacer
}

func (mq MysqlQuoter) unQuote(s string) string {
	if strings.IndexByte(s, quoteRune) == -1 {
		return s
	}
	return mq.replacer.Replace(s)
}

func (mq MysqlQuoter) writeName(w queryWriter, n string) {
	w.WriteByte(quoteRune)
	w.WriteString(mq.unQuote(n))
	w.WriteByte(quoteRune)
}

func (mq MysqlQuoter) writeQualifierName(w queryWriter, q, n string) {
	mq.writeName(w, q)
	w.WriteByte('.')
	mq.writeName(w, n)
}

// exprAlias appends to the provided `expression` the quote alias name, e.g.:
// 		exprAlias("(e.price*x.tax*t.weee)", "final_price") // (e.price*x.tax*t.weee) AS `final_price`
func (mq MysqlQuoter) exprAlias(expression, alias string) string {
	if alias == "" {
		return expression
	}
	return expression + " AS " + quote + mq.unQuote(alias) + quote
}

// Name quotes securely a name.
// 		Name("tableName") => `tableName`
// 		Name("table`Name") => `tableName`
// https://dev.mysql.com/doc/refman/5.7/en/identifier-qualifiers.html
func (mq MysqlQuoter) Name(n string) string {
	return quote + mq.unQuote(n) + quote
}

// QualifierName quotes securely a qualifier and its name.
// 		QualifierName("dbName", "tableName") => `dbName`.`tableName`
// 		QualifierName("db`Name", "`tableName`") => `dbName`.`tableName`
// https://dev.mysql.com/doc/refman/5.7/en/identifier-qualifiers.html
func (mq MysqlQuoter) QualifierName(q, n string) string {
	if q == "" {
		return mq.Name(n)
	}
	// return mq.Name(q) + "." + mq.Name(n) <-- too slow, too many allocs
	return quote + mq.unQuote(q) + quote + "." + quote + mq.unQuote(n) + quote
}

// WriteQualifierName same as function QualifierName but writes into w.
func (mq MysqlQuoter) WriteQualifierName(w queryWriter, q, n string) {
	if q == "" {
		mq.writeName(w, n)
		return
	}
	mq.writeName(w, q)
	w.WriteByte('.')
	mq.writeName(w, n)
}

// NameAlias quotes with back ticks and splits at a dot in the name. First
// argument table and/or column name (separated by a dot) and second argument
// can be an alias. Both parts will get quoted.
//		NameAlias("f", "g") 			// "`f` AS `g`"
//		NameAlias("e.entity_id", "ee") 	// `e`.`entity_id` AS `ee`
//		NameAlias("e.entity_id", "") 	// `e`.`entity_id`
func (mq MysqlQuoter) NameAlias(name, alias string) string {
	buf := bufferpool.Get()
	mq.WriteNameAlias(buf, name, alias)
	x := buf.String()
	bufferpool.Put(buf)
	return x
}

// WriteExpressionAlias writes an expression with an optinal quoted alias (which
// can be empty) into w.
func (mq MysqlQuoter) WriteExpressionAlias(w queryWriter, expression, alias string) {
	w.WriteString(expression)
	if alias != "" {
		w.WriteString(" AS ")
		mq.writeName(w, alias)
	}
}

// WriteNameAlias same as NameAlias but writes into w. It quotes always and each
// part. If a string contains quotes, they won't get stripped.
func (mq MysqlQuoter) WriteNameAlias(w queryWriter, name, alias string) {

	// checks if there are quotes at the beginning and at the end. no white spaces allowed.
	nameHasQuote := strings.HasPrefix(name, quote) && strings.HasSuffix(name, quote)
	nameHasDot := strings.IndexByte(name, '.') >= 0

	//fmt.Printf("lp %d expr %mq nameHasQuote %t nameHasDot %t | %#v\n", lp, expr, nameHasQuote, nameHasDot, expressionAlias)

	switch {
	case alias == "" && nameHasQuote:
		// already quoted
		w.WriteString(name)
		return
	case alias == "" && !nameHasQuote && !nameHasDot:
		// must be quoted
		mq.writeName(w, name)
		return
	case alias == "" && !nameHasQuote && nameHasDot:
		mq.splitDotAndQuote(w, name)
		return
	case name == "" && alias == "":
		// just an empty string
		return
	}

	mq.splitDotAndQuote(w, name)

	if alias != "" {
		w.WriteString(" AS ")
		mq.writeName(w, alias)
	}
}

func (mq MysqlQuoter) splitDotAndQuote(w queryWriter, part string) {
	dotIndex := strings.IndexByte(part, '.')
	if dotIndex > 0 { // dot at a beginning of a string at illegal
		mq.writeName(w, part[:dotIndex])
		w.WriteByte('.')
		if a := part[dotIndex+1:]; a == sqlStar {
			w.WriteByte('*')
		} else {
			mq.writeName(w, part[dotIndex+1:])
		}
		return
	}
	mq.writeName(w, part)
}

// TableColumnAlias prefixes all columns in the slice `cols` with a table
// name/alias and puts quotes around them. If a column name has already been
// prefixed by a name or an alias it will be ignored. This functions modifies
// the argument slice `cols`.
func (mq MysqlQuoter) TableColumnAlias(t string, cols ...string) []string {
	for i, c := range cols {
		switch {
		case strings.IndexByte(c, quoteRune) >= 0:
			cols[i] = c
		case strings.IndexByte(c, '.') > 0:
			cols[i] = mq.NameAlias(c, "")
		default:
			cols[i] = mq.QualifierName(t, c)
		}
	}
	return cols
}

// maxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const maxIdentifierLength = 64
const dummyQualifier = "X" // just a dummy value, can be optimized later

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters.
//
// Returns 0 if the identifier is valid.
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func isValidIdentifier(objectName string) int8 {
	if objectName == sqlStar {
		return 0
	}
	qualifier := dummyQualifier
	if i := strings.IndexByte(objectName, '.'); i >= 0 {
		qualifier = objectName[:i]
		objectName = objectName[i+1:]
	}

	validQualifier := isNameValid(qualifier)
	if validQualifier == 0 && objectName == sqlStar {
		return 0
	}
	if validQualifier > 0 {
		return validQualifier
	}
	return isNameValid(objectName)
}

// isNameValid returns 0 if the name is valid or an error number identifying
// where the name becomes invalid.
func isNameValid(name string) int8 {
	if name == dummyQualifier {
		return 0
	}

	ln := len(name)
	if ln > maxIdentifierLength || name == "" {
		return 1 //errors.NewNotValidf("[csdb] Incorrect identifier. Too long or empty: %q", name)
	}
	pos := 0
	for pos < ln {
		r, w := utf8.DecodeRuneInString(name[pos:])
		pos += w
		if !mapAlNum(r) {
			return 2 // errors.NewNotValidf("[csdb] Invalid character in name %q", name)
		}
	}
	return 0
}

func mapAlNum(r rune) bool {
	var ok bool
	switch {
	case '0' <= r && r <= '9':
		ok = true
	case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
		ok = true
	case r == '$', r == '_':
		ok = true
	}
	return ok
}
