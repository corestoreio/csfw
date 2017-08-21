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

package dbr_test

import (
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// Make sure that type productEntity implements interface.
var _ dbr.ArgumentsAppender = (*productEntity)(nil)

// productEntity represents just a demo record.
type productEntity struct {
	EntityID       int64 // Auto Increment
	AttributeSetID int64
	TypeID         string
	SKU            dbr.NullString
	HasOptions     bool
}

func (pe productEntity) AppendArgs(args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	l := len(columns)
	if l == 1 {
		// Most commonly used case
		return pe.appendBind(args, columns[0])
	}
	if l == 0 {
		// This case gets executed when an INSERT statement doesn't contain any
		// columns.
		return args.Int64(pe.EntityID).Int64(pe.AttributeSetID).Str(pe.TypeID).NullString(pe.SKU).Bool(pe.HasOptions), nil
	}
	// This case gets executed when an INSERT statement requests specific columns.
	for _, col := range columns {
		var err error
		if args, err = pe.appendBind(args, col); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

func (pe productEntity) appendBind(args dbr.Arguments, column string) (_ dbr.Arguments, err error) {
	switch column {
	case "attribute_set_id":
		args = args.Int64(pe.AttributeSetID)
	case "type_id":
		args = args.Str(pe.TypeID)
	case "sku":
		args = args.NullString(pe.SKU)
	case "has_options":
		args = args.Bool(pe.HasOptions)
	default:
		return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", column)
	}
	return args, nil
}

// ExampleInsert_Bind inserts new data into table
// `catalog_product_entity`. First statement by specifying the exact column
// names. In the second example all columns values are getting inserted and you
// specify the number of place holders per record.
func ExampleInsert_Bind() {

	objs := []productEntity{
		{1, 5, "simple", dbr.MakeNullString("SOA9"), false},
		{2, 5, "virtual", dbr.NullString{}, true},
	}

	i := dbr.NewInsert("catalog_product_entity").AddColumns("attribute_set_id", "type_id", "sku", "has_options").
		Bind(objs[0]).Bind(objs[1])
	writeToSQLAndInterpolate(i)

	fmt.Print("\n\n")
	i = dbr.NewInsert("catalog_product_entity").SetRecordValueCount(5).Bind(objs[0]).Bind(objs[1])
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_entity`
	//(`attribute_set_id`,`type_id`,`sku`,`has_options`) VALUES (?,?,?,?),(?,?,?,?)
	//Arguments: [5 simple SOA9 false 5 virtual <nil> true]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_entity`
	//(`attribute_set_id`,`type_id`,`sku`,`has_options`) VALUES
	//(5,'simple','SOA9',0),(5,'virtual',NULL,1)
	//
	//Prepared Statement:
	//INSERT INTO `catalog_product_entity` VALUES (?,?,?,?,?),(?,?,?,?,?)
	//Arguments: [1 5 simple SOA9 false 2 5 virtual <nil> true]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_entity` VALUES
	//(1,5,'simple','SOA9',0),(2,5,'virtual',NULL,1)
}