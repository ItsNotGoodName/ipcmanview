package sqlite

import "github.com/ItsNotGoodName/ipcmanview/internal/repo"

// this is dumb

type DBTx struct {
	repo.DBTX
	*repo.Queries
}

func (db DB) C() DBTx {
	return DBTx{
		DBTX:    db,
		Queries: repo.New(db),
	}
}

func (tx Tx) C() DBTx {
	return DBTx{
		DBTX:    tx,
		Queries: repo.New(tx),
	}
}
