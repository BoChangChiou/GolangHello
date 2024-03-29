package db

import (
	"database/sql"
	"log"
	"regexp"
	formatter "server/util"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFunNeedStartWithTest(t *testing.T) {
	println("abcde")
}

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func newRowForTestQuery(id int, name string, age int64, sex bool, height decimal.Decimal) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"Id", "Name", "Age", "Sex", "Height"}).AddRow(id, name, age, sex, height)
}

func TestQuery(t *testing.T) {
	db, mock := newMock()
	defer func() {
		db.Close()
	}()

	type Result struct {
		Data []MemberInfo
		Err  error
	}

	sName := "Steve"
	var sAge int64 = 33
	sSex := true
	sId := 1
	sHeight, _ := decimal.NewFromString("173.3")

	results := []Result{}
	result := Result{}

	// we should new sqlmock.Rows each time because sqlmock.Rows instance cannot be reused
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM member WHERE name=?")).WithArgs(sName).WillReturnRows(newRowForTestQuery(sId, sName, sAge, sSex, sHeight))
	result.Data, result.Err = Get(db, "Steve", &sAge, &sSex, &sHeight)
	results = append(results, result)

	result = Result{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM member WHERE age=?")).WithArgs(sAge).WillReturnRows(newRowForTestQuery(sId, sName, sAge, sSex, sHeight))
	result.Data, result.Err = Get(db, "", &sAge, &sSex, &sHeight)
	results = append(results, result)

	result = Result{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM member WHERE sex=?")).WithArgs(sSex).WillReturnRows(newRowForTestQuery(sId, sName, sAge, sSex, sHeight))
	result.Data, result.Err = Get(db, "", nil, &sSex, &sHeight)
	results = append(results, result)

	result = Result{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM member WHERE height=?")).WithArgs(sHeight).WillReturnRows(newRowForTestQuery(sId, sName, sAge, sSex, sHeight))
	result.Data, result.Err = Get(db, "", nil, nil, &sHeight)
	results = append(results, result)

	for _, r := range results {
		assert.NotNil(t, r.Data)
		assert.Equal(t, len(r.Data), sId)
		assert.Equal(t, r.Data[0].Name, sName)
		assert.Equal(t, r.Data[0].Age, formatter.Int64ToInt(sAge))
		assert.Equal(t, r.Data[0].Sex, sSex)
		assert.Equal(t, r.Data[0].Height, sHeight)
		assert.NoError(t, r.Err)
	}

	result = Result{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM member WHERE name=?")).WithArgs("Aaa").WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "Age", "Sex", "Height"}))
	result.Data, result.Err = Get(db, "Aaa", nil, nil, nil)
	assert.Nil(t, result.Data)
	assert.NoError(t, result.Err)
}
