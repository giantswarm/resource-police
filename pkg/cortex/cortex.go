// Package cortex provides methods to query data from our central
// Prometheus data store.
package cortex

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

const (
	query = `aggregation:giantswarm:cluster_release_version{pipeline="testing"}`
)

type Config struct {
	// Base URL of the prometheus API server.
	URL string
	// User name.
	UserName string
	// Password (API token).
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

// QueryClusters queries cortex for a list of workload clusters
// that existed in time series at a given point in time.
// Returns a sorted slice of strings with format "<installation>/<cluster_id>".
func (s Service) QueryClusters(t time.Time) ([]string, error) {
	clusters := []string{}

	v1api := v1.NewAPI(s.client)

	value, warnings, err := v1api.Query(context.Background(), query, t)
	if err != nil {
		return clusters, microerror.Mask(err)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings while querying: %#v\n", warnings)
	}

	if value == nil {
		return clusters, microerror.Mask(fmt.Errorf("query %s returned nil value for time %s", query, t))
	}

	vector, ok := value.(model.Vector)
	if ok {
		if vector.Len() == 0 {
			return clusters, microerror.Mask(fmt.Errorf("query %s returned an empty result for time %s", query, t))
		}

		for i := 0; i < vector.Len(); i++ {
			installation := ""
			clusterID := ""

			if val, ok := vector[i].Metric["installation"]; ok {
				installation = string(val)
			} else {
				return clusters, microerror.Mask(errors.New("could not find required label 'installation' in sample"))
			}

			if val, ok := vector[i].Metric["cluster_id"]; ok {
				clusterID = string(val)
			} else {
				return clusters, microerror.Mask(errors.New("could not find required label 'cluster_id' in sample"))
			}

			clusters = append(clusters, fmt.Sprintf("%s/%s", installation, clusterID))
		}
	} else {
		return clusters, microerror.Mask(fmt.Errorf("query %s did not return a vector for time %s", query, t))
	}

	sort.Strings(clusters)

	return clusters, nil
}
