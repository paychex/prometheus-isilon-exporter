package isiclient

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/tidwall/gjson"

	"github.com/prometheus/common/log"
)

// ISIClient is used to connect to an EMC Isilon Cluster
type ISIClient struct {
	UserName       string
	Password       string
	authToken      string
	ClusterAddress string
	ClusterName    string
	ISIVersion     string
	NumNodes       int64
	ErrorCount     float64
}

// CallIsiAPI uses the client auth to call against an API endpoint and returns the string response
func (c *ISIClient) CallIsiAPI(request string, retryAttempts int) string {
	client := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}}
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.UserName, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("\n - Error connecting to Isilon: %s", err)
	}
	defer resp.Body.Close()
	respText, _ := ioutil.ReadAll(resp.Body)
	s := string(respText)
	if resp.StatusCode == 200 {
		log.Debugln(s)
	} else {
		if retryAttempts >= 1 {
			log.Infof("Got unknown code: %v when accessing URL: %s\n Body text is: %s\n", resp.StatusCode, request, respText)
			// to do need to re-auth to get back into system
			log.Info("retrying command")
			// now lets recursively call ourselves and hopefully we get in again
			s = c.CallIsiAPI(request, retryAttempts-1)
			c.ErrorCount++
		} else {
			s = ""
		}
	}
	return s
}

// NewIsiClient returns an initialized Isilon Client.
func NewIsiClient(user string, pass string, target string) (*ISIClient, error) {

	log.Debugln("Init ISI Client")

	c := ISIClient{
		UserName:       user,
		Password:       pass,
		ClusterAddress: target,
	}

	reqStatusURL := "https://" + c.ClusterAddress + ":8080/platform/1/cluster/config"

	// make a quick call to the API and ensure that it works
	s := c.CallIsiAPI(reqStatusURL, 2)
	if s != "" {
		c.ClusterName = gjson.Get(s, "name").String()
		c.ISIVersion = gjson.Get(s, "onefs_version.release").String()
		c.NumNodes = gjson.Get(s, "devices.#").Int()

		return &c, nil
	}

	return nil, errors.New("Error creating connection")

}
