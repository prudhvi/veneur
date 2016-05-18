package veneur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Flush takes the slices of metrics, combines then and marshals them to json
// for posting to Datadog.
func Flush(postMetrics [][]DDMetric) {
	totalCount := 0
	var finalMetrics []DDMetric
	// TODO This seems very inefficient
	for _, metrics := range postMetrics {
		totalCount += len(metrics)
		finalMetrics = append(finalMetrics, metrics...)
	}
	// Check to see if we have anything to do
	if totalCount > 0 {
		// TODO Watch this error
		postJSON, _ := json.Marshal(map[string][]DDMetric{
			"series": finalMetrics,
		})

		resp, err := http.Post(fmt.Sprintf("%s/api/v1/series?api_key=%s", Config.APIHostname, Config.Key), "application/json", bytes.NewBuffer(postJSON))
		defer resp.Body.Close()
		if err != nil {
			Stats.Count("flush.error_total", int64(totalCount), nil, 1.0)
			log.WithError(err).Error("Error posting")
		} else {
			log.WithField("metrics", len(finalMetrics)).Info("Completed flush to Datadog")
		}
		if log.GetLevel() == log.DebugLevel {
			// TODO Watch this error
			body, _ := ioutil.ReadAll(resp.Body)
			log.WithFields(log.Fields{
				"json":     string(postJSON),
				"status":   resp.Status,
				"headers":  resp.Header,
				"response": body,
			}).Debug("POSTing JSON")
		}
	} else {
		log.Info("Nothing to flush, skipping.")
	}
}
