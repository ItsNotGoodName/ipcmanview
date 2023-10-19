package auth

import (
	"crypto/rand"
	"fmt"
	"net/netip"
	"time"
)

type Session struct {
	UserID     int64
	Token      string
	ClientID   string
	IpAddress  netip.Addr
	IssuedAt   time.Time
	LastUsedAt time.Time
	ExpiredAt  time.Time
}

func (s Session) Valid() bool {
	return s.ExpiredAt.After(time.Now())
}

func NewSession(clientID string, ipAddress string, userID int64, duration time.Duration) (Session, error) {
	ip, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return Session{}, err
	}

	now := time.Now()

	return Session{
		UserID:     userID,
		Token:      NewSessionToken(),
		ClientID:   clientID,
		IpAddress:  ip,
		IssuedAt:   now,
		LastUsedAt: now,
		ExpiredAt:  now.Add(duration),
	}, nil
}

func NewSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
