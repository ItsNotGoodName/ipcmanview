package auth

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dbgen/postgres/public/model"
	. "github.com/ItsNotGoodName/ipcmanview/internal/dbgen/postgres/public/table"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	. "github.com/go-jet/jet/v2/postgres"
)

type dbT struct{}

var DB dbT

var dbSessionProjection ProjectionList = []Projection{
	UserSessions.UserID.AS("user_id"),
	UserSessions.Token.AS("token"),
	UserSessions.ClientID.AS("client_id"),
	UserSessions.IPAddress.AS("ip_address"),
	UserSessions.ExpiredAt.AS("expired_at"),
	UserSessions.LastUsedAt.AS("last_used_at"),
	UserSessions.IssuedAt.AS("issued_at"),
}

func (dbT) SessionCreate(ctx context.Context, db qes.Querier, r Session) (Session, error) {
	var session Session
	err := qes.ScanOne(ctx, db, &session, UserSessions.
		INSERT(
			UserSessions.UserID,
			UserSessions.Token,
			UserSessions.ClientID,
			UserSessions.IPAddress,
			UserSessions.IssuedAt,
			UserSessions.ExpiredAt,
			UserSessions.LastUsedAt,
		).
		MODEL(model.UserSessions{
			UserID:     int32(r.UserID),
			Token:      r.Token,
			ClientID:   r.ClientID,
			IPAddress:  r.IpAddress.String(),
			ExpiredAt:  r.ExpiredAt,
			LastUsedAt: r.LastUsedAt,
			IssuedAt:   r.IssuedAt,
		}).
		RETURNING(dbSessionProjection))
	return session, err
}

func (dbT) SessionGet(ctx context.Context, db qes.Querier, sessionToken string) (Session, error) {
	var session Session
	err := qes.ScanOne(ctx, db, &session, UserSessions.
		UPDATE(UserSessions.LastUsedAt).
		SET(UserSessions.LastUsedAt.SET(TimestampzT(time.Now()))).
		WHERE(UserSessions.Token.EQ(String(sessionToken))).
		RETURNING(dbSessionProjection))
	return session, err
}

func (dbT) SessionRemove(ctx context.Context, db qes.Querier, sessionToken string) error {
	_, err := qes.Exec(ctx, db, UserSessions.
		DELETE().
		WHERE(UserSessions.Token.EQ(String(sessionToken))))
	return err
}
