package handler_test

import (
	"testing"

	"companies/handler"
	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	crud := handler.NewMongoCrud("mongodb://localhost:27017", "companies", "S3cret")
	crud.Connect()
	testObj := handler.Company{
		Name:              "test",
		Description:       "test",
		AmountOfEmployers: 1,
		Registered:        false,
		Type:              "Corporations",
	}
	testRes, err := crud.Create(testObj)
	assert.NoError(t, err)
	assert.NotEmpty(t, testRes.ID)
	objRead, errRead := crud.Read(testRes.Name)
	assert.NoError(t, errRead)
	assert.Equal(t, testObj.Name, objRead.Name)
	testObj.Description = "test Update"
	resUpd, errUpd := crud.Update(testObj)
	assert.NoError(t, errUpd)
	assert.Equal(t, resUpd, int64(1))
	resDel, errDel := crud.Delete(testObj.Name)
	assert.NoError(t, errDel)
	assert.Equal(t, resDel, int64(1))
}
