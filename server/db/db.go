package db

import (
	"database/sql"
	"errors"
	"fmt"
	formatter "server/util"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

const P_NAME = "Name"
const P_AGE = "Age"
const P_SEX = "Sex"
const P_HEIGHT = "Height"

type MemberInfo struct {
	Id     int
	Name   string
	Age    int
	Sex    bool
	Height decimal.Decimal
}

func InitMySQL() *sql.DB {
	db, errSql := sql.Open("mysql", "steven:123456@/gotest") // account:password@/dbname
	if errSql != nil {
		panic(errSql)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	fmt.Println("After MySQL init")

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS member (id int NOT NULL AUTO_INCREMENT, name VARCHAR(10), age int, sex BINARY, height DECIMAL(4,2), PRIMARY KEY(id))")
	if err != nil {
		fmt.Println("create member table err:", err)
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS jwt (account VARCHAR(10), token VARCHAR(1024), PRIMARY KEY(account))")
	if err != nil {
		fmt.Println("create jwt table err:", err)
		panic(err)
	}

	testMySQL()
	return db
}

func testMySQL() {
	// Query
	// rows, err := DB.Query("SELECT * FROM member WHERE name = ?", "Bob")
	// if err != nil {
	// 	fmt.Println("db exception")
	// 	// 教學文件說要特別例外處理ErrNoRows但測試不會收到
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("empty row")
	// 	} else {
	// 		fmt.Println("exception", err)
	// 	}
	// 	return
	// }

	// defer rows.Close()
	// memberInfo := MemberInfo{}
	// for rows.Next() {
	// 	errScan := rows.Scan(&memberInfo.Id, &memberInfo.Name, &memberInfo.Age, &memberInfo.Sex, &memberInfo.Height)
	// 	if errScan != nil {
	// 		fmt.Println("exception2 ", errScan)
	// 	}
	// 	fmt.Println(memberInfo)

	// 	response := BaseResponse[MemberInfo]{Status: 1, Data: memberInfo}
	// 	fmt.Println("baseResponse: ", response)
	// 	b, err := json.Marshal(response)
	// 	if err != nil {
	// 		fmt.Println("to json string fail", err)
	// 	} else {
	// 		fmt.Println("b len:", len(b))
	// 		fmt.Println("Get Json String:", string(b))
	// 	}
	// }

	// INSERT
	// result, err := DB.Exec("INSERT INTO member (name, age, height, sex) VALUES (?, ?, ?, ?)", "Ken2", 30, 160.3, false)
	// if err != nil {
	// 	fmt.Println("insert error: ", err)
	// } else {
	// 	row, resultErr := result.RowsAffected()
	// 	if resultErr != nil {
	// 		fmt.Println("RowsAffected err", resultErr)
	// 	} else {
	// 		fmt.Println("row effect row: ", row)
	// 	}
	// }

	// Update
	// result, err := DB.Exec("UPDATE member SET name=? where id=?", "new name", 2)
	// if err != nil {
	// 	fmt.Println("Update Exec fail:", err)
	// } else {
	// 	row, resultErr := result.RowsAffected()
	// 	if resultErr != nil {
	// 		fmt.Println("RowsAffected error", resultErr)
	// 		return
	// 	} else {
	// 		fmt.Println("row effect row: ", row)
	// 	}
	// }
}

func Update(db *sql.DB, id, name, age, height, sex string) (int64, error) {
	var result sql.Result
	var err error
	if name != "" {
		result, err = db.Exec("UPDATE member SET name=? where id=?", name, id)
	} else if age != "" {
		result, err = db.Exec("UPDATE member SET age=? where id=?", age, id)
	} else if height != "" {
		result, err = db.Exec("UPDATE member SET height=? where id=?", height, id)
	} else if sex != "" {
		sb, _ := formatter.StringToBool(age)
		result, err = db.Exec("UPDATE member SET sex=? where id=?", sb, id)
	} else {
		return -1, errors.New("no available data to update")
	}

	if err != nil {
		return -1, err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return row, nil
}

func Insert(db *sql.DB, member *MemberInfo) (int64, error) {
	result, err := db.Exec("INSERT INTO member (name, age, height, sex) VALUES (?, ?, ?, ?)", member.Name, member.Age, member.Height, member.Sex)
	if err != nil {
		fmt.Println("Insert Exec fail:", err)
		return -1, err
	}

	row, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Insert RowsAffected fail:", err)
		return -1, err
	}

	fmt.Println("row effect row: ", row)
	return row, nil
}

func Del(db *sql.DB, id string) (int64, error) {
	result, err := db.Exec("DELETE FROM member WHERE id = ?", id)
	if err != nil {
		fmt.Println("Exec del err:", err)
		return -1, err
	}

	row, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Delete RowsAffected fail:", err)
		return -1, err
	} else {
		fmt.Println("Delete effect row: ", row)
		return row, nil
	}
}

func Get(db *sql.DB, name string, age *int64, sex *bool, height *decimal.Decimal) ([]MemberInfo, error) {
	var rows *sql.Rows
	var err error

	if len(name) > 0 {
		rows, err = db.Query("SELECT * FROM member WHERE name=?", name)
	} else if age != nil {
		rows, err = db.Query("SELECT * FROM member WHERE age=?", age)
	} else if sex != nil {
		rows, err = db.Query("SELECT * FROM member WHERE sex=?", sex)
	} else if height != nil {
		rows, err = db.Query("SELECT * FROM member WHERE height=?", height)
	} else {
		fmt.Println("no available query")
		return nil, errors.New("no available query")
	}

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		fmt.Println("Query exec fail: ", err)
		return nil, err
	}

	var data []MemberInfo
	for rows.Next() {
		memberInfo := MemberInfo{}
		err = rows.Scan(&memberInfo.Id, &memberInfo.Name, &memberInfo.Age, &memberInfo.Sex, &memberInfo.Height)
		if err != nil {
			fmt.Println("Query scan fail: ", err)
			return nil, err
		}
		data = append(data, memberInfo)
	}

	return data, nil
}

var jwtCacheMap = make(map[string]string)

func CheckJwt(db *sql.DB, account, tokenString string) bool {
	cacheToken, isExist := jwtCacheMap[account]
	if isExist {
		fmt.Printf("Find %s token in cache\n", account)
		return tokenString == cacheToken
	}

	rows, err := db.Query("SELECT token FROM jwt WHERE account=?", account)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		fmt.Printf("Query jwt table fail with account %s, %v", account, err)
		return false
	}

	result := false
	for rows.Next() {
		dbToken := ""
		rows.Scan(&dbToken)
		fmt.Println("dbToken ", dbToken)
		result = dbToken == tokenString
	}
	return result
}

func InsertJwt(db *sql.DB, account, tokenString string) (int64, error) {
	result, err := db.Exec("INSERT INTO jwt (account, token) VALUES (?, ?) ON DUPLICATE KEY UPDATE token=?", account, tokenString, tokenString)
	if err != nil {
		fmt.Println("Insert jwt Exec fail:", err)
		return -1, err
	}

	row, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Insert jwt RowsAffected fail:", err)
		return -1, err
	}

	jwtCacheMap[account] = tokenString
	return row, nil
}
