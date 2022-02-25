package common

import "github.com/gofrs/uuid"

func IsUUID(sid string) bool {
	_, err := uuid.FromString(sid)
	return err == nil
}
