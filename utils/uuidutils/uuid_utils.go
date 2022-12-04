package uuidutils

import uuid "github.com/google/uuid"

// NewUUID 返回随机生成的 UUID
func NewUUID() string {
	return uuid.New().String()
}
