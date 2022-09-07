package UuidUtils

import uuid "github.com/satori/go.uuid"

// NewUUID 返回随机生成的 UUID
func NewUUID() string {
	return uuid.NewV4().String()
}

// NewSpaceUUID 返回指定命名的 UUID
func NewSpaceUUID(space string) string {
	return uuid.NewV5(uuid.NewV4(), space).String()
}
