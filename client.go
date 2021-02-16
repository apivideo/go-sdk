package apivideosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client type handles communicating with the api.video API
type Client struct {
	BaseURL    *url.URL
	APIKey     string
	httpClient *http.Client
	chunkSize  int64
	Token      *Token

	Videos       VideosServiceI
	Livestreams  LivestreamsServiceI
	UploadTokens UploadTokensServiceI
	Captions     CaptionsServiceI
	Chapters     ChaptersServiceI
	Players      PlayersServiceI
	Statistics   StatisticsServiceI
	Account      AccountServiceI
}

// Token contains token for connecting to the api.video API
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpireTime   time.Time
}

// ErrorResponse contains an error from the api.video API
type ErrorResponse struct {
	Response *http.Response
	Type     string `json:"type"`
	Title    string `json:"title"`
	Name     string `json:"name"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"[%d]: %v %v\nType: %v\nTitle: %v\nName: %v",
		r.Response.StatusCode,
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Type,
		r.Title,
		r.Name,
	)
}

const (
	defaultBaseURL        = "https://ws.api.video/"
	defaultSandboxBaseURL = "https://sandbox.api.video/"
)

// NewClient returns a new api.video API client for production
func NewClient(apiKey string) *Client {
	return newClient(apiKey, defaultBaseURL)
}

// NewSandboxClient returns a new api.video API client for sandbox environment
func NewSandboxClient(apiKey string) *Client {
	return newClient(apiKey, defaultSandboxBaseURL)
}

func newClient(apiKey, envURL string) *Client {

	baseURL, _ := url.Parse(envURL)

	c := &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		httpClient: http.DefaultClient,
		chunkSize:  128 * 1024 * 1024,
	}

	c.Videos = &VideosService{client: c}
	c.Livestreams = &LivestreamsService{client: c}
	c.UploadTokens = &UploadTokensService{client: c}
	c.Captions = &CaptionsService{client: c}
	c.Chapters = &ChaptersService{client: c}
	c.Players = &PlayersService{client: c}
	c.Statistics = &StatisticsService{client: c}
	c.Account = &AccountService{client: c}

	return c
}

//ChunkSize cganges chunk size for video upload, by default its 128MB
func (c *Client) ChunkSize(size int64) {
	c.chunkSize = size
}

func (c *Client) prepareRequest(method, urlStr string, body interface{}) (*http.Request, error) {

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {

		switch v := body.(type) {
		case *bytes.Buffer:
			buf = v
		default:
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "api.video SDK (go; v:1.0)")

	req, err = c.auth(req)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) prepareRangeRequests(urlStr string, filePath string) ([]*http.Request, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	stat, _ := file.Stat()

	var bufSize int64
	if stat.Size() > c.chunkSize && c.chunkSize != 0 {
		bufSize = c.chunkSize
	} else {
		bufSize = stat.Size()
	}

	buf := make([]byte, bufSize)
	requests := []*http.Request{}
	startByte := 0
	for {
		bytesread, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", file.Name())
		if err != nil {
			return nil, err
		}
		part.Write(buf)

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req, err := c.prepareRequest(http.MethodPost, urlStr, body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		if stat.Size() > c.chunkSize && c.chunkSize != 0 {
			ranges := fmt.Sprintf("bytes %d-%d/%d", startByte, (startByte+bytesread)-1, stat.Size())
			req.Header.Set("Content-Range", ranges)
			startByte = startByte + bytesread
		}

		if err != nil {
			return nil, err
		}

		requests = append(requests, req)
	}
	return requests, nil
}

func (c *Client) prepareUploadRequest(urlStr string, filePath string, extraFields map[string]string) (*http.Request, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	if extraFields != nil {
		for key, val := range extraFields {
			err = writer.WriteField(key, val)
			if err != nil {
				return nil, err
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := c.prepareRequest(http.MethodPost, urlStr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (c *Client) auth(req *http.Request) (*http.Request, error) {

	if c.Token == nil || time.Now().After(c.Token.ExpireTime) {
		u, err := c.BaseURL.Parse("/auth/api-key")
		if err != nil {
			return nil, err
		}

		payload := map[string]string{"apiKey": c.APIKey}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(payload)
		resp, err := c.httpClient.Post(u.String(), "application/json", buf)

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		err = checkResponse(resp)
		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(resp.Body).Decode(&c.Token)
		if err != nil {
			return nil, err
		}

		c.Token.ExpireTime = time.Now().Add(time.Duration(c.Token.ExpiresIn) * time.Second)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	return req, nil
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)

	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Title = string(data)
		}
	}

	return errorResponse
}
