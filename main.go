package main

import (
	"encoding/json"
	"flag"
	"fmt"
	errorLog "log"
	"net/http"
	"os"

	"github.com/athlonxpgzw/alert-kvm-adapter/pkg/kvm"
	"github.com/go-kit/log"
)

type handler struct {
	Logger log.Logger
}

func main() {
	address := flag.String("listen-address", ":6725", "The address to listen on for HTTP requests.")
	json := flag.Bool("json", true, "enable json logging")
	flag.Parse()

	lw := log.NewSyncWriter(os.Stdout)
	var logger log.Logger
	if *json {
		logger = log.NewJSONLogger(lw)
	} else {
		logger = log.NewLogfmtLogger(lw)
	}
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)

	http.Handle("/", &handler{
		Logger: logger,
	})
	if err := http.ListenAndServe(*address, nil); err != nil {
		errorLog.Fatalf("failed to start http server: %v", err)
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var alerts kvm.Data
	err := json.NewDecoder(r.Body).Decode(&alerts)
	if err != nil {
		errorLog.Printf("cannot parse content because of %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(alerts)

	err = logAlerts(alerts, h.Logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func logAlerts(alerts kvm.Data, logger log.Logger) error {
	err := logger.Log("Topic", alerts.Topic, "zone", alerts.Params.Zone, "level", alerts.Params.Level, "message", alerts.Params.Mesg,
		"node_id", alerts.Params.NodeID, "node_ip", alerts.Params.NodeIP, "duration", alerts.Params.Duration, "duration_raw",
		alerts.Params.DurationRaw, "synopsis", alerts.Params.Synopsis, "node_type", alerts.Params.NodeType, "alert_type", alerts.Params.AlertType,
		"suggestion", alerts.Params.Suggestion, "event_count", alerts.Params.EventCount, "aggs_event_id", alerts.Params.AggsEventID, "send_duration",
		alerts.Params.SendDuration, "send_duration_raw", alerts.Params.SendDurationRaw, "first_notify_time", alerts.Params.FirstNotify,
		"last_notify_time", alerts.Params.LastestNotify,
		"status_time", alerts.Params.StatusTime, "timestamp", alerts.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

/*
func logWith(values map[string]string, logger log.Logger) log.Logger {
	for k, v := range values {
		logger = log.With(logger, k, v)
	}
	return logger
}
*/
