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

func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	_,err := Get(key)
	if errors.Is(err,ErrorNoSuchKey){
		http.Error(w,err.Error(),http.StatusNotFound)
	}

	Delete(key)

	_,err = Get(key)
	if !errors.Is(err,ErrorNoSuchKey){
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE key=%s\n",key)
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}",KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}",KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}",KeyValueDeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000",r))
}