package apivideosdk

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const videosBasePath = "videos"

// VideosServiceI is an interface representing the Videos
// endpoints of the api.video API
// See: https://docs.api.video/5.1/videos
type VideosServiceI interface {
	Get(videoID string) (*Video, error)
	List(opts *VideoOpts) (*VideoList, error)
	Create(createRequest *VideoRequest) (*Video, error)
	Update(videoID string, updateRequest *VideoRequest) (*Video, error)
	Delete(videoID string) error
	Upload(videoID string, filePath string) (*Video, error)
	Status(videoID string) (*VideoStatus, error)
	PickThumbnail(videoID string, timecode string) (*Video, error)
	UploadThumbnail(videoID string, filePath string) (*Video, error)
}

// VideosService communicating with the Videos
// endpoints of the api.video API
type VideosService struct {
	client *Client
}

// Video represents an api.video Video
type Video struct {
	VideoID     string     `json:"videoId,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	PublishedAt string     `json:"publishedAt,omitempty"`
	UpdatedAt   string     `json:"updatedAt,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Metadata    []Metadata `json:"metadata,omitempty"`
	Source      *Source    `json:"source,omitempty"`
	Assets      *Assets    `json:"assets,omitempty"`
	PlayerID    string     `json:"playerId,omitempty"`
	Public      bool       `json:"public,omitempty"`
	Panoramic   bool       `json:"panoramic,omitempty"`
	Mp4Support  bool       `json:"mp4Support,omitempty"`
}

// VideoRequest represents a request to create / update a Video
type VideoRequest struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Metadata    []Metadata `json:"metadata,omitempty"`
	Source      string     `json:"source,omitempty"`
	PlayerID    string     `json:"playerId,omitempty"`
	Public      bool       `json:"public,omitempty"`
	Panoramic   bool       `json:"panoramic,omitempty"`
	Mp4Support  bool       `json:"mp4Support,omitempty"`
}

//VideoStatus represents the encoding status of one video
type VideoStatus struct {
	Ingest   *Ingest   `json:"ingest,omitempty"`
	Encoding *Encoding `json:"encoding,omitempty"`
}

//Ingest represents the ingest status of one video
type Ingest struct {
	Status        string              `json:"status,omitempty"`
	Filesize      int                 `json:"filesize,omitempty"`
	ReceivedBytes []ReceivedBytesItem `json:"receivedBytes,omitempty"`
}

//ReceivedBytesItem represents a received bytes item
type ReceivedBytesItem struct {
	To    int `json:"to,omitempty"`
	From  int `json:"from,omitempty"`
	Total int `json:"total,omitempty"`
}

//Encoding represents the encoding status of one video
type Encoding struct {
	Playable  bool `json:"playable,omitempty"`
	Qualities []Quality
	Metadata  *EncodingMetadata
}

//Quality represents a quality
type Quality struct {
	Quality string `json:"quality,omitempty"`
	Status  string `json:"status,omitempty"`
}

//EncodingMetadata represents a encoding metadata
type EncodingMetadata struct {
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Bitrate     int    `json:"bitrate,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	Framerate   int    `json:"framerate,omitempty"`
	Samplerate  int    `json:"samplerate,omitempty"`
	VideoCodec  string `json:"videoCodec,omitempty"`
	AudioCodec  string `json:"audioCodec,omitempty"`
	AspectRatio string `json:"aspectRatio,omitempty"`
}

//VideoList represents a list of videos
type VideoList struct {
	Data       []Video     `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

//VideoOpts represents a query string to search on videos
type VideoOpts struct {
	CurrentPage  int               `url:"currentPage,omitempty"`
	PageSize     int               `url:"pageSize,omitempty"`
	SortBy       string            `url:"sortBy,omitempty"`
	SortOrder    string            `url:"sortOrder,omitempty"`
	Title        string            `url:"title,omitempty"`
	Tags         []string          `url:"tags,brackets,omitempty"`
	Description  string            `url:"description,omitempty"`
	LivestreamID string            `url:"livestreamId,omitempty"`
	Metadata     map[string]string `url:"-"`
}

//Metadata represents a Video metadata
type Metadata struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

//Source represents a Video source
type Source struct {
	URI        string           `json:"uri,omitempty"`
	Type       string           `json:"type,omitempty"`
	Livestream *VideoLivestream `json:"livestream,omitempty"`
}

//VideoLivestream represents a Video livestream
type VideoLivestream struct {
	LivestreamID string `json:"livestreamId,omitempty"`
	Links        []Link `json:"links,omitempty"`
}

//Assets represents Video assets
//Assets.Mp4 is available if Mp4Support is true and the video is encoded, else it's empty string
type Assets struct {
	Hls       string `json:"hls,omitempty"`
	Iframe    string `json:"iframe,omitempty"`
	Player    string `json:"player,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Mp4       string `json:"mp4,omitempty"`
}

//Get returns a Video by id
func (s *VideosService) Get(videoID string) (*Video, error) {

	err := checkVideoID(videoID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", videosBasePath, videoID)

	req, err := s.client.prepareRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	v := new(Video)
	_, err = s.client.do(req, v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

//List returns a VideoList containing all videos matching VideoOpts
func (s *VideosService) List(opts *VideoOpts) (*VideoList, error) {

	err := checkOpts(opts)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(opts)

	if err != nil {
		return nil, err
	}
	qs := v.Encode()

	if len(opts.Metadata) > 0 {
		u := url.Values{}
		for k, v := range opts.Metadata {
			u.Add(fmt.Sprintf("metadata[%s]", k), v)
		}
		qs = fmt.Sprintf("%s&%s", qs, u.Encode())
	}

	path := fmt.Sprintf("%s?%s", videosBasePath, qs)

	req, err := s.client.prepareRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	vl := new(VideoList)
	_, err = s.client.do(req, vl)

	if err != nil {
		return nil, err
	}

	return vl, nil
}

//Create a video container and returns it
func (s *VideosService) Create(createRequest *VideoRequest) (*Video, error) {

	req, err := s.client.prepareRequest(http.MethodPost, videosBasePath, createRequest)
	if err != nil {
		return nil, err
	}

	v := new(Video)
	_, err = s.client.do(req, v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

//Update a video container and returns it
func (s *VideosService) Update(videoID string, updateRequest *VideoRequest) (*Video, error) {

	err := checkVideoID(videoID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", videosBasePath, videoID)

	req, err := s.client.prepareRequest(http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, err
	}

	v := new(Video)
	_, err = s.client.do(req, v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

//Delete a video container
func (s *VideosService) Delete(videoID string) error {

	err := checkVideoID(videoID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", videosBasePath, videoID)

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

//Upload a video in a container.
//The upload is chuncked if the file size is more than 128MB
func (s *VideosService) Upload(videoID string, filePath string) (*Video, error) {

	path := fmt.Sprintf("%s/%s/source", videosBasePath, videoID)

	requests, err := s.client.prepareRangeRequests(path, filePath)

	if err != nil {
		return nil, err
	}

	v := new(Video)

	for _, req := range requests {
		_, err = s.client.do(req, v)

		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

//Status returns the  status of encoding and ingest of a video
func (s *VideosService) Status(videoID string) (*VideoStatus, error) {

	err := checkVideoID(videoID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s/status", videosBasePath, videoID)

	req, err := s.client.prepareRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	vs := new(VideoStatus)
	_, err = s.client.do(req, vs)

	if err != nil {
		return nil, err
	}

	return vs, nil
}

//PickThumbnail change the thumbnail of a video with a timecode
func (s *VideosService) PickThumbnail(videoID string, timecode string) (*Video, error) {

	err := checkVideoID(videoID)
	if err != nil {
		return nil, err
	}

	err = checkTimecode(timecode)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s/thumbnail", videosBasePath, videoID)

	body := map[string]string{
		"timecode": timecode,
	}

	req, err := s.client.prepareRequest(http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}

	v := new(Video)
	_, err = s.client.do(req, v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

//UploadThumbnail upload the thumbnail of a video
func (s *VideosService) UploadThumbnail(videoID string, filePath string) (*Video, error) {

	err := checkVideoID(videoID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s/thumbnail", videosBasePath, videoID)

	req, err := s.client.prepareUploadRequest(path, filePath, nil)

	if err != nil {
		return nil, err
	}

	v := new(Video)

	_, err = s.client.do(req, v)

	if err != nil {
		return nil, err
	}

	return v, nil
}
