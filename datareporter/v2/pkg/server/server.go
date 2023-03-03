package server

import (
	"encoding/json"
	"io"
	golangLog "log"
	"net/http"
	// "github.com/gorilla/mux"

	"github.com/redhat-marketplace/redhat-marketplace-operator/datareporter/v2/pkg/events"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("meteric_generator")

// type DataReporterHandler struct {
// 	WriteTimeout time.Duration
// 	ReadTimeout  time.Duration
// }

//	func NewDataReporterHandler() *DataReporterHandler {
//		return &DataReporterHandler{
//			WriteTimeout: time.Second * 30,
//			ReadTimeout:  time.Second * 30,
//		}
//	}

func NewDataReporterHandler(eventEngine *events.EventEngine) *http.ServeMux {
	router := http.NewServeMux()

	if eventEngine == nil {
		golangLog.Fatal("eventEngine is nil")
	}

	router.HandleFunc("/v1/event", func(w http.ResponseWriter, r *http.Request) {
		EventHandler(eventEngine, w, r)
	})
	router.HandleFunc("/v1/status", StatusHandler)
	return router
}

func EventHandler(eventEngine *events.EventEngine, w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("apiKey")

	eventKeyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !json.Valid(eventKeyBytes) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rawMessage := json.RawMessage(eventKeyBytes)
	event := events.Event{Key: events.Key(apiKey), RawMessage: rawMessage}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventEngine.EventChan <- event

	w.WriteHeader(http.StatusOK)
	out, _ := json.Marshal(event)
	w.Write(out)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
