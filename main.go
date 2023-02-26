package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Author struct {
	Title   string
	Content string
	Author  string
	Id      int
}

type DataId struct {
	Id int `json:"id"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Author `json:"data"`
}

type DispalyAllAuthor struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []Author `json:"data"`
}

type UserRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type GetArticle struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	AuthorInfo Author `json:"data"`
}

type ApiResp struct {
	ResponseCode        int32  `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
}

func GetUnmarshallErrorString(unMarshalErr error) error {
	if ute, ok := unMarshalErr.(*json.UnmarshalTypeError); ok {
		return errors.New("Input " + ute.Value + " for field " + ute.Field + " is incorrect.")
	} else {
		return unMarshalErr
	}
}

func ReadAndParseInput(w http.ResponseWriter, r *http.Request, input interface{}) error {
	const MAX_REST_API_PAYLOAD int64 = 1073741824
	const RESPONSE_FAILED = 1
	body, err := io.ReadAll(io.LimitReader(r.Body, MAX_REST_API_PAYLOAD))
	if err != nil {
		fmt.Errorf("Error in Reading request Body %+v", err)
		return err
	}
	if err1 := r.Body.Close(); err1 != nil {
		fmt.Errorf("Error in Closing body %s\n", err1.Error())
		return err1
	}

	if err2 := json.Unmarshal(body, input); err2 != nil {

		err2 = GetUnmarshallErrorString(err2)
		fmt.Errorf("Unmarshalling Error. %+v ", err2)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		output := ApiResp{
			ResponseCode:        RESPONSE_FAILED,
			ResponseDescription: err2.Error(),
		}
		if err3 := json.NewEncoder(w).Encode(output); err3 != nil {
			fmt.Errorf("Json Encoding Error. %+v", err3)
		}
		return err2
	}
	return nil
}

func createHttpServer() {
	var server http.Server
	// var err error
	router := mux.NewRouter()

	server.Addr = "127.0.0.1:8080"
	server.Handler = router

	fmt.Println("### HTTPS Server listening on \n", server.Addr)

	router.HandleFunc("/articles", insertIntoDb).Methods(http.MethodPost)
	router.HandleFunc("/articles", printAuthorInfo).Methods(http.MethodGet)
	fmt.Println(server.ListenAndServe())
}

func printAuthorInfo(w http.ResponseWriter, r *http.Request) {
	idNumber := r.FormValue("id")
	if idNumber != "" {
		getSingleAuthorInfo(w, idNumber)
	} else {
		getAllAuthorInfo(w)
	}
}

func getAllAuthorInfo(w http.ResponseWriter) {
	var response DispalyAllAuthor
	var err error
	defer func() {
		if err != nil {
			response.Message = err.Error()
			response.Status = http.StatusInternalServerError
		} else {
			response.Status = http.StatusAccepted
			response.Message = "Success"
		}
		out, _ := json.Marshal(response)
		w.Write(out)
	}()
	response.Data, err = dispalyData(db)
	if err != nil {
		fmt.Println("Error getting user info ")
		return
	}
	fmt.Println("AUTHOR: ", response.Data)
	return
}

func getSingleAuthorInfo(w http.ResponseWriter, idNumber string) {
	var response Response
	var err error
	defer func() {
		if err != nil {
			response.Message = err.Error()
			response.Status = http.StatusInternalServerError
		} else {
			response.Status = http.StatusAccepted
			response.Message = "Success"
		}
		out, _ := json.Marshal(response)
		w.Write(out)
	}()

	idNo, err := strconv.Atoi(idNumber)
	if err != nil {
		fmt.Println("Id number should be valid integer")
		return
	}
	response.Data, err = selectAuthorInfo(idNo)
	if err != nil {
		fmt.Println("Error getting user info ")
		return
	}
	return
}

func initialiseDatabaeServer() error {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "new_db_kbv",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr.Error())
		return pingErr
	}
	return nil
}

func insertIntoDb(w http.ResponseWriter, r *http.Request) {
	var err error
	var out []byte
	var id int
	resp := new(Response)
	defer func() {
		if err == nil {
			resp.Status = 201
			resp.Message = "Success"
			resp.Data.Id = id
		} else {
			resp.Status = 500
			resp.Message = err.Error()
			resp.Data.Id = 0
		}
		out, _ = json.Marshal(resp)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		_, err := w.Write(out)
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}()

	var req UserRequest
	if err = ReadAndParseInput(w, r, &req); err != nil {
		return
	}
	fmt.Println("Request is: ", req)
	id, _ = inserIntoTable(db, req)
	dispalyData(db)
}

func main() {
	err := initialiseDatabaeServer()
	if err != nil {
		fmt.Println("Error in initialising database")
		return
	}
	fmt.Println("Connected!")
	createProductTable(db)
	createHttpServer()
}

func createProductTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS articles(Id int primary key auto_increment, Title text, 
		Content text, Author text)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when getting rows affected", err)
		return err
	}
	fmt.Printf("Rows affected when creating table: %d\n", rows)
	return nil
}

func inserIntoTable(db *sql.DB, userInfo UserRequest) (int, error) {
	query := `INSERT INTO articles (Title, Content, Author) VALUES
  ('` + userInfo.Title + `', '` + userInfo.Content + `', '` + userInfo.Author + `');`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when getting rows affected", err)
		return 0, err
	}
	fmt.Printf("Rows affected when creating table: %d\n", rows)
	id := ExecuteQuery("")
	fmt.Println("Id is: ", id)
	return id, nil
}

func dispalyData(db *sql.DB) ([]Author, error) {
	var author Author
	var allauthors []Author
	query := `select * from articles;`
	res, err := db.Query(query)
	// var result Result
	// db.Raw("DESCRIBE TABLE_NAME").Scan(&result)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return allauthors, err
	}
	for res.Next() {

		err := res.Scan(&author.Id, &author.Title, &author.Content, &author.Author)

		if err != nil {
			fmt.Println(err)
			return allauthors, err
		}

		fmt.Printf("Author info %v\n", author)
		fmt.Printf("Author type %v\n", reflect.TypeOf(author))
		allauthors = append(allauthors, author)
	}

	return allauthors, nil
}

func selectAuthorInfo(id int) (Author, error) {
	fmt.Println("Inside selectAuthorInfo 1")
	var author Author
	query := `select * from articles where id=` + fmt.Sprintf("%d", id) + `;`
	res, err := db.Query(query)
	respData := res.Next()
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return author, err
	}
	if respData {
		err := res.Scan(&author.Id, &author.Title, &author.Content, &author.Author)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", author)
		if author.Id == 0 {
			return author, errors.New("Couldn't not find blog of the id number")
		}
	} else {
		fmt.Println("FALSE")
		return author, errors.New("Couldn't not find blog of the id number")
	}
	return author, nil
}

func ExecuteQuery(queryInfo string) int {
	var id int
	query := `SELECT Id FROM articles ORDER BY id DESC LIMIT 1;`
	res, err := db.Query(query)
	// var result Result
	// db.Raw("DESCRIBE TABLE_NAME").Scan(&result)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return 0
	}
	for res.Next() {
		err := res.Scan(&id)
		if err != nil {
			fmt.Println("Error getting id number\n", err)
			return 0
		}
	}
	return id
}
