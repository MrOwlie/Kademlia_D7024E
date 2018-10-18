package serverd

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

type apiServer struct {
	recivingChannel chan []string
	sendingChannel  chan []string
}

func NewAPIServer(recivingChannel chan []string, sendingChannel chan []string) *apiServer {
	return &apiServer{recivingChannel: recivingChannel, sendingChannel: sendingChannel}
}

//Låta callern specifiera vilken port servern ska köra på?
func (server *apiServer) ListenApiServer( /*int serverPort*/ ) {
	mux := http.NewServeMux()
	mux.HandleFunc("/pin", server.pinFile)
	mux.HandleFunc("/unpin", server.unpinFile)
	mux.HandleFunc("/fetch", server.fetchFile)
	mux.HandleFunc("/upload", server.uploadFile)

	http.ListenAndServe(":80", mux)
}

func (server *apiServer) pinFile(response http.ResponseWriter, request *http.Request) {
	hash := request.FormValue("hash")
	response.Header().Set("Access-Control-Allow-Origin", "")
    response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH")
    response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if hash == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	} else {
		message := []string{"pin", hash}
		server.sendingChannel <- message
		results := <-server.recivingChannel
		response.Write([]byte(results[1]))
	}
}

func (server *apiServer) unpinFile(response http.ResponseWriter, request *http.Request) {
	hash := request.FormValue("hash")
	response.Header().Set("Access-Control-Allow-Origin", "")
    response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH")
    response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if hash == "" {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		message := []string{"unpin", hash}
		server.sendingChannel <- message
		results := <-server.recivingChannel
		response.Write([]byte(results[1]))
	}
}

func (server *apiServer) fetchFile(response http.ResponseWriter, request *http.Request) {
	hash := request.FormValue("hash")
	response.Header().Set("Access-Control-Allow-Origin", "")
    response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH")
    response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if hash == "" {
		//Om ingen fil specificeras kan man kanske skicka tillbaka hashes för alla filer man har?
		response.WriteHeader(http.StatusBadRequest)
	} else {
		message := []string{"cat", hash}
		server.sendingChannel <- message
		results := <-server.recivingChannel
		if results[0] == "fail" {
			response.WriteHeader(http.StatusNotFound)
		} else {
			file, err := os.Open(results[2])
			if err != nil {
				fmt.Println(err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			io.Copy(response, file)
		}
	}
}

func (server *apiServer) uploadFile(response http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(math.MaxInt64)
	response.Header().Set("Access-Control-Allow-Origin", "")
    response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH")
    response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
		newFile, err := os.Create(basePath + fileHeader[0].Filename)
		io.Copy(newFile, file)
		message := []string{"store", newFile.Name()}

		server.sendingChannel <- message
		results := <-server.recivingChannel
		response.Write([]byte(results[1]))		
	}	
}
