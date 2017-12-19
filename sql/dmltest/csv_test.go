// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dmltest_test

import (
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
)

func TestLoadCSVWithFile(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithFile("testdata", "core_config_data1.csv"),
		dmltest.WithTestMode(),
	)
	assert.NoError(t, err)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 20)

	want := "[][]driver.Value{[]driver.Value{text.Chars(`1`), text.Chars(`default`), text.Chars(`0`), text.Chars(`cms/wysiwyg/enabled`), text.Chars(`disabled`)}, []driver.Value{text.Chars(`2`), text.Chars(`default`), text.Chars(`0`), text.Chars(`general/region/display_all`), text.Chars(`1`)}, []driver.Value{text.Chars(`3`), text.Chars(`default`), text.Chars(`0`), text.Chars(`general/region/state_required`), text.Chars(`AT,CA,CH,DE,EE,ES,FI,FR,LT,LV,RO,US`)}, []driver.Value{text.Chars(`3`), text.Chars(`stores`), text.Chars(`2`), text.Chars(`general/region/state_required`), text.Chars(`AT`)}, []driver.Value{text.Chars(`5`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/url/redirect_to_base`), text.Chars(`1`)}, []driver.Value{text.Chars(`7`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/unsecure/base_url`), text.Chars(`http://magento-1-8.local/`)}, []driver.Value{text.Chars(`7`), text.Chars(`websites`), text.Chars(`1`), text.Chars(`web/unsecure/base_url`), text.Chars(`http://magento-1-8a.dev/`)}, []driver.Value{text.Chars(`8`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/unsecure/base_link_url`), text.Chars(`{{unsecure_base_url}}`)}, []driver.Value{text.Chars(`9`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/unsecure/base_skin_url`), text.Chars(`{{unsecure_base_url}}skin/`)}, []driver.Value{text.Chars(`10`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/unsecure/base_media_url`), text.Chars(`http://localhost:4711/media/`)}, []driver.Value{text.Chars(`11`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/unsecure/base_js_url`), text.Chars(`{{unsecure_base_url}}js/`)}, []driver.Value{text.Chars(`12`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/base_url`), text.Chars(`http://magento-1-8.local/`)}, []driver.Value{text.Chars(`13`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/base_link_url`), text.Chars(`{{secure_base_url}}`)}, []driver.Value{text.Chars(`14`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/base_skin_url`), text.Chars(`{{secure_base_url}}skin/`)}, []driver.Value{text.Chars(`15`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/base_media_url`), text.Chars(`http://localhost:4711/media/`)}, []driver.Value{text.Chars(`16`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/base_js_url`), text.Chars(`{{secure_base_url}}js/`)}, []driver.Value{text.Chars(`17`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/use_in_frontend`), text.Chars(`0`)}, []driver.Value{text.Chars(`18`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/use_in_adminhtml`), text.Chars(`0`)}, []driver.Value{text.Chars(`19`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/secure/offloader_header`), text.Chars(`SSL_OFFLOADED`)}, []driver.Value{text.Chars(`20`), text.Chars(`default`), text.Chars(`0`), text.Chars(`web/default/front`), driver.Value(nil)}}"
	assert.Exactly(t, want, fmt.Sprintf("%#v", rows))
}

func TestLoadCSVWithReaderConfig(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithTestMode(),
		dmltest.WithFile("testdata", "core_config_data3.csv"),
		dmltest.WithReaderConfig(dmltest.CSVConfig{Comma: '|'}),
	)
	assert.NoError(t, err)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 5)

	want := "[][]driver.Value{[]driver.Value{text.Chars(`1`), text.Chars(`default`), text.Chars(`0`), text.Chars(`cms/wysiwyg/enabled`), text.Chars(`disabled`)}, []driver.Value{text.Chars(`2`), text.Chars(`default`), text.Chars(`0`), text.Chars(`general/region/display_all`), text.Chars(`1`)}, []driver.Value{text.Chars(`3`), text.Chars(`default`), text.Chars(`0`), text.Chars(`general/region/state_required`), text.Chars(`AT,CA,CH,DE,EE,ES,FI,FR,LT,LV,RO,US`)}, []driver.Value{text.Chars(`3`), text.Chars(`stores`), text.Chars(`2`), text.Chars(`general/region/state_required`), text.Chars(`AT`)}, []driver.Value{text.Chars(`5`), text.Chars(`default`), text.Chars(`0`), driver.Value(nil), text.Chars(`1`)}}"
	assert.Exactly(t, want, fmt.Sprintf("%#v", rows))
}

func TestLoadCSVFileError(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithTestMode(),
		dmltest.WithFile("testdata", "core_config_dataXX.csv"),
	)
	assert.Nil(t, cols)
	assert.Nil(t, rows)
	assert.Contains(t, err.Error(), "core_config_dataXX.csv: no such file or directory")
}

func TestLoadCSVReadError(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithFile("testdata", "core_config_data2.csv"),
		dmltest.WithTestMode(),
	)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 5)
	assert.EqualError(t, err, "[cstesting] csvReader.Read: record on line 7; parse error on line 8, column 0: extraneous or missing \" in quoted-field")
}

func TestMockRowsError(t *testing.T) {
	t.Parallel()
	r, err := dmltest.MockRows(dmltest.WithFile("non", "existent.csv"))
	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "non/existent.csv: no such file or directory")
}

func TestMockRowsLoaded(t *testing.T) {
	t.Parallel()
	rows, err := dmltest.MockRows(
		dmltest.WithReaderConfig(dmltest.CSVConfig{Comma: '|'}),
		dmltest.WithFile("testdata", "core_config_data3.csv"),
		dmltest.WithTestMode(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, rows)

	// Sorry for this test, but they removed the .Columns() function
	assert.Contains(t, fmt.Sprintf("%#v", rows), `sqlmock.Rows{cols:[]string{"config_id", "scope", "scope_id", "path", "value"}`)
}

func TestMustMockRows(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(error).Error(), "non/existent.csv: no such file or directory")
		} else {
			t.Fatal("Expecting a panic")
		}
	}()

	r := dmltest.MustMockRows(dmltest.WithFile("non", "existent.csv"))
	assert.Nil(t, r)
}