package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"

	json "github.com/dutchcoders/elasticofs/json"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("elastico:client")

type Client struct {
	Client  *http.Client
	BaseURL *url.URL
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.Reader
	if body == nil {
	} else if v, ok := body.(io.Reader); ok {
		buf = v
	} else if v, ok := body.(json.M); ok {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf.(io.ReadWriter)).Encode(v); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("not supported type: %s", reflect.TypeOf(body))
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/json; charset=UTF-8")
	req.Header.Add("Accept", "text/json")
	return req, nil
}

func New(u string) (*Client, error) {
	baseURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:  http.DefaultClient,
		BaseURL: baseURL,
	}, nil
}

type elasticsearchResponse struct {
}

/*
{"acknowledged":true}map[string]interface {}{
    "acknowledged": bool(true),
}⏎
*/
/*
   "error": json.M{
       "root_cause": []interface {}{
           json.M{
               "type":   "illegal_argument_exception",
               "reason": "Malformed action/metadata line [1], expected START_OBJECT or END_OBJECT but found [null]",
           },
       },
       "type":   "illegal_argument_exception",
       "reason": "Malformed action/metadata line [1], expected START_OBJECT or END_OBJECT but found [null]",
   },


   {
  "acknowledged": true
}⏎

*/

type Response struct {
	Error  *Error `json:"error"`
	Status int64  `json:"status"`
}

type Error struct {
	FailedShards []struct {
		Index  string `json:"index"`
		Node   string `json:"node"`
		Reason struct {
			Reason string `json:"reason"`
			Type   string `json:"type"`
		} `json:"reason"`
		Shard int64 `json:"shard"`
	} `json:"failed_shards"`
	Grouped   bool   `json:"grouped"`
	Phase     string `json:"phase"`
	Reason    string `json:"reason"`
	RootCause []struct {
		Reason string `json:"reason"`
		Type   string `json:"type"`
	} `json:"root_cause"`
	Type string `json:"type"`
}

func (e Error) Error() string {
	return e.RootCause[0].Reason
}

func (wd *Client) do(req *http.Request, v interface{}) error {
	if b, err := httputil.DumpRequest(req, true); err == nil {
		log.Debug(string(b))
	}

	resp, err := wd.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if b, err := httputil.DumpResponse(resp, true); err == nil {
		log.Debug(string(b))
	}

	var r io.Reader = resp.Body

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < 300 {
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Not found")
	} else {
		resp := Response{}
		err = json.NewDecoder(r).Decode(&resp)
		return resp.Error
	}

	switch v := v.(type) {
	case io.Writer:
		io.Copy(v, resp.Body)
	case interface{}:
		return json.NewDecoder(r).Decode(&v)
	}

	return nil
}

func (wd *Client) Do(req *http.Request, v interface{}) error {
	return wd.do(req, v)
}
