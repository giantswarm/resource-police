package installation

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	gsclient "github.com/giantswarm/gsclientgen/v2/client"
	"github.com/giantswarm/gsclientgen/v2/client/auth_tokens"
	"github.com/giantswarm/gsclientgen/v2/client/clusters"
	"github.com/giantswarm/gsclientgen/v2/models"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hako/durafmt"
	"sigs.k8s.io/yaml"
)

const (
	DatetimeLayout = "2006-01-02T15:04:05Z"
	DateLayout     = "2006-01-02"
)

var AgeLimit = time.Hour * 8

type Cluster struct {
	ID               string
	Name             string
	CreateDate       time.Time
	Labels           map[string]string
	Age              time.Duration
	AgeString        string
	Creator          string
	KeepUntil        time.Time
	IsV5             bool
	InstallationName string
}

type Config struct {
	Logger micrologger.Logger

	// Installations configuration
	InstallationsConfigFile string
}

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"`
}

type Installation struct {
	Name        string      `json:"name"`
	APIEndpoint string      `json:"apiEndpoint"`
	Credentials Credentials `json:"credentials"`

	Client *gsclient.Gsclientgen

	Clusters []*Cluster
}

type TemplateData struct {
	Clusters []*Cluster
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

	for i, inst := range installations {
		u, uerr := url.Parse(inst.APIEndpoint)
		if uerr != nil {
			return nil, microerror.Maskf(invalidConfigError, "API endpoint URL %s could not be parsed", inst.APIEndpoint)
		}

		tlsConfig := &tls.Config{}
		transport := httptransport.New(u.Host, "", []string{u.Scheme})
		transport.Transport = &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: tlsConfig,
		}
		installations[i].Client = gsclient.New(transport, strfmt.Default)
	}

	return installations, nil
}

func getAuthorization(i Installation) (runtime.ClientAuthInfoWriter, error) {
	authHeader := "giantswarm " + i.Credentials.Token
	return httptransport.APIKeyAuth("Authorization", "header", authHeader), nil
}

func ListClusters(i Installation) ([]*Cluster, error) {
	authToken, err := authorize(i)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	i.Credentials.Token = authToken

	authWriter, err := getAuthorization(i)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	fmt.Printf("Listing clusters for installation %s\n", i.Name)

	params := clusters.NewGetClustersParams()
	response, err := i.Client.Clusters.GetClusters(params, authWriter)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var clustersToDelete []*Cluster
	for _, cluster := range response.Payload {
		if cluster.DeleteDate != nil {
			// We skip clusters that are already in deletion.
			continue
		}

		created, err := time.Parse(DatetimeLayout, cluster.CreateDate)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		age := time.Since(created)

		if age < AgeLimit {
			// We skip clusters that are younger than 3 hours.
			continue
		}

		isV5 := false
		if strings.Contains(cluster.Path, "/v5/") {
			isV5 = true
		}

		c := &Cluster{
			ID:               cluster.ID,
			Name:             cluster.Name,
			CreateDate:       created,
			IsV5:             isV5,
			Age:              age,
			AgeString:        durafmt.ParseShort(age).String(),
			Labels:           cluster.Labels,
			InstallationName: i.Name,
		}

		if val, ok := c.Labels["creator"]; ok {
			c.Creator = val
		}
		if val, ok := c.Labels["keep-until"]; ok {
			c.KeepUntil, err = time.Parse(DateLayout, val)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		// If required labels are set, we look at the keep-until value
		if c.Creator != "" {
			if time.Until(c.KeepUntil) > 0 {
				continue
			}
		}

		clustersToDelete = append(clustersToDelete, c)
	}

	return clustersToDelete, nil
}

func RenderReport(clusters []*Cluster) (string, error) {
	fmt.Println("Rendering report")

	t := template.Must(template.New("report-template").Parse(reportTemplate))

	myData := TemplateData{
		Clusters: clusters,
	}

	var renderedReport bytes.Buffer
	err := t.Execute(&renderedReport, myData)
	if err != nil {
		return "", microerror.Mask(err)
	}

	fmt.Println("Report has been rendered")

	return renderedReport.String(), nil
}

func authorize(i Installation) (string, error) {
	fmt.Printf("Authorizing for installation %s\n", i.Name)
	params := auth_tokens.NewCreateAuthTokenParams().WithBody(&models.V4CreateAuthTokenRequest{
		Email:          i.Credentials.User,
		PasswordBase64: base64.StdEncoding.EncodeToString([]byte(i.Credentials.Password)),
	})
	response, err := i.Client.AuthTokens.CreateAuthToken(params, nil)
	if err != nil {
		return "", microerror.Mask(err)
	}

	fmt.Printf("Successfully authorized for installation %s\n", i.Name)

	return response.Payload.AuthToken, nil
}
