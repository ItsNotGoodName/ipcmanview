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

var ScanQueueTasks = newScanQueueTasksTable("dahua", "scan_queue_tasks", "")

type scanQueueTasksTable struct {
	postgres.Table

	// Columns
	CameraID postgres.ColumnInteger
	Kind     postgres.ColumnString
	Range    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ScanQueueTasksTable struct {
	scanQueueTasksTable

	EXCLUDED scanQueueTasksTable
}

// AS creates new ScanQueueTasksTable with assigned alias
func (a ScanQueueTasksTable) AS(alias string) *ScanQueueTasksTable {
	return newScanQueueTasksTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ScanQueueTasksTable with assigned schema name
func (a ScanQueueTasksTable) FromSchema(schemaName string) *ScanQueueTasksTable {
	return newScanQueueTasksTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ScanQueueTasksTable with assigned table prefix
func (a ScanQueueTasksTable) WithPrefix(prefix string) *ScanQueueTasksTable {
	return newScanQueueTasksTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ScanQueueTasksTable with assigned table suffix
func (a ScanQueueTasksTable) WithSuffix(suffix string) *ScanQueueTasksTable {
	return newScanQueueTasksTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newScanQueueTasksTable(schemaName, tableName, alias string) *ScanQueueTasksTable {
	return &ScanQueueTasksTable{
		scanQueueTasksTable: newScanQueueTasksTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newScanQueueTasksTableImpl("", "excluded", ""),
	}
}

func newScanQueueTasksTableImpl(schemaName, tableName, alias string) scanQueueTasksTable {
	var (
		CameraIDColumn = postgres.IntegerColumn("camera_id")
		KindColumn     = postgres.StringColumn("kind")
		RangeColumn    = postgres.StringColumn("range")
		allColumns     = postgres.ColumnList{CameraIDColumn, KindColumn, RangeColumn}
		mutableColumns = postgres.ColumnList{CameraIDColumn, KindColumn, RangeColumn}
	)

	return scanQueueTasksTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		CameraID: CameraIDColumn,
		Kind:     KindColumn,
		Range:    RangeColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
