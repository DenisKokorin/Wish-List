package models

type InvitationToken struct {
	UserID    int64
	GroupID   int64
	ExpiresAt int64
}
