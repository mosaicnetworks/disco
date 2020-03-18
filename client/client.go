package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mosaicnetworks/disco/group"
	"github.com/sirupsen/logrus"
)

// DiscoClient is a client for the Discovery API
type DiscoClient struct {
	url      string
	certFile string
	client   *http.Client
	logger   *logrus.Entry
}

// NewDiscoClient creates a new DiscoClient for a server hosted at the provided
// url
func NewDiscoClient(url string, certFile string, logger *logrus.Entry) (*DiscoClient, error) {
	tlscfg := &tls.Config{}

	// If told to trust a certificate.
	if certFile != "" {
		// Load PEM-encoded certificate to trust.
		certPEM, err := ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}

		// Create CertPool containing the certificate to trust.
		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM(certPEM) {
			return nil, errors.New("failed to import certificate to trust")
		}

		// Trust the certificate by putting it into the pool of root CAs.
		tlscfg.RootCAs = roots

		// Decode and parse the server cert to extract the subject info.
		block, _ := pem.Decode(certPEM)
		if block == nil {
			return nil, errors.New("failed to decode certificate to trust")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		logger.Debugf("Trusting certificate %s with CN: %s", certFile, cert.Subject.CommonName)

		// Set ServerName in TLS config to CN from trusted cert so that
		// certificate will validate if CN does not match DNS name.
		tlscfg.ServerName = cert.Subject.CommonName
	}

	res := &DiscoClient{
		url:      fmt.Sprintf("https://%s", url),
		certFile: certFile,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlscfg,
			},
		},
		logger: logger,
	}

	return res, nil
}

// GetAllGroups returs all groups in a map where the key is the ID of the group
// and the value is a pointer to the corresponding Group object.
func (c *DiscoClient) GetAllGroups() (map[string]*group.Group, error) {
	path := fmt.Sprintf("%s/groups", c.url)
	fmt.Println("path: ", path)

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var allGroups map[string]*group.Group
	err = json.Unmarshal(body, &allGroups)
	if err != nil {
		return nil, fmt.Errorf("Error parsing groups: %v", err)
	}

	return allGroups, nil
}

// GetGroupByID gets a single group by ID
func (c *DiscoClient) GetGroupByID(id string) (*group.Group, error) {
	path := fmt.Sprintf("%s/groups/%s", c.url, id)
	fmt.Println("path: ", path)

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var group *group.Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, fmt.Errorf("Error parsing group: %v", err)
	}

	return group, nil
}

// CreateGroup adds a group to the discovery server. The group's ID field should
// be empty as it will be set by the server.
func (c *DiscoClient) CreateGroup(group group.Group) (string, error) {
	path := fmt.Sprintf("%s/group", c.url)
	fmt.Println("path: ", path)

	jsonValue, err := json.Marshal(group)
	if err != nil {
		return "", fmt.Errorf("Error marshalling group: %v", err)
	}

	resp, err := http.Post(path, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var id string
	err = json.Unmarshal(body, &id)
	if err != nil {
		return "", fmt.Errorf("Error parsing id: %v", err)
	}

	return id, nil
}
