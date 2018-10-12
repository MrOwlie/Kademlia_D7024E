package serverd

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"math"
)

type apiServer struct {
	recivingChannel chan string
	sendingChannel chan string 
}

func NewAPIServer(recivingChannel chan string, sendingChannel chan string) *apiServer{
	return &apiServer{recivingChannel: recivingChannel, sendingChannel: sendingChannel}
}

//Låta callern specifiera vilken port servern ska köra på?
func (server *apiServer) ListenApiServer(/*int serverPort*/){
	mux := http.NewServeMux()
	mux.HandleFunc("/pin", server.pinFile)
	mux.HandleFunc("/unpin", server.unpinFile)
	mux.HandleFunc("/fetch", server.fetchFile)
	mux.HandleFunc("/upload", server.uploadFile)

	http.ListenAndServe(":80", mux)
}

func (server *apiServer) pinFile(response http.ResponseWriter, request *http.Request){
	hash := request.FormValue("hash")
	if hash == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	} else {
		server.sendingChannel<- "pin "+hash
		responseBody := <-server.recivingChannel
		response.Write([]byte(responseBody))
	}
}

func (server *apiServer) unpinFile(response http.ResponseWriter, request *http.Request){
	hash := request.FormValue("hash")
	if hash == "" {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		server.sendingChannel<- "unpin "+hash
		responseBody := <-server.recivingChannel
		response.Write([]byte(responseBody))
	}
}

func (server *apiServer) fetchFile(response http.ResponseWriter, request *http.Request){
	hash := request.FormValue("hash")
	if hash == "" {
		//Om ingen fil specificeras kan man kanske skicka tillbaka hashes för alla filer man har?
		response.WriteHeader(http.StatusBadRequest)
	} else {
		server.sendingChannel<- "fetch "+hash
		filePath := <-server.recivingChannel
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(response, file)
	}
}

func (server *apiServer) uploadFile(response http.ResponseWriter, request *http.Request){
	err := request.ParseMultipartForm(math.MaxInt64)
	if err != nil{
		fmt.Println(err)
		response.WriteHeader(http.StatusBadRequest)
	} else {
		basePath, err := filepath.Abs("../downloads/")
		fileHeader := request.MultipartForm.File["file"]
		file, err := fileHeader[0].Open()
		if err != nil {
			fmt.Println(err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		newFile, err := os.Create(basePath+fileHeader[0].Filename)
		io.Copy(newFile, file)

		server.sendingChannel<- "store "+newFile.Name()
		hash := <-server.recivingChannel
		response.Write([]byte(hash))
	}	
}