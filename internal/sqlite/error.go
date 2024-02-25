package sqlite

import (
	"slices"
	"strings"
)

const (
	CONSTRAINT_CHECK      = 275
	CONSTRAINT_COMMITHOOK = 531
	CONSTRAINT_DATATYPE   = 3091
	CONSTRAINT_FOREIGNKEY = 787
	CONSTRAINT_FUNCTION   = 1043
	CONSTRAINT_NOTNULL    = 1299
	CONSTRAINT_PINNED     = 2835
	CONSTRAINT_PRIMARYKEY = 1555
	CONSTRAINT_ROWID      = 2579
	CONSTRAINT_TRIGGER    = 1811
	CONSTRAINT_UNIQUE     = 2067
	CONSTRAINT_VTAB       = 2323
)

type Error struct {
	Code int
	Msg  string
}

func AsConstraintError(err error, code int, codes ...int) (ConstraintError, bool) {
	e, ok := AsError(err)
	if !(ok || e.Code == code || slices.Contains(codes, e.Code)) {
		return ConstraintError{}, false
	}
	return ConstraintError(e), true
}

type ConstraintError Error

func (e ConstraintError) IsField(field string) bool {
	return strings.Contains(string(e.Msg), field)
}
