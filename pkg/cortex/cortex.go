// Package cortex provides methods to query data from our central
// Prometheus data store.
package cortex

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

const (
	query = `aggregation:giantswarm:cluster_release_version{pipeline="testing"}`

	// Amount of time to look back. The longer the time frame, the slower
	// and more expensive the query.
	timeRange = 7 * 24 * time.Hour

	// If a cluster has been last seen more than this much time ago,
	// it is considered deleted. Be careful to make this at least as
	// large as the stepInterval.
	lastSeenDurationThreshold = time.Hour

	// Step interval to query. Should be large enough to avoid returning a huge
	// result with every query.
	stepInterval = time.Hour
)

type Config struct {
	// URL of the prometheus API server.
	URL string
	// UserName for authentication with the API server.
	UserName string
	// Password (API token) for authentication with the API server.
	Password string
}

type Service struct {
	config Config
	client api.Client
}

func New(conf Config) (*Service, error) {
	apiConfig := api.Config{
		Address: conf.URL,
		RoundTripper: config.NewBasicAuthRoundTripper(
			conf.UserName,
			config.Secret(conf.Password),
			"",
			api.DefaultRoundTripper),
	}

	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	s := &Service{
		config: conf,
		client: client,
	}

	return s, nil
}

// Cluster represents a workload cluster.
type Cluster struct {
	Installation   string
	ID             string
	Release        string
	Provider       string
	FirstTimestamp time.Time
}

// QueryClusters queries cortex for a list of workload clusters
// that existed in time series at a given point in time.
// Returns a slice of Clusters.
func (s Service) QueryClusters() ([]Cluster, error) {
	clusters := []Cluster{}

	v1api := v1.NewAPI(s.client)

	now := time.Now().UTC()
	queryRange := v1.Range{
		Start: now.Add(-timeRange),
		End:   now,
		Step:  stepInterval,
	}

	value, warnings, err := v1api.QueryRange(context.Background(), query, queryRange)
	if err != nil {
		return clusters, microerror.Mask(err)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings while querying: %#v\n", warnings)
	}

	if value == nil {
		return clusters, microerror.Maskf(executionFailedError, "query %s returned nil value", query)
	}

	matrix, ok := value.(model.Matrix)
	if ok {
		if matrix.Len() == 0 {
			log.Printf("Returned result is empty.\n")
			return clusters, microerror.Maskf(executionFailedError, "query %s returned an empty result", query)
		} else {
			log.Printf("Returned result is matrix with %d series.\n", matrix.Len())
		}

		for i := 0; i < matrix.Len(); i++ {
			installation := ""
			clusterID := ""
			release := ""
			provider := ""

			if val, ok := matrix[i].Metric["installation"]; ok {
				installation = string(val)
			} else {
				return clusters, microerror.Maskf(executionFailedError, "could not find required label 'installation' in sample")
			}

			if val, ok := matrix[i].Metric["cluster_id"]; ok {
				clusterID = string(val)
			} else {
				return clusters, microerror.Maskf(executionFailedError, "could not find required label 'cluster_id' in sample")
			}

			if val, ok := matrix[i].Metric["release_version"]; ok {
				release = string(val)
			} else {
				return clusters, microerror.Maskf(executionFailedError, "could not find required label 'release_version' in sample")
			}
			if val, ok := matrix[i].Metric["provider"]; ok {
				provider = string(val)
			} else {
				return clusters, microerror.Maskf(executionFailedError, "could not find required label 'provider' in sample")
			}

			first := int64(matrix[i].Values[0].Timestamp)
			latest := int64(matrix[i].Values[len(matrix[i].Values)-1].Timestamp)

			// Determine whether cluster still exists
			lastSeenDuration := now.Sub(time.Unix(latest/1000, 0))
			if lastSeenDuration > lastSeenDurationThreshold {
				// Skip as not existing any more.
				continue
			}

			c := Cluster{
				Installation:   installation,
				ID:             clusterID,
				Release:        release,
				Provider:       provider,
				FirstTimestamp: time.Unix(first/1000, 0),
			}
			clusters = append(clusters, c)
		}
	} else {
		return clusters, microerror.Maskf(executionFailedError, "query %s did not return a matrix", query)
	}

	return clusters, nil
}
