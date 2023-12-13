package main
import (
	"errors"
	"io"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func KeyValuePutHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	value,err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	err = Put(key,string(value))
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusCreated)

}

func KeyValueGetHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	value,err := Get(key)
	
	if errors.Is(err,ErrorNoSuchKey){
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}

	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))
	
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}",KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}",KeyValueGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000",r))
}