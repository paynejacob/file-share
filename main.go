package main

import (
	_ "embed"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

//go:embed web/upload.html
var UploadPage []byte

//go:embed web/file.html
var FilePage []byte

type FileShare struct {
	mu          sync.RWMutex
	files       map[string][]byte
	contentType map[string]string
	fileTTL     time.Duration
}

func (f *FileShare) writeFile(in io.Reader, size int64, contentType string) (key string, err error) {
	t := time.NewTimer(f.fileTTL)
	key = randstr.Hex(6)

	f.mu.Lock()
	defer f.mu.Unlock()

	f.files[key] = make([]byte, size)
	f.contentType[key] = contentType

	_, err = in.Read(f.files[key])

	go func() {
		<-t.C
		f.deleteFile(key)
	}()

	return
}

func (f *FileShare) readFile(key string, out io.Writer) (contentType string, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	contentType, _ = f.contentType[key]
	if contentType == "" {
		return
	}

	_, err = out.Write(f.files[key])

	return
}

func (f *FileShare) deleteFile(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.files, key)
	delete(f.contentType, key)
}

func (f *FileShare) Download(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	var contentType string
	var err error

	if _, ok := r.URL.Query()["download"]; !ok {
		_, _ = w.Write(FilePage)
		w.Header().Add("Content-Type", "text/html")
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+key)

	contentType, err = f.readFile(key, w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if contentType == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Add("Content-Type", contentType)
}

func (f *FileShare) Upload(w http.ResponseWriter, r *http.Request) {
	var key string

	if r.Method == http.MethodGet {
		_, _ = w.Write(UploadPage)
		w.Header().Add("Content-Type", "text/html")
		return
	}

	// parse the form data
	data, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	key, err = f.writeFile(data, header.Size, header.Header.Get("Content-Type"))

	http.Redirect(w, r, "/"+key+"/", http.StatusTemporaryRedirect)
}

func main() {
	r := mux.NewRouter()

	fileShare := &FileShare{
		mu:          sync.RWMutex{},
		files:       make(map[string][]byte, 0),
		contentType: make(map[string]string, 0),
		fileTTL:     15 * time.Minute,
	}

	r.Path("/{key}/").Methods("GET", "POST").HandlerFunc(fileShare.Download)
	r.Path("/").Methods("GET", "POST").HandlerFunc(fileShare.Upload)

	log.Fatal(http.ListenAndServe(":80", r))
}
