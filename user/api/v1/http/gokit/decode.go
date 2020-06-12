package gokit

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeGetByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return getByUserIDRequest{UserID: id}, nil
}

func DecodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	cursor := r.URL.Query().Get("cursor")
	num := r.URL.Query().Get("num")
	n, _ := strconv.Atoi(num)
	return getAllRequest{
		Cursor: cursor,
		Number: n,
	}, nil
}

func DecodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request registerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request updateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return deleteRequest{ID: id}, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
