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

var CameraSoftwares = newCameraSoftwaresTable("dahua", "camera_softwares", "")

type cameraSoftwaresTable struct {
	postgres.Table

	// Columns
	CameraID                postgres.ColumnInteger
	Build                   postgres.ColumnString
	BuildDate               postgres.ColumnString
	SecurityBaseLineVersion postgres.ColumnString
	Version                 postgres.ColumnString
	WebVersion              postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CameraSoftwaresTable struct {
	cameraSoftwaresTable

	EXCLUDED cameraSoftwaresTable
}

// AS creates new CameraSoftwaresTable with assigned alias
func (a CameraSoftwaresTable) AS(alias string) *CameraSoftwaresTable {
	return newCameraSoftwaresTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CameraSoftwaresTable with assigned schema name
func (a CameraSoftwaresTable) FromSchema(schemaName string) *CameraSoftwaresTable {
	return newCameraSoftwaresTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CameraSoftwaresTable with assigned table prefix
func (a CameraSoftwaresTable) WithPrefix(prefix string) *CameraSoftwaresTable {
	return newCameraSoftwaresTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CameraSoftwaresTable with assigned table suffix
func (a CameraSoftwaresTable) WithSuffix(suffix string) *CameraSoftwaresTable {
	return newCameraSoftwaresTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCameraSoftwaresTable(schemaName, tableName, alias string) *CameraSoftwaresTable {
	return &CameraSoftwaresTable{
		cameraSoftwaresTable: newCameraSoftwaresTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newCameraSoftwaresTableImpl("", "excluded", ""),
	}
}

func newCameraSoftwaresTableImpl(schemaName, tableName, alias string) cameraSoftwaresTable {
	var (
		CameraIDColumn                = postgres.IntegerColumn("camera_id")
		BuildColumn                   = postgres.StringColumn("build")
		BuildDateColumn               = postgres.StringColumn("build_date")
		SecurityBaseLineVersionColumn = postgres.StringColumn("security_base_line_version")
		VersionColumn                 = postgres.StringColumn("version")
		WebVersionColumn              = postgres.StringColumn("web_version")
		allColumns                    = postgres.ColumnList{CameraIDColumn, BuildColumn, BuildDateColumn, SecurityBaseLineVersionColumn, VersionColumn, WebVersionColumn}
		mutableColumns                = postgres.ColumnList{CameraIDColumn, BuildColumn, BuildDateColumn, SecurityBaseLineVersionColumn, VersionColumn, WebVersionColumn}
	)

	return cameraSoftwaresTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		CameraID:                CameraIDColumn,
		Build:                   BuildColumn,
		BuildDate:               BuildDateColumn,
		SecurityBaseLineVersion: SecurityBaseLineVersionColumn,
		Version:                 VersionColumn,
		WebVersion:              WebVersionColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
