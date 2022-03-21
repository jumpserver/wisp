package common

import "github.com/gofrs/uuid"

func IsUUID(sid string) bool {
	_, err := uuid.FromString(sid)
	return err == nil
}

func UUID() string {
	ret, _ := uuid.NewV4()
	return ret.String()
}
