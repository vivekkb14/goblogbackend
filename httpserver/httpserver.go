package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vivekkb14/goblogbackend/common"
	"github.com/vivekkb14/goblogbackend/dbops"
)

func CreateHttpServer() {
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

func insertIntoDb(w http.ResponseWriter, r *http.Request) {
	var err error
	var out []byte
	var id int
	resp := new(common.Response)
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

	var req common.UserRequest
	if err = common.ReadAndParseInput(w, r, &req); err != nil {
		return
	}
	fmt.Println("Request is: ", req)
	id, _ = dbops.GlobalDatabase.InserIntoTable(req)
	// dispalyData(db)
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
	var response common.DispalyAllAuthor
	var err error
	var res interface{}
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
	res, err = dbops.GlobalDatabase.DispalyData()
	if err != nil {
		fmt.Println("Error getting user info ")
		return
	}
	response.Data = res.([]common.Author)
	fmt.Println("AUTHOR: ", res)
	return
}

func getSingleAuthorInfo(w http.ResponseWriter, idNumber string) {
	var response common.Response
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
	response.Data, err = dbops.GlobalDatabase.SelectAuthorInfo(idNo)
	if err != nil {
		fmt.Println("Error getting user info ")
		return
	}
	return
}
