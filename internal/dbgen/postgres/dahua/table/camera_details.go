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

var CameraDetails = newCameraDetailsTable("dahua", "camera_details", "")

type cameraDetailsTable struct {
	postgres.Table

	// Columns
	CameraID        postgres.ColumnInteger
	Sn              postgres.ColumnString
	DeviceClass     postgres.ColumnString
	DeviceType      postgres.ColumnString
	HardwareVersion postgres.ColumnString
	MarketArea      postgres.ColumnString
	ProcessInfo     postgres.ColumnString
	Vendor          postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CameraDetailsTable struct {
	cameraDetailsTable

	EXCLUDED cameraDetailsTable
}

// AS creates new CameraDetailsTable with assigned alias
func (a CameraDetailsTable) AS(alias string) *CameraDetailsTable {
	return newCameraDetailsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CameraDetailsTable with assigned schema name
func (a CameraDetailsTable) FromSchema(schemaName string) *CameraDetailsTable {
	return newCameraDetailsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CameraDetailsTable with assigned table prefix
func (a CameraDetailsTable) WithPrefix(prefix string) *CameraDetailsTable {
	return newCameraDetailsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CameraDetailsTable with assigned table suffix
func (a CameraDetailsTable) WithSuffix(suffix string) *CameraDetailsTable {
	return newCameraDetailsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCameraDetailsTable(schemaName, tableName, alias string) *CameraDetailsTable {
	return &CameraDetailsTable{
		cameraDetailsTable: newCameraDetailsTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newCameraDetailsTableImpl("", "excluded", ""),
	}
}

func newCameraDetailsTableImpl(schemaName, tableName, alias string) cameraDetailsTable {
	var (
		CameraIDColumn        = postgres.IntegerColumn("camera_id")
		SnColumn              = postgres.StringColumn("sn")
		DeviceClassColumn     = postgres.StringColumn("device_class")
		DeviceTypeColumn      = postgres.StringColumn("device_type")
		HardwareVersionColumn = postgres.StringColumn("hardware_version")
		MarketAreaColumn      = postgres.StringColumn("market_area")
		ProcessInfoColumn     = postgres.StringColumn("process_info")
		VendorColumn          = postgres.StringColumn("vendor")
		allColumns            = postgres.ColumnList{CameraIDColumn, SnColumn, DeviceClassColumn, DeviceTypeColumn, HardwareVersionColumn, MarketAreaColumn, ProcessInfoColumn, VendorColumn}
		mutableColumns        = postgres.ColumnList{CameraIDColumn, SnColumn, DeviceClassColumn, DeviceTypeColumn, HardwareVersionColumn, MarketAreaColumn, ProcessInfoColumn, VendorColumn}
	)

	return cameraDetailsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		CameraID:        CameraIDColumn,
		Sn:              SnColumn,
		DeviceClass:     DeviceClassColumn,
		DeviceType:      DeviceTypeColumn,
		HardwareVersion: HardwareVersionColumn,
		MarketArea:      MarketAreaColumn,
		ProcessInfo:     ProcessInfoColumn,
		Vendor:          VendorColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
