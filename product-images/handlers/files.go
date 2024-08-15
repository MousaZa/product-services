package handlers

import (
	"github.com/hashicorp/go-hclog"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/MousaZa/product-services/product-images/files"
	"github.com/gorilla/mux"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// ServeHTTP implements the http.Handler interface
func (f *Files) UploadRest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// check that the filepath is a valid name and file
	if id == "" || fn == "" {
		f.invalidURI(r.URL.String(), rw)
		return
	}

	f.saveFile(id, fn, rw, r.Body)
}

func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad Request", "error", err)
		http.Error(rw, "Expected Multipart form data", http.StatusBadRequest)
		return
	}
	id, idErr := strconv.Atoi(r.FormValue("id"))

	f.log.Info("Process form for id ", id)

	if idErr != nil {
		f.log.Error("Bad Request", "error", err)
		http.Error(rw, "Expected integer id", http.StatusBadRequest)
		return
	}
	fl, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Bad Request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}
	f.saveFile(r.FormValue("id"), mh.Filename, rw, fl)

}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

// get the file from the store and returns it to the user in a gzipped format
func (f *Files) getFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Get file for product", "id", id, "path", path)

	/*
		fp := filepath.Join(id, path)
		fr, err := f.store.Get(fp)
		if err != nil {
			f.log.Error("Unable to get file", "file", fp, "error", err)
			http.Error(rw, "Unable to find file", http.StatusNotFound)
			return
		}
		defer fr.Close()

		// set the filetpe header
		// DetectContentType() function only uses the first 512 bytes
		d := make([]byte, 512)
		_, err = fr.Read(d)
		if err != nil {
			f.log.Error("Unable to read file headers", "file", fp, "error", err)
			http.Error(rw, "Unable to find file", http.StatusInternalServerError)
			return
		}

		// roll back the stream
		fr.Seek(0, 0)

		// detect the content type
		ct := http.DetectContentType(d)
		if ct != "" {
			// detected content type
			f.log.Debug("Detected content type", "type", ct, "file", fp)
			rw.Header().Add("Content-Type", ct)
		} else {
			// fall back to default
			f.log.Debug("Unable to detect content type", "file", fp)
			rw.Header().Add("Content-Type", "application/octet-stream")
		}

		// if client cant handle gzip send plain
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			f.log.Debug("Unable to handle gzipped", "file", fp)
			io.Copy(rw, fr)
			return
		}

		// client can handle gziped content send gzipped to speed up transfer
		// set the content encoding header for gzip
		rw.Header().Add("Content-Encoding", "gzip")

		// wrap the default writer in a gzip writer
		gw := gzip.NewWriter(rw)
		defer gw.Close()

		// write the file
		io.Copy(gw, fr)
	*/
}
