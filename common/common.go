package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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
