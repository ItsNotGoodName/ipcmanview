//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	CameraDetails = CameraDetails.FromSchema(schema)
	CameraFiles = CameraFiles.FromSchema(schema)
	CameraLicenses = CameraLicenses.FromSchema(schema)
	CameraSoftwares = CameraSoftwares.FromSchema(schema)
	Cameras = Cameras.FromSchema(schema)
	ScanActiveTasks = ScanActiveTasks.FromSchema(schema)
	ScanCompleteTasks = ScanCompleteTasks.FromSchema(schema)
	ScanCursors = ScanCursors.FromSchema(schema)
	ScanQueueTasks = ScanQueueTasks.FromSchema(schema)
	ScanSeeds = ScanSeeds.FromSchema(schema)
}