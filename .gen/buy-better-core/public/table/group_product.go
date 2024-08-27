//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var GroupProduct = newGroupProductTable("public", "group_product", "")

type groupProductTable struct {
	postgres.Table

	// Columns
	ID   postgres.ColumnInteger
	Name postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type GroupProductTable struct {
	groupProductTable

	EXCLUDED groupProductTable
}

// AS creates new GroupProductTable with assigned alias
func (a GroupProductTable) AS(alias string) *GroupProductTable {
	return newGroupProductTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new GroupProductTable with assigned schema name
func (a GroupProductTable) FromSchema(schemaName string) *GroupProductTable {
	return newGroupProductTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new GroupProductTable with assigned table prefix
func (a GroupProductTable) WithPrefix(prefix string) *GroupProductTable {
	return newGroupProductTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new GroupProductTable with assigned table suffix
func (a GroupProductTable) WithSuffix(suffix string) *GroupProductTable {
	return newGroupProductTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newGroupProductTable(schemaName, tableName, alias string) *GroupProductTable {
	return &GroupProductTable{
		groupProductTable: newGroupProductTableImpl(schemaName, tableName, alias),
		EXCLUDED:          newGroupProductTableImpl("", "excluded", ""),
	}
}

func newGroupProductTableImpl(schemaName, tableName, alias string) groupProductTable {
	var (
		IDColumn       = postgres.IntegerColumn("id")
		NameColumn     = postgres.StringColumn("name")
		allColumns     = postgres.ColumnList{IDColumn, NameColumn}
		mutableColumns = postgres.ColumnList{NameColumn}
	)

	return groupProductTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:   IDColumn,
		Name: NameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
