package service

import (
	"fmt"
	"github.com/hopeio/utils/log"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"net/http"
)

func init() {
	// Create a new FileStore instance which is responsible for
	// storing the uploaded file on disk in the specified directory.
	// This path _must_ exist before tusd will store uploads in it.
	// If you want to save them on a different medium, for example
	// a remote FTP server, you can implement your own storage backend
	// by implementing the tusd.DataStore interface.

	store := NewFileStore("./data/uploads")

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	const prefix = "/api/v2/files/"
	// Create a new HTTP handler for the tusd server by providing a configuration.
	// The StoreComposer property must be set to allow the handler to function.
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              prefix,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		log.Fatalf("Unable to create handler: %s", err)
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go TusdHandler(handler)

	http.Handle(prefix, http.StripPrefix(prefix, handler))
}

func TusdHandler(handler *tusd.Handler) {
	for {
		event := <-handler.CompleteUploads
		fmt.Printf("Upload %s finished\n", event.Upload.ID)
	}
}
