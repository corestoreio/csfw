// Code generated by codegen. DO NOT EDIT.
// Generated by sql/dmlgen. DO NOT EDIT.
package dmltestgenerated4

import (
	"context"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
	"sort"
	"testing"
	"time"
)

func TestNewDBManagerDB_48a8450c0b62e880b2d40acd0bbbd0dc(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)
	defer dmltest.SQLDumpLoad(t, "../testdata/test_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	}).Deferred()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	tbls, err := NewDBManager(ctx, &DBMOption{TableOptions: []ddl.TableOption{ddl.WithConnPool(db)}})
	assert.NoError(t, err)
	tblNames := tbls.Tables.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{"core_configuration", "sales_order_status_state", "view_customer_auto_increment"}, tblNames)
	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de", FloatMaxDecimals: 6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
	)
	t.Run("CoreConfiguration_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCoreConfiguration)
		selOneRow := tbl.Select("*").Where(
			dml.Column("config_id").Equal().PlaceHolder(),
		)
		selTenRows := tbl.Select("*").Where(
			dml.Column("config_id").LessOrEqual().Int(10),
		)
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx, tbl.Insert().BuildValues())
		for i := 0; i < 9; i++ {
			entIn := new(CoreConfiguration)
			assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CoreConfiguration_Entity")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entIn)))
			entINSERTStmtA.Reset()
			entOut := new(CoreConfiguration)
			rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.ConfigID, entOut.ConfigID, "IDX%d: ConfigID should match", lID)
			assert.ExactlyLength(t, 8, &entIn.Scope, &entOut.Scope, "IDX%d: Scope should match", lID)
			assert.Exactly(t, entIn.ScopeID, entOut.ScopeID, "IDX%d: ScopeID should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Path, &entOut.Path, "IDX%d: Path should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.Value, &entOut.Value, "IDX%d: Value should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)
		entCol := NewCoreConfigurations()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())
		lID := dmltest.CheckLastInsertID(t, "Error:  CoreConfigurations ")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))
		t.Logf("Last insert ID into: %d", lID)
	})
	t.Run("SalesOrderStatusState_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameSalesOrderStatusState)
		selOneRow := tbl.Select("*").Where()
		selTenRows := tbl.Select("*").Where()
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		// this table/view does not support auto_increment
		entCol := NewSalesOrderStatusStates()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
	})
	t.Run("ViewCustomerAutoIncrement_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameViewCustomerAutoIncrement)
		selOneRow := tbl.Select("*").Where()
		selTenRows := tbl.Select("*").Where()
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		// this table/view does not support auto_increment
		entCol := NewViewCustomerAutoIncrements()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
	})
	// Uncomment the next line for debugging to see all the queries.
	// t.Logf("queries: %#v", tbls.ConnPool.CachedQueries())
}
