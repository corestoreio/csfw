// Code generated by codegen. DO NOT EDIT.
// Generated by sql/dmlgen. DO NOT EDIT.
package dmltestgenerated2

import (
	"fmt"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"io"
	"time"
)

// CoreConfiguration represents a single row for DB table core_configuration.
// Auto generated.
// Table comment: Config Data
type CoreConfiguration struct {
	ConfigID  uint32      // config_id int(10) unsigned NOT NULL PRI  auto_increment "Id"
	Scope     string      // scope varchar(8) NOT NULL MUL DEFAULT ''default''  "Scope"
	ScopeID   int32       // scope_id int(11) NOT NULL  DEFAULT '0'  "Scope Id"
	Expires   null.Time   // expires datetime NULL  DEFAULT 'NULL'  "Value expiration time"
	Path      string      // path varchar(255) NOT NULL    "Path"
	Value     null.String // value text NULL  DEFAULT 'NULL'  "Value"
	VersionTs time.Time   // version_ts timestamp(6) NOT NULL   STORED GENERATED "Timestamp Start Versioning"
	VersionTe time.Time   // version_te timestamp(6) NOT NULL PRI  STORED GENERATED "Timestamp End Versioning"
}

// Copy copies the struct and returns a new pointer. TODO use deepcopy tool to
// generate code afterwards
func (e *CoreConfiguration) Copy() *CoreConfiguration {
	if e == nil {
		return &CoreConfiguration{}
	}
	e2 := *e // for now a shallow copy
	return &e2
}

// Empty empties all the fields of the current object. Also known as Reset.
func (e *CoreConfiguration) Empty() *CoreConfiguration { *e = CoreConfiguration{}; return e }

// This variable can be set in another file to provide a custom validator.
var validateCoreConfiguration func(*CoreConfiguration) error

// Validate runs internal consistency tests.
func (e *CoreConfiguration) Validate() error {
	if e == nil {
		return errors.NotValid.Newf("Type %T cannot be nil", e)
	}
	if validateCoreConfiguration != nil {
		return validateCoreConfiguration(e)
	}
	return nil
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (e *CoreConfiguration) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.
	n2, err := fmt.Fprint(w,
		"config_id:", e.ConfigID, "\n",
		"scope:", e.Scope, "\n",
		"scope_id:", e.ScopeID, "\n",
		"expires:", e.Expires, "\n",
		"path:", e.Path, "\n",
		"value:", e.Value, "\n",
		"version_ts:", e.VersionTs, "\n",
		"version_te:", e.VersionTe, "\n",
	)
	return int64(n2), err
}

// CoreConfigurations represents a collection type for DB table
// core_configuration
// Not thread safe. Auto generated.
type CoreConfigurations struct {
	Data []*CoreConfiguration `json:"data,omitempty"`
}

// NewCoreConfigurations  creates a new initialized collection. Auto generated.
func NewCoreConfigurations() *CoreConfigurations {
	return &CoreConfigurations{
		Data: make([]*CoreConfiguration, 0, 5),
	}
}

// Append will add a new item at the end of * CoreConfigurations . Auto generated
// via dmlgen.
func (cc *CoreConfigurations) Append(n ...*CoreConfiguration) *CoreConfigurations {
	cc.Data = append(cc.Data, n...)
	return cc
}

// Cut will remove items i through j-1. Auto generated via dmlgen.
func (cc *CoreConfigurations) Cut(i, j int) *CoreConfigurations {
	z := cc.Data // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this avoids the memory leak
	}
	z = z[:len(z)-j+i]
	cc.Data = z
	return cc
}

// Delete will remove an item from the slice. Auto generated via dmlgen.
func (cc *CoreConfigurations) Delete(i int) *CoreConfigurations {
	z := cc.Data // copy the slice header
	end := len(z) - 1
	cc.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	cc.Data = z
	return cc
}

// Each will run function f on all items in []* CoreConfiguration . Auto
// generated via dmlgen.
func (cc *CoreConfigurations) Each(f func(*CoreConfiguration)) *CoreConfigurations {
	if cc == nil {
		return nil
	}
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Filter filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (cc *CoreConfigurations) Filter(f func(*CoreConfiguration) bool) *CoreConfigurations {
	if cc == nil {
		return nil
	}
	b, i := cc.Data[:0], 0
	for _, e := range cc.Data {
		if f(e) {
			b = append(b, e)
		}
		i++
	}
	for i := len(b); i < len(cc.Data); i++ {
		cc.Data[i] = nil // this should avoid the memory leak
	}
	cc.Data = b
	return cc
}

// Insert will place a new item at position i. Auto generated via dmlgen.
func (cc *CoreConfigurations) Insert(n *CoreConfiguration, i int) *CoreConfigurations {
	z := cc.Data // copy the slice header
	z = append(z, &CoreConfiguration{})
	copy(z[i+1:], z[i:])
	z[i] = n
	cc.Data = z
	return cc
}

// Swap will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *CoreConfigurations) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }

// Len will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *CoreConfigurations) Len() int {
	if cc == nil {
		return 0
	}
	return len(cc.Data)
}

// Validate runs internal consistency tests on all items.
func (cc *CoreConfigurations) Validate() (err error) {
	if len(cc.Data) == 0 {
		return nil
	}
	for i, ld := 0, len(cc.Data); i < ld && err == nil; i++ {
		err = cc.Data[i].Validate()
	}
	return
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (cc *CoreConfigurations) WriteTo(w io.Writer) (n int64, err error) {
	for i, d := range cc.Data {
		n2, err := d.WriteTo(w)
		if err != nil {
			return 0, errors.Wrapf(err, "[dmltestgenerated2] WriteTo failed at index %d", i)
		}
		n += n2
	}
	return n, nil
}

// SalesOrderStatusState represents a single row for DB table
// sales_order_status_state. Auto generated.
// Table comment: Sales Order Status Table
type SalesOrderStatusState struct {
	Status         string // status varchar(32) NOT NULL PRI   "Status"
	State          string // state varchar(32) NOT NULL PRI   "Label"
	IsDefault      bool   // is_default smallint(5) unsigned NOT NULL  DEFAULT '0'  "Is Default"
	VisibleOnFront uint16 // visible_on_front smallint(5) unsigned NOT NULL  DEFAULT '0'  "Visible on front"
}

// Copy copies the struct and returns a new pointer. TODO use deepcopy tool to
// generate code afterwards
func (e *SalesOrderStatusState) Copy() *SalesOrderStatusState {
	if e == nil {
		return &SalesOrderStatusState{}
	}
	e2 := *e // for now a shallow copy
	return &e2
}

// Empty empties all the fields of the current object. Also known as Reset.
func (e *SalesOrderStatusState) Empty() *SalesOrderStatusState {
	*e = SalesOrderStatusState{}
	return e
}

// This variable can be set in another file to provide a custom validator.
var validateSalesOrderStatusState func(*SalesOrderStatusState) error

// Validate runs internal consistency tests.
func (e *SalesOrderStatusState) Validate() error {
	if e == nil {
		return errors.NotValid.Newf("Type %T cannot be nil", e)
	}
	if validateSalesOrderStatusState != nil {
		return validateSalesOrderStatusState(e)
	}
	return nil
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (e *SalesOrderStatusState) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.
	n2, err := fmt.Fprint(w,
		"status:", e.Status, "\n",
		"state:", e.State, "\n",
		"is_default:", e.IsDefault, "\n",
		"visible_on_front:", e.VisibleOnFront, "\n",
	)
	return int64(n2), err
}

// SalesOrderStatusStates represents a collection type for DB table
// sales_order_status_state
// Not thread safe. Auto generated.
type SalesOrderStatusStates struct {
	Data []*SalesOrderStatusState `json:"data,omitempty"`
}

// NewSalesOrderStatusStates  creates a new initialized collection. Auto
// generated.
func NewSalesOrderStatusStates() *SalesOrderStatusStates {
	return &SalesOrderStatusStates{
		Data: make([]*SalesOrderStatusState, 0, 5),
	}
}

// Append will add a new item at the end of * SalesOrderStatusStates . Auto
// generated via dmlgen.
func (cc *SalesOrderStatusStates) Append(n ...*SalesOrderStatusState) *SalesOrderStatusStates {
	cc.Data = append(cc.Data, n...)
	return cc
}

// Cut will remove items i through j-1. Auto generated via dmlgen.
func (cc *SalesOrderStatusStates) Cut(i, j int) *SalesOrderStatusStates {
	z := cc.Data // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this avoids the memory leak
	}
	z = z[:len(z)-j+i]
	cc.Data = z
	return cc
}

// Delete will remove an item from the slice. Auto generated via dmlgen.
func (cc *SalesOrderStatusStates) Delete(i int) *SalesOrderStatusStates {
	z := cc.Data // copy the slice header
	end := len(z) - 1
	cc.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	cc.Data = z
	return cc
}

// Each will run function f on all items in []* SalesOrderStatusState . Auto
// generated via dmlgen.
func (cc *SalesOrderStatusStates) Each(f func(*SalesOrderStatusState)) *SalesOrderStatusStates {
	if cc == nil {
		return nil
	}
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Filter filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (cc *SalesOrderStatusStates) Filter(f func(*SalesOrderStatusState) bool) *SalesOrderStatusStates {
	if cc == nil {
		return nil
	}
	b, i := cc.Data[:0], 0
	for _, e := range cc.Data {
		if f(e) {
			b = append(b, e)
		}
		i++
	}
	for i := len(b); i < len(cc.Data); i++ {
		cc.Data[i] = nil // this should avoid the memory leak
	}
	cc.Data = b
	return cc
}

// Insert will place a new item at position i. Auto generated via dmlgen.
func (cc *SalesOrderStatusStates) Insert(n *SalesOrderStatusState, i int) *SalesOrderStatusStates {
	z := cc.Data // copy the slice header
	z = append(z, &SalesOrderStatusState{})
	copy(z[i+1:], z[i:])
	z[i] = n
	cc.Data = z
	return cc
}

// Swap will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *SalesOrderStatusStates) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }

// Len will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *SalesOrderStatusStates) Len() int {
	if cc == nil {
		return 0
	}
	return len(cc.Data)
}

// Validate runs internal consistency tests on all items.
func (cc *SalesOrderStatusStates) Validate() (err error) {
	if len(cc.Data) == 0 {
		return nil
	}
	for i, ld := 0, len(cc.Data); i < ld && err == nil; i++ {
		err = cc.Data[i].Validate()
	}
	return
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (cc *SalesOrderStatusStates) WriteTo(w io.Writer) (n int64, err error) {
	for i, d := range cc.Data {
		n2, err := d.WriteTo(w)
		if err != nil {
			return 0, errors.Wrapf(err, "[dmltestgenerated2] WriteTo failed at index %d", i)
		}
		n += n2
	}
	return n, nil
}
