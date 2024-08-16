package domain

import "time"

type User struct {
	ID         int64
	RoleID     int64
	Email      string
	Password   string
	Name       string
	LastAccess *time.Time
}

type Role struct {
	ID   int64
	Name string
}

type RoleRight struct {
	ID      int64
	RoleID  int64
	RCreate int
	RRead   int
	RUpdate int
	RDelete int
}
