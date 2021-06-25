package xgenid

import (
	"fmt"
	"github.com/duchiporexia/goutils/xconfig"
	"github.com/duchiporexia/goutils/xerr"
	"github.com/duchiporexia/goutils/xlog"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const retryMaxInterval = 3000

var (
	schemaName              string
	idGeneratorExpireInSecs int
)

type idGeneratorService struct {
	db                     *gorm.DB
	idGenerator            *idGenerator
	expireInSeconds        int
	disableAutoGetWorkerId bool
}

var s idGeneratorService

func Init(cfg *IdGeneratorConfig, db *gorm.DB) {
	if cfg == nil {
		cfg = &IdGeneratorConfig{
			SchemaName:   "",
			ExpireInSecs: "600",
		}
	}
	schemaName = cfg.SchemaName
	idGeneratorExpireInSecs = xconfig.GetInt(cfg.ExpireInSecs, 600)
	s.db = db
	s.expireInSeconds = idGeneratorExpireInSecs
	if !s.disableAutoGetWorkerId {
		ensureHasWorkerId()
		go refreshWorkerIdLoop()
	}
}

func waitFor(ms int, retry int) (int, int) {
	// -1 means reach the max
	if retry == -1 {
		return retryMaxInterval, -1
	}
	newMs := 100 + retry*100 + rand.Intn(100)
	if newMs >= retryMaxInterval {
		return retryMaxInterval, -1
	}
	return newMs, retry + 1
}

func ensureHasWorkerId() {
	var ms, retry int
	for getNewWorkerId() != nil {
		ms, retry = waitFor(ms, retry)
		xlog.Warn(fmt.Sprintf("Get worker id delayed %v ms", ms))
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func getNewWorkerId() error {
	var err error
	s.idGenerator, err = repoGetOne()
	if err != nil {
		return err
	}
	err = updateVersion(true)
	if err == nil {
		setServerId(int(s.idGenerator.Id))
	}
	return err
}

func updateVersion(resetVersion bool) error {
	newVersion, expireTime, err := repoUpdateVersion(s.idGenerator.Id, s.idGenerator.Version, resetVersion, s.expireInSeconds)
	if err != nil {
		return err
	}
	s.idGenerator.Version = newVersion
	s.idGenerator.ExpireTime = expireTime
	return nil
}

func nextTimeInMs() int {
	valueRange := s.expireInSeconds >> 2
	start := s.expireInSeconds - valueRange
	secs := start + rand.Intn(valueRange>>1)
	return secs * 1000
}

func refreshWorkerIdLoop() {
	ms := nextTimeInMs()
	xlog.Info(fmt.Sprintf("ServerId will be refreshed in %v ms", ms))
	time.Sleep(time.Millisecond * time.Duration(ms))

	var retry int
	for {
		err := refreshWorkerId()
		if err == nil {
			break
		}
		if err == xerr.ErrInvalidServerId || err == xerr.ErrNoRowsAffected {
			resetServerId()
			s.idGenerator = nil
			ensureHasWorkerId()
			break
		}

		ms, retry = waitFor(ms, retry)
		xlog.Warn(fmt.Sprintf("Refresh worker id failed, will retry in %v ms", ms))
		time.Sleep(time.Millisecond * time.Duration(ms))
	}

	refreshWorkerIdLoop()
}

func refreshWorkerId() error {
	if s.idGenerator == nil || s.idGenerator.ExpireTime.Before(time.Now()) {
		return xerr.ErrInvalidServerId
	}
	return updateVersion(false)
}
