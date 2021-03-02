package updater

import (
	"crypto"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/inconshreveable/go-update"
)

const baseURL = "https://dl.paesa.es/krypton/"

var ErrNotAvailable = errors.New("no update available")

func Update(version string) error {
	resp, err := checkForUpdate(version)
	if err != nil {
		if errors.Is(err, ErrNotAvailable) {
			fmt.Println("Krypton is up to date.")
			return nil
		}
		return err
	}

	fmt.Printf("New version %s is available.\n", resp.Version)
	err = resp.Apply()
	if err != nil {
		return err
	}

	fmt.Printf("Updating Krypton to %s.\n", resp.Version)
	return nil
}

type Release struct {
	Version  string
	checksum []byte
}

type serverResponse struct {
	Version string `json:"version"`
	Sha256  string `json:"sha256"`
}

var (
	client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,

			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 15 * time.Second,
	}
)

func checkForUpdate(version string) (r Release, err error) {
	resp, err := client.Get(baseURL + runtime.GOOS + "-" + runtime.GOARCH + ".json")
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("server responded with %s to request %s", resp.Status, resp.Request.URL)
	}

	var serverResp serverResponse
	err = json.NewDecoder(resp.Body).Decode(&serverResp)
	if err != nil {
		return r, err
	}

	if serverResp.Version == version {
		return r, ErrNotAvailable
	}

	r.Version = serverResp.Version
	r.checksum, err = hex.DecodeString(serverResp.Sha256)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (r Release) Apply() error {
	opts := update.Options{
		Checksum: r.checksum,
		Hash:     crypto.SHA256,
	}

	resp, err := r.fetchUpdate()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, opts)
	if err != nil {
		return err
	}
	return nil
}

func (r Release) fetchUpdate() (*http.Response, error) {
	resp, err := client.Get(baseURL + r.Version + "/" + runtime.GOOS + "_" + runtime.GOARCH)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("server responded with %s to request %s", resp.Status, resp.Request.URL)
	}
	return resp, nil
}
