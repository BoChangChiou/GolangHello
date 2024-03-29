package http

import (
	"context"
	"database/sql"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"server/db"
	"server/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shopspring/decimal"
)

const RES_SUCCESS = 1
const RES_FAIL = 0
const RES_JWT_FAIL = 2

var secretKey = []byte("golang_abcdef_123456")

// -----BEGIN PUBLIC KEY-----
// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLaxSRB24P7C4pzeSReaQ8IdIb
// 7os/7d9rJADxw63B89tMWn4B6PNkvlWhHgOPqWY/5IeGRRfWdvrb8qnQuFCoHbMQ
// vpxYhVxni/z1RCvg3CTXuW38CXyDBBPMU1G9cEJ7bR+m8o5X6zh8YC/1siSHKINj
// cT9ZnK4dNvW7jN1c1wIDAQAB
// -----END PUBLIC KEY-----

const privateKey string = `-----BEGIN PRIVATE KEY-----
MIICXAIBAAKBgQDLaxSRB24P7C4pzeSReaQ8IdIb7os/7d9rJADxw63B89tMWn4B
6PNkvlWhHgOPqWY/5IeGRRfWdvrb8qnQuFCoHbMQvpxYhVxni/z1RCvg3CTXuW38
CXyDBBPMU1G9cEJ7bR+m8o5X6zh8YC/1siSHKINjcT9ZnK4dNvW7jN1c1wIDAQAB
AoGAcvbWzcyEMK2LvYamym0UHAQFSlH8EypuHZBglELCPh6C71kpZAzzGhnULVXY
L2ZO6odO7Ny5xzTBPHOd899nfUjmUYYVFxJF7g6i0kA4y7C5GWM9JMvIX7qqji1G
8i71pdvXIC0OXEieO1XHg60aEMxVPti4JW3NGRK4O/npI1kCQQDOb+fg5DjPgiRM
JsS/zN2Ng5/WQMfoOwKvMRQ6/kkuo+rThTC0ce6rvh/ihpL4yHv72hPfC4G8upTp
/TgazAJFAkEA/EGhFZMnCMAuzQO5YYW5v/0puuRTvoWh9NbvlTDS20QkJW0CyFN8
0mUfRkJQ8zKoUIoTbuekN7YqNr9rJnJiawJARauS0F11puK/KUw0Pp7/btErUn3O
edvgjgu8TiSfwjPj/rsGsv94k1G5JRRR6dCPt3HkHvSdNnqp40ZodvK/GQJBAJ8a
5fctsVkbnmlBCBQyvE4T59YxXYC12MkNKF/5Q4V5HTNd5ntj7T7m+SrfeR9rvC3Q
aSFyiWl6RHXzlinRy7ECQD5UyagQxFhg4hyZlB5/IX0dUAMwZcpJv3DShfIiHuIP
dcUeF3RysdIo8pZo3FtBKDzGgr9FRqkEmrIHZshrKvQ=
-----END PRIVATE KEY-----`

type Env struct {
	db *sql.DB
}

var StopChannel = make(chan interface{})

func StartServer(ctx context.Context) {
	db := db.InitMySQL()

	err := db.Ping()
	if err != nil {
		panic("open DB err")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	env := &Env{db: db}

	r.Get("/", defaultPath)
	r.Get("/sqlGet", env.sqlGet)
	r.Get("/ws", myws)
	r.Get("/stop", env.stop)
	r.Post("/login", env.login)
	r.Post("/sqlInsert", env.sqlInsert)
	r.Patch("/sqlUpdate/{id}", env.sqlUpdate)
	r.Delete("/sqlDelete/{id}", env.sqlDel)

	srv := &http.Server{
		Addr:         ":9090", //設定監聽的埠
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
	}

	// go httpApiClientGet(ctx)
	go wsHub.run()
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func checkJwtToken(env *Env, w http.ResponseWriter, r *http.Request) bool {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		fmt.Println("No JWT token")
		writeNormalFail(w, "No JWT token")
		return false
	}
	tokenString = tokenString[len("Bearer "):]

	// check JWT is expire or not
	token, claims, err := util.VerifyJwtToken(tokenString, secretKey)
	if err != nil {
		fmt.Println("JWT verify err: ", err)
		writeJwtTokenFail(w)
		return false
	}

	if !token.Valid {
		fmt.Println("JWT expire")
		writeJwtTokenFail(w)
		return false
	}

	account, isExist := claims["account"]
	if !isExist {
		fmt.Println("JWT doesn't cotain account")
		writeJwtTokenFail(w)
		return false
	}

	if !db.CheckJwt(env.db, account.(string), tokenString) {
		fmt.Println("JWT check in memory or db fail")
		writeJwtTokenFail(w)
		return false
	}

	return true
}

func defaultPath(w http.ResponseWriter, r *http.Request) {
	printRequest(r)
	fmt.Fprintf(w, "Hello astaxie!") //這個寫入到 w 的是輸出到客戶端的
}

func (env *Env) login(w http.ResponseWriter, r *http.Request) {
	printRequest(r)

	eAccount := r.PostFormValue("Account")
	ePassword := r.PostFormValue("Password")

	if eAccount == "" || ePassword == "" {
		writeNormalFail(w, "沒有帳號或密碼")
		return
	}

	account, err := util.RsaDecryptBase(eAccount, privateKey)
	fmt.Println("get Account ", account)
	if err != nil {
		fmt.Println("get Account fail", err)
		writeNormalFail(w, "帳號或密碼錯誤1")
		return
	}

	password, err := util.RsaDecryptBase(ePassword, privateKey)
	fmt.Println("get Password ", password)
	if err != nil {
		fmt.Println("get Password fail", err)
		writeNormalFail(w, "帳號或密碼錯誤2")
		return
	}

	if account != "Steven" || password != "123456" {
		writeNormalFail(w, "帳號或密碼錯誤3")
		return
	}
	token, err := util.CreateJwtToken(account, secretKey)
	if err != nil {
		fmt.Println("createToken fail: ", err)
		writeNormalFail(w, "登入失敗")
		return
	}

	rows, err := db.InsertJwt(env.db, account, token)
	if err != nil {
		fmt.Println("Insert Jwt fail", err)
		writeNormalFail(w, "登入失敗，寫入JWT失敗1")
		return
	}
	if rows == 0 {
		fmt.Println("Insert Jwt fail2")
		writeNormalFail(w, "登入失敗，寫入JWT失敗2")
		return
	}
	writeSuccessResponse(w, token)
}

type BaseResponse struct {
	Status int
	Data   any    `json:",omitempty"`
	Msg    string `json:",omitempty"`
}

// 從Body(Json)取得JSON資料轉成struct
func (env *Env) sqlInsert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sqlInsert")
	printRequest(r)

	if !checkJwtToken(env, w, r) {
		return
	}

	var member db.MemberInfo
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		println(err)
		writeNormalFail(w, "Insert bad data")
		return
	}

	if rows, err := db.Insert(env.db, &member); err != nil {
		println(err)
		writeNormalFail(w, "Insert fail")
		return
	} else {
		writeSuccessResponse(w, rows)
	}
}

func writeSuccessResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BaseResponse{Status: RES_SUCCESS, Data: data})
}

func writeNormalFail(w http.ResponseWriter, msg string) {
	writeFailResponse(w, RES_FAIL, msg)
}

func writeJwtTokenFail(w http.ResponseWriter) {
	writeFailResponse(w, RES_JWT_FAIL, "JWT token fail")
}

func writeFailResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BaseResponse{Status: status, Msg: msg})
}

// 從url path取得id
func (env *Env) sqlDel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sqlDel")
	printRequest(r)

	if !checkJwtToken(env, w, r) {
		return
	}

	id := r.PathValue("id")
	_, idError := util.StringToInt64(id)
	if idError != nil {
		println(idError)
		writeNormalFail(w, "ID is not allow")
		return
	}

	result, err := db.Del(env.db, id)
	if err != nil {
		println(err)
		writeNormalFail(w, "Delete fail")
		return
	}
	writeSuccessResponse(w, result)
}

// 從Body(Form Data)取得更新資訊
func (env *Env) sqlUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sqlUpdate")
	printRequest(r)

	if !checkJwtToken(env, w, r) {
		return
	}

	var id = r.PathValue("id")
	if _, idError := util.StringToInt64(id); idError != nil {
		println(idError)
		writeNormalFail(w, "ID is not allow")
		return
	}

	var name = r.PostFormValue(db.P_NAME)
	var age = r.PostFormValue(db.P_AGE)
	var height = r.PostFormValue(db.P_HEIGHT)
	var sex = r.PostFormValue(db.P_SEX)

	var err error
	var row int64
	if name != "" {
		row, err = db.Update(env.db, id, name, "", "", "")
	} else if age != "" {
		if _, intErr := util.StringToInt64(age); intErr != nil {
			println(intErr)
			writeNormalFail(w, "Parse age parse fail")
			return
		}
		row, err = db.Update(env.db, id, "", age, "", "")
	} else if height != "" {
		if _, floatErr := util.StringToFloat64(height); floatErr != nil {
			println(floatErr)
			writeNormalFail(w, "Parse height parse fail")
			return
		}
		row, err = db.Update(env.db, id, "", "", height, "")
	} else if sex != "" {
		if _, sexErr := util.StringToBool(sex); sexErr != nil {
			println(sexErr)
			writeNormalFail(w, "Parse sex parse fail")
			return
		}
		row, err = db.Update(env.db, id, "", "", "", sex)
	} else {
		writeNormalFail(w, "No available data to update")
		return
	}

	if err != nil {
		println(err)
		writeNormalFail(w, "Update fail")
		return
	}
	writeSuccessResponse(w, row)
}

// 從URL arguments 取搜尋條件
func (env *Env) sqlGet(w http.ResponseWriter, r *http.Request) {
	log.Println("sqlGet")
	printRequest(r)

	if !checkJwtToken(env, w, r) {
		return
	}

	qName := r.URL.Query().Get(db.P_NAME)
	qAge := r.URL.Query().Get(db.P_AGE)
	qSex := r.URL.Query().Get(db.P_SEX)
	qHeight := r.URL.Query().Get(db.P_HEIGHT)

	var age *int64 = nil
	var sex *bool = nil
	var height *decimal.Decimal = nil

	if qAge != "" {
		a, err := util.StringToInt64(qAge)
		if err == nil {
			age = &a
		}
	}

	if qSex != "" {
		s, err := util.StringToBool(qSex)
		if err == nil {
			sex = &s
		}
	}

	if qHeight != "" {
		h, err := util.StringToDecimal(qHeight)
		if err == nil {
			height = &h
		}
	}

	if data, err := db.Get(env.db, qName, age, sex, height); err != nil {
		writeNormalFail(w, "Get Fail")
	} else {
		writeSuccessResponse(w, data)
	}
}

func printRequest(r *http.Request) {
	r.ParseForm() //解析參數，預設是不會解析的
	// fmt.Println("------NewRequest-------")
	// fmt.Println("r.Form=", r.Form) //這些資訊是輸出到伺服器端的列印資訊
	// fmt.Println("path=", r.URL.Path)
	// fmt.Println("scheme=", r.URL.Scheme)
	// fmt.Println("GetQueryArgument[url_long]=", r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Print("key:", k)
	// 	fmt.Println(" val:", strings.Join(v, ""))
	// }
}

func (env *Env) stop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !checkJwtToken(env, w, r) {
		return
	}
	fmt.Println("Receive stop request")
	StopChannel <- struct{}{}
	writeSuccessResponse(w, "Server will close in 5 seconds......")
}
