package recorderstorage

import (
	"github/jumpserver/wisp/pkg/jms-sdk-go/model"
)

func NewNullStorage() (storage NullStorage) {
	storage = NullStorage{}
	return
}

type NullStorage struct{}

func (f NullStorage) BulkSave(commands []*model.Command) (err error) {
	return
}

func (f NullStorage) Upload(gZipFile, target string) (err error) {

	return
}

func (f NullStorage) TypeName() string {
	return "null"
}
