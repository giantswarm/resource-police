package installation

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/ghodss/yaml"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/resource-police/internal/key"
)

const (
	AGE_LIMIT   = 8
	DATE_LAYOUT = "2006-01-02T15:04:05"
)

type AuthResponse struct {
	AuthToken string `json:"auth_token"`
}

type Cluster struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreateDate string `json:"create_date"`

	// internal field
	Age string
}

type Config struct {
	Logger micrologger.Logger

	// Installations configuration
	InstallationsConfigFile string
}

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Installation struct {
	Name        string      `json:"name"`
	APIEndpoint string      `json:"apiEndpoint"`
	Credentials Credentials `json:"credentials"`

	// internal field
	Clusters []Cluster
}

func New(config Config) (installations []Installation, err error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.InstallationsConfigFile == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.InstallationsConfigFile must not be empty", config)
	}

	b, err := ioutil.ReadFile(config.InstallationsConfigFile)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = yaml.Unmarshal(b, &installations)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return installations, nil
}

func ListClusters(i Installation) ([]Cluster, error) {
	authEndpoint := fmt.Sprintf("%s/%s", i.APIEndpoint, key.GSAPIAuthEndpoint)
	authToken, err := authorize(authEndpoint, i.Credentials)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	clustersEndpoint := fmt.Sprintf("%s/%s", i.APIEndpoint, key.GSAPIListClustersEndpoint)
	req, err := http.NewRequest("GET", clustersEndpoint, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("giantswarm %s", authToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var clusters []Cluster
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var clustersWithAge []Cluster
	for _, cluster := range clusters {
		d := strings.Split(cluster.CreateDate, "Z")
		createdAt, err := time.Parse(DATE_LAYOUT, d[0])
		if err != nil {
			return nil, microerror.Mask(err)
		}
		now := time.Now()
		timeDiff := now.Sub(createdAt)
		if timeDiff.Hours() > AGE_LIMIT {
			cluster.Age = fmt.Sprintf("%.0fh", timeDiff.Hours())
			clustersWithAge = append(clustersWithAge, cluster)
		}
	}

	return clustersWithAge, nil
}

func RenderReport(installations []Installation) (string, error) {
	fmt.Println("Rendering report...")

	t := template.Must(template.New("report-template").Parse(ReportTemplate))

	var renderedReport bytes.Buffer
	err := t.Execute(&renderedReport, installations)
	if err != nil {
		return "", microerror.Mask(err)
	}

	fmt.Println("Report has been rendered.")

	return renderedReport.String(), nil
}

func authorize(endpoint string, credentials Credentials) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"email":           credentials.User,
		"password_base64": base64.URLEncoding.EncodeToString([]byte(credentials.Password)),
	})
	if err != nil {
		return "", microerror.Mask(err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", microerror.Mask(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", microerror.Mask(err)
	}

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return authResponse.AuthToken, nil
}
