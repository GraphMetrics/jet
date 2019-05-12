package tests

import (
	"github.com/sub0zero/go-sqlbuilder/sqlbuilder"
	"github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/dvds/model"
	"gotest.tools/assert"
	"strings"
	"testing"
	"time"
)

func assertQuery(t *testing.T, query sqlbuilder.Statement, expectedQuery string, expectedArgs ...interface{}) {
	_, args, err := query.Sql()
	assert.NilError(t, err)
	//assert.Equal(t, queryStr, expectedQuery)
	assert.DeepEqual(t, args, expectedArgs)

	debuqSql, err := query.DebugSql()
	assert.NilError(t, err)
	assert.Equal(t, debuqSql, expectedQuery, args)
}

func int16Ptr(i int16) *int16 {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func timeWithoutTimeZone(t string, precision int) *time.Time {

	precisionStr := ""

	if precision > 0 {
		precisionStr = "." + strings.Repeat("9", precision)
	}

	time, err := time.Parse("2006-01-02 15:04:05"+precisionStr+" +0000", t+" +0000")

	if err != nil {
		panic(err)
	}

	return &time
}

var customer0 = model.Customer{
	CustomerID: 1,
	StoreID:    1,
	FirstName:  "Mary",
	LastName:   "Smith",
	Email:      stringPtr("mary.smith@sakilacustomer.org"),
	Address:    nil,
	Activebool: true,
	CreateDate: *timeWithoutTimeZone("2006-02-14 00:00:00", 0),
	LastUpdate: timeWithoutTimeZone("2013-05-26 14:49:45.738", 3),
	Active:     int32Ptr(1),
}

var customer1 = model.Customer{
	CustomerID: 2,
	StoreID:    1,
	FirstName:  "Patricia",
	LastName:   "Johnson",
	Email:      stringPtr("patricia.johnson@sakilacustomer.org"),
	Address:    nil,
	Activebool: true,
	CreateDate: *timeWithoutTimeZone("2006-02-14 00:00:00", 0),
	LastUpdate: timeWithoutTimeZone("2013-05-26 14:49:45.738", 3),
	Active:     int32Ptr(1),
}

var lastCustomer = model.Customer{
	CustomerID: 599,
	StoreID:    2,
	FirstName:  "Austin",
	LastName:   "Cintron",
	Email:      stringPtr("austin.cintron@sakilacustomer.org"),
	Address:    nil,
	Activebool: true,
	CreateDate: *timeWithoutTimeZone("2006-02-14 00:00:00", 0),
	LastUpdate: timeWithoutTimeZone("2013-05-26 14:49:45.738", 3),
	Active:     int32Ptr(1),
}
