package apivideosdk

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const livestreamsBasePath = "live-streams"

// LivestreamsServiceI is an interface representing the Videos
// endpoints of the api.video API
// See: https://docs.api.video/5.1/live
type LivestreamsServiceI interface {
	Get(livestreamID string) (*Livestream, error)
	List(opts *LivestreamOpts) (*LivestreamList, error)
	Create(createRequest *LivestreamRequest) (*Livestream, error)
	Update(livestreamID string, updateRequest *LivestreamRequest) (*Livestream, error)
	Delete(livestreamID string) error
	UploadThumbnail(livestreamID string, filePath string) (*Livestream, error)
	DeleteThumbnail(livestreamID string) (*Livestream, error)
}

// LivestreamsService communicating with the Livestream
// endpoints of the api.video API
type LivestreamsService struct {
	client *Client
}

// Livestream represents an api.video Livestream
type Livestream struct {
	LivestreamID string  `json:"liveStreamId,omitempty"`
	Name         string  `json:"name,omitempty"`
	StreamKey    string  `json:"streamKey,omitempty"`
	Record       bool    `json:"record,omitempty"`
	Assets       *Assets `json:"assets,omitempty"`
	PlayerID     string  `json:"playerId,omitempty"`
	Broadcasting bool    `json:"broadcasting,omitempty"`
	Public       bool    `json:"public,omitempty"`
}

//LivestreamList represents a list of livestream
type LivestreamList struct {
	Data       []Livestream `json:"data,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
}

// LivestreamRequest represents a request to create / update a Livestream
type LivestreamRequest struct {
	Name     string `json:"name,omitempty"`
	Record   bool   `json:"record"`
	PlayerID string `json:"playerId,omitempty"`
}

//LivestreamOpts represents a query string to search on livestreams
type LivestreamOpts struct {
	CurrentPage int    `url:"currentPage,omitempty"`
	PageSize    int    `url:"pageSize,omitempty"`
	StreamKey   string `url:"streamKey,omitempty"`
}

//Get returns a Livestream by id
func (s *LivestreamsService) Get(livestreamID string) (*Livestream, error) {

	err := checkLivestreamID(livestreamID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", livestreamsBasePath, livestreamID)

	req, err := s.client.prepareRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	l := new(Livestream)
	_, err = s.client.do(req, l)

	if err != nil {
		return nil, err
	}

	return l, nil
}

//List returns a LivestreamList containing all livestreams
func (s *LivestreamsService) List(opts *LivestreamOpts) (*LivestreamList, error) {

	v, err := query.Values(opts)

	if err != nil {
		return nil, err
	}
	qs := v.Encode()

	path := fmt.Sprintf("%s?%s", livestreamsBasePath, qs)

	req, err := s.client.prepareRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	ll := new(LivestreamList)
	_, err = s.client.do(req, ll)

	if err != nil {
		return nil, err
	}

	return ll, nil
}

//Create a livestream container and returns it
func (s *LivestreamsService) Create(createRequest *LivestreamRequest) (*Livestream, error) {

	req, err := s.client.prepareRequest(http.MethodPost, livestreamsBasePath, createRequest)
	if err != nil {
		return nil, err
	}

	l := new(Livestream)
	_, err = s.client.do(req, l)

	if err != nil {
		return nil, err
	}

	return l, nil
}

//Update a video container and returns it
func (s *LivestreamsService) Update(livestreamID string, updateRequest *LivestreamRequest) (*Livestream, error) {

	err := checkLivestreamID(livestreamID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", livestreamsBasePath, livestreamID)

	req, err := s.client.prepareRequest(http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, err
	}

	l := new(Livestream)
	_, err = s.client.do(req, l)

	if err != nil {
		return nil, err
	}

	return l, nil
}

//Delete a livestream container
func (s *LivestreamsService) Delete(livestreamID string) error {

	err := checkLivestreamID(livestreamID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", livestreamsBasePath, livestreamID)

	req, err := s.client.prepareRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)

	if err != nil {
		return err
	}

	return nil
}

//UploadThumbnail upload the thumbnail of a livestream
func (s *LivestreamsService) UploadThumbnail(livestreamID string, filePath string) (*Livestream, error) {

	err := checkLivestreamID(livestreamID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s/thumbnail", livestreamsBasePath, livestreamID)

	req, err := s.client.prepareUploadRequest(path, filePath, nil)

	if err != nil {
		return nil, err
	}

	l := new(Livestream)

	_, err = s.client.do(req, l)

	if err != nil {
		return nil, err
	}

	return l, nil
}

//DeleteThumbnail upload the thumbnail of a livestream
func (s *LivestreamsService) DeleteThumbnail(livestreamID string) (*Livestream, error) {

	err := checkLivestreamID(livestreamID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s/thumbnail", livestreamsBasePath, livestreamID)

	req, err := s.client.prepareRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	l := new(Livestream)
	_, err = s.client.do(req, l)

	if err != nil {
		return nil, err
	}
	return l, nil
}
