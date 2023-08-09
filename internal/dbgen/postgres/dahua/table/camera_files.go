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

var CameraFiles = newCameraFilesTable("dahua", "camera_files", "")

type cameraFilesTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	CameraID  postgres.ColumnInteger
	FilePath  postgres.ColumnString
	Kind      postgres.ColumnString
	Size      postgres.ColumnInteger
	StartTime postgres.ColumnTimestampz
	EndTime   postgres.ColumnTimestampz
	Duration  postgres.ColumnInteger
	ScannedAt postgres.ColumnTimestampz
	Events    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CameraFilesTable struct {
	cameraFilesTable

	EXCLUDED cameraFilesTable
}

// AS creates new CameraFilesTable with assigned alias
func (a CameraFilesTable) AS(alias string) *CameraFilesTable {
	return newCameraFilesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CameraFilesTable with assigned schema name
func (a CameraFilesTable) FromSchema(schemaName string) *CameraFilesTable {
	return newCameraFilesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CameraFilesTable with assigned table prefix
func (a CameraFilesTable) WithPrefix(prefix string) *CameraFilesTable {
	return newCameraFilesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CameraFilesTable with assigned table suffix
func (a CameraFilesTable) WithSuffix(suffix string) *CameraFilesTable {
	return newCameraFilesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCameraFilesTable(schemaName, tableName, alias string) *CameraFilesTable {
	return &CameraFilesTable{
		cameraFilesTable: newCameraFilesTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newCameraFilesTableImpl("", "excluded", ""),
	}
}

func newCameraFilesTableImpl(schemaName, tableName, alias string) cameraFilesTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		CameraIDColumn  = postgres.IntegerColumn("camera_id")
		FilePathColumn  = postgres.StringColumn("file_path")
		KindColumn      = postgres.StringColumn("kind")
		SizeColumn      = postgres.IntegerColumn("size")
		StartTimeColumn = postgres.TimestampzColumn("start_time")
		EndTimeColumn   = postgres.TimestampzColumn("end_time")
		DurationColumn  = postgres.IntegerColumn("duration")
		ScannedAtColumn = postgres.TimestampzColumn("scanned_at")
		EventsColumn    = postgres.StringColumn("events")
		allColumns      = postgres.ColumnList{IDColumn, CameraIDColumn, FilePathColumn, KindColumn, SizeColumn, StartTimeColumn, EndTimeColumn, DurationColumn, ScannedAtColumn, EventsColumn}
		mutableColumns  = postgres.ColumnList{CameraIDColumn, FilePathColumn, KindColumn, SizeColumn, StartTimeColumn, EndTimeColumn, ScannedAtColumn, EventsColumn}
	)

	return cameraFilesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		CameraID:  CameraIDColumn,
		FilePath:  FilePathColumn,
		Kind:      KindColumn,
		Size:      SizeColumn,
		StartTime: StartTimeColumn,
		EndTime:   EndTimeColumn,
		Duration:  DurationColumn,
		ScannedAt: ScannedAtColumn,
		Events:    EventsColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
