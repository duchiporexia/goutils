package xgenid

import (
	"github.com/duchiporexia/goutils/xerr"
	"time"
)

func repoGetOne() (*idGenerator, error) {
	var m idGenerator
	result := s.db.Where("expire_time <= ?", time.Now()).First(&m)
	return &m, result.Error
}

func repoUpdateVersion(id int, version int, resertVersion bool, expireInSeconds int) (newVersion int, expireTime time.Time, err error) {
	newVersion = version + 1
	if resertVersion {
		newVersion = 1
	}
	expireTime = time.Now().Add(time.Second * time.Duration(expireInSeconds))

	var m idGenerator
	m.Id = id
	result := s.db.Model(&m).Where("version = ?", version).Updates(idGenerator{Version: newVersion, ExpireTime: expireTime})

	if result.Error != nil {
		err = result.Error
		return
	}
	if result.RowsAffected != 1 {
		err = xerr.ErrNoRowsAffected
	}
	return
}
