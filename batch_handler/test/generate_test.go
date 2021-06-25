package test

import (
	"github.com/stretchr/testify/assert"
	bh "servicehub/common/batch_handler"
	"testing"
)

func TestGenerateAllHandlers(t *testing.T) {
	paramsType := "servicehub/common/batch_handler/test/dto.DogPo"
	valueType := "*servicehub/common/batch_handler/test/dto.Dog"

	assert.NoError(t, bh.Generate("TCreateNvHandler", "", paramsType, "", true, bh.HandlerTypeCreate))
	assert.NoError(t, bh.Generate("TCreateHvHandler", "", paramsType, valueType, true, bh.HandlerTypeCreate))

	assert.NoError(t, bh.Generate("TReadNpHandler", "string", "", valueType, true, bh.HandlerTypeRead))
	assert.NoError(t, bh.Generate("TReadHpHandler", "string", paramsType, valueType, true, bh.HandlerTypeRead))

	assert.NoError(t, bh.Generate("TUpdateNpNvHandler", "int", "", "", true, bh.HandlerTypeUpdate))
	assert.NoError(t, bh.Generate("TUpdateHpNvHandler", "int", paramsType, "", true, bh.HandlerTypeUpdate))
	assert.NoError(t, bh.Generate("TUpdateNpHvHandler", "int", "", valueType, true, bh.HandlerTypeUpdate))
	assert.NoError(t, bh.Generate("TUpdateHpHvHandler", "int", paramsType, valueType, true, bh.HandlerTypeUpdate))

	assert.NoError(t, bh.Generate("TDeleteNpNvHandler", "int", "", "", true, bh.HandlerTypeDelete))
	assert.NoError(t, bh.Generate("TDeleteHpNvHandler", "int", paramsType, "", true, bh.HandlerTypeDelete))
	assert.NoError(t, bh.Generate("TDeleteNpHvHandler", "int", "", valueType, true, bh.HandlerTypeDelete))
	assert.NoError(t, bh.Generate("TDeleteHpHvHandler", "int", paramsType, valueType, true, bh.HandlerTypeDelete))

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	assert.NoError(t, bh.Generate("SCreateNvHandler", "", paramsType, "", false, bh.HandlerTypeCreate))
	assert.NoError(t, bh.Generate("SCreateHvHandler", "", paramsType, valueType, false, bh.HandlerTypeCreate))

	assert.NoError(t, bh.Generate("SReadNpHandler", "string", "", valueType, false, bh.HandlerTypeRead))
	assert.NoError(t, bh.Generate("SReadHpHandler", "string", paramsType, valueType, false, bh.HandlerTypeRead))

	assert.NoError(t, bh.Generate("SUpdateNpNvHandler", "int", "", "", false, bh.HandlerTypeUpdate))
	assert.NoError(t, bh.Generate("SUpdateHpNvHandler", "int", paramsType, "", false, bh.HandlerTypeUpdate))
	assert.NoError(t, bh.Generate("SUpdateNpHvHandler", "int", "", valueType, false, bh.HandlerTypeUpdate))
	assert.NoError(t, bh.Generate("SUpdateHpHvHandler", "int", paramsType, valueType, false, bh.HandlerTypeUpdate))

	assert.NoError(t, bh.Generate("SDeleteNpNvHandler", "int", "", "", false, bh.HandlerTypeDelete))
	assert.NoError(t, bh.Generate("SDeleteHpNvHandler", "int", paramsType, "", false, bh.HandlerTypeDelete))
	assert.NoError(t, bh.Generate("SDeleteNpHvHandler", "int", "", valueType, false, bh.HandlerTypeDelete))
	assert.NoError(t, bh.Generate("SDeleteHpHvHandler", "int", paramsType, valueType, false, bh.HandlerTypeDelete))
}
