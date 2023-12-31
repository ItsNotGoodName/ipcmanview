package migrations

import (
	"context"
	_ "embed"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

//go:embed normalize.sql
var normalizeSQL string

func Normalize(ctx context.Context, db repo.DB) error {
	_, err := db.ExecContext(ctx, normalizeSQL)
	if err != nil {
		return err
	}

	{
		c := dahua.NewFileCursor()
		err := db.NormalizeDahuaFileCursor(context.Background(), repo.NormalizeDahuaFileCursorParams{
			QuickCursor: c.QuickCursor,
			FullCursor:  c.FullCursor,
			FullEpoch:   c.FullEpoch,
			Scan:        c.Scan,
			ScanPercent: c.ScanPercent,
			ScanType:    c.ScanType,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
