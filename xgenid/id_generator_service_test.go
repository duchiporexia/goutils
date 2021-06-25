package xgenid

import (
	"database/sql"
	"github.com/duchiporexia/goutils/xlog"
	"github.com/duchiporexia/goutils/xutils/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	testutils.RegisterDBForTest()
	s.disableAutoGetWorkerId = true
}

func TestIdGeneratorService(t *testing.T) {
	xlog.Init("test")
	db, err := sql.Open("txdb", "identifier")
	defer db.Close()
	assert.NoError(t, err)
	err = testutils.InitIdGeneratorTable(db)
	assert.NoError(t, err)

	gdb := testutils.NewGormDBForTesting()
	Init(nil, gdb)
	err = getNewWorkerId()
	assert.NoError(t, err)
	assert.Equal(t, s.idGenerator.Id, 1)
	err = getNewWorkerId()
	assert.NoError(t, err)
	assert.Equal(t, s.idGenerator.Id, 2)
}
