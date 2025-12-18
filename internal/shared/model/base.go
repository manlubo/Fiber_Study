package model

import "time"

type BaseModel struct {
	Status    Status     `db:"status" json:"status"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type Status string

const (
	StatusReady    Status = "READY"
	StatusActive   Status = "ACTIVE"
	StatusDisabled Status = "DISABLED"
	StatusDeleted  Status = "DELETED"
)
