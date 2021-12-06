package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	errorLog "log"
	"net/http"
	"os"
	"strconv"

	"github.com/athlonxpgzw/alert-kvm-adapter/pkg/firefly"
	"github.com/athlonxpgzw/alert-kvm-adapter/pkg/kvm"
	"github.com/athlonxpgzw/alert-kvm-adapter/pkg/metrics"
	"github.com/go-kit/log"
)

type handler struct {
	Logger log.Logger
}

var bindAddress *string = flag.String("bindaddress", ":6725", "The address to listen on for HTTP requests.")
var sendURL *string = flag.String("sendURL", "", "The address to send to.")
var alertKey *string = flag.String("alertKey", "", "Firefly alert API certificate")
var alertID *string = flag.String("alertId", "", "Firefly alert API ID")
var logFmt *bool = flag.Bool("json", true, "enable json logging")

func main() {
	flag.Parse()

	lw := log.NewSyncWriter(os.Stdout)
	var logger log.Logger
	if *logFmt {
		logger = log.NewJSONLogger(lw)
	} else {
		logger = log.NewLogfmtLogger(lw)
	}
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)

	http.Handle("/alert", &handler{
		Logger: logger,
	})

	http.Handle("/metrics", metrics.Handler())
	http.HandleFunc("/healthz", healthzHandler)

	if err := http.ListenAndServe(*bindAddress, nil); err != nil {
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
	}

	sendAlert := makeFFAlert(alerts)
	ffAlerts, err := json.Marshal(*sendAlert)
	if err != nil {
		errorLog.Printf("cannot encode content because of %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(ffAlerts)
	ffRequest, err := http.NewRequest("POST", *sendURL, reader)
	if err != nil {
		errorLog.Printf("Cannot make http request because of %s", err)
		return
	}

	client := &http.Client{}
	ffResponse, err := client.Do(ffRequest)
	if err != nil {
		errorLog.Printf("Cannot Do the post request because of %s", err)
		return
	}
	defer ffResponse.Body.Close()
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

func makeFFAlert(alerts kvm.Data) *firefly.Data {
	var ffAlert firefly.Data
	if alerts.Topic == "alert_resolved" {
		ffAlert.Status = "resolved"
	} else {
		ffAlert.Status = "firing"
		ffAlert.Name = alerts.Params.Synopsis
		ffAlert.Desc = alerts.Params.Mesg
		ffAlert.Level = alerts.Params.Level
	}
	ffAlert.ApplyType = "custom"
	ffAlert.Key = *alertKey
	ffAlert.Id = *alertID
	ffAlert.MsgId = strconv.FormatInt(alerts.Params.AggsEventID, 10)
	return &ffAlert
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK\n")
}
