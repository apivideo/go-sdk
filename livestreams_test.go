package apivideosdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var liveJSONResponses = []string{
	`{
		"liveStreamId": "li2FgWk8CyBKFIGyDK1SimnL",
		"streamKey": "bd003ee9-9f71-44fd-a191-14d6a5063181",
		"name": "My livestream",
		"record": false,
		"broadcasting": false,
		"public": true,
		"assets": {
			"iframe": "<iframe src=\"https://embed.api.video/live/li2FgWk8CyBKFIGyDK1SimnL\" width=\"100%\" height=\"100%\" frameborder=\"0\" scrolling=\"no\" allowfullscreen=\"\"></iframe>",
			"player": "https://embed.api.video/live/li2FgWk8CyBKFIGyDK1SimnL",
			"hls": "https://live.api.video/li2FgWk8CyBKFIGyDK1SimnL.m3u8",
			"thumbnail": "https://cdn.api.video/live/li2FgWk8CyBKFIGyDK1SimnL/thumbnail.jpg"
		}
	}`,
	`{
		"liveStreamId": "liXyizQXIaWWenlD8pc3WLS",
		"streamKey": "ffee715c-72b8-4ed6-ab98-8d4602ced98f",
		"name": "My livestream 2",
		"record": false,
		"broadcasting": false,
		"public": false,
		"assets": {
			"iframe": "<iframe src=\"https://embed.api.video/live/liXyizQXIaWWenlD8pc3WLS\" width=\"100%\" height=\"100%\" frameborder=\"0\" scrolling=\"no\" allowfullscreen=\"\"></iframe>",
			"player": "https://embed.api.video/live/liXyizQXIaWWenlD8pc3WLS",
			"hls": "https://live.api.video/liXyizQXIaWWenlD8pc3WLS.m3u8",
			"thumbnail": "https://cdn.api.video/live/liXyizQXIaWWenlD8pc3WLS/thumbnail.jpg"
		}
	}`,
}
var liveStructs = []Livestream{
	Livestream{
		LivestreamID: "li2FgWk8CyBKFIGyDK1SimnL",
		Name:         "My livestream",
		StreamKey:    "bd003ee9-9f71-44fd-a191-14d6a5063181",
		Record:       false,
		Public:       true,
		Assets: &Assets{
			Hls:       "https://live.api.video/li2FgWk8CyBKFIGyDK1SimnL.m3u8",
			Iframe:    "<iframe src=\"https://embed.api.video/live/li2FgWk8CyBKFIGyDK1SimnL\" width=\"100%\" height=\"100%\" frameborder=\"0\" scrolling=\"no\" allowfullscreen=\"\"></iframe>",
			Thumbnail: "https://cdn.api.video/live/li2FgWk8CyBKFIGyDK1SimnL/thumbnail.jpg",
			Player:    "https://embed.api.video/live/li2FgWk8CyBKFIGyDK1SimnL",
		},
		Broadcasting: false,
	},
	Livestream{
		LivestreamID: "liXyizQXIaWWenlD8pc3WLS",
		Name:         "My livestream 2",
		StreamKey:    "ffee715c-72b8-4ed6-ab98-8d4602ced98f",
		Record:       false,
		Public:       false,
		Assets: &Assets{
			Hls:       "https://live.api.video/liXyizQXIaWWenlD8pc3WLS.m3u8",
			Iframe:    "<iframe src=\"https://embed.api.video/live/liXyizQXIaWWenlD8pc3WLS\" width=\"100%\" height=\"100%\" frameborder=\"0\" scrolling=\"no\" allowfullscreen=\"\"></iframe>",
			Thumbnail: "https://cdn.api.video/live/liXyizQXIaWWenlD8pc3WLS/thumbnail.jpg",
			Player:    "https://embed.api.video/live/liXyizQXIaWWenlD8pc3WLS",
		},
		Broadcasting: false,
	},
}

var liveRequestJSON = `{
	"name": "Test live",
	"record": true,
	"playerId": "pl4f4ferf5erfr5zed4fsdd"
  }`

var liveRequestStruct = LivestreamRequest{
	Name:     "Test live",
	Record:   true,
	PlayerID: "pl4f4ferf5erfr5zed4fsdd",
}

func TestLivestreams_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/live-streams/li2FgWk8CyBKFIGyDK1SimnL", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, liveJSONResponses[0])
	})

	live, err := client.Livestreams.Get("li2FgWk8CyBKFIGyDK1SimnL")
	if err != nil {
		t.Errorf("Livestreams.Get error: %v", err)
	}

	expected := &liveStructs[0]
	if !reflect.DeepEqual(live, expected) {
		t.Errorf("Livestreams.Get\n got=%#v\nwant=%#v", live, expected)
	}
}

func TestLivestreams_List(t *testing.T) {
	setup()
	defer teardown()
	JSONResp := fmt.Sprintf(
		`{"data":[%s,%s], "pagination":%s}`,
		liveJSONResponses[0],
		liveJSONResponses[1],
		paginationJSON)

	mux.HandleFunc("/live-streams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expectedQuery := url.Values{
			"currentPage": []string{"1"},
			"pageSize":    []string{"25"},
		}
		if !reflect.DeepEqual(r.URL.Query(), expectedQuery) {
			t.Errorf("Request querystring\n got=%#v\nwant=%#v", r.URL.Query(), expectedQuery)
		}
		fmt.Fprint(w, JSONResp)
	})

	opts := &LivestreamOpts{
		CurrentPage: 1,
		PageSize:    25,
	}
	livestreams, err := client.Livestreams.List(opts)
	if err != nil {
		t.Errorf("Livestreams.List error: %v", err)
	}

	expected := &LivestreamList{
		Data:       liveStructs,
		Pagination: &paginationStruct,
	}
	if !reflect.DeepEqual(livestreams, expected) {
		t.Errorf("Livestreams.List\n got=%#v\nwant=%#v", livestreams, expected)
	}
}

func TestLivestreams_Create(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/live-streams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		expectedBody := map[string]interface{}{
			"name":     "Test live",
			"record":   true,
			"playerId": "pl4f4ferf5erfr5zed4fsdd",
		}
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expectedBody) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expectedBody)
		}
		fmt.Fprint(w, liveJSONResponses[0])
	})

	livestream, err := client.Livestreams.Create(&liveRequestStruct)
	if err != nil {
		t.Errorf("Livestreams.Create error: %v", err)
	}

	expected := &liveStructs[0]
	if !reflect.DeepEqual(livestream, expected) {
		t.Errorf("Livestreams.Create\n got=%#v\nwant=%#v", livestream, expected)
	}
}

func TestLivestreams_Update(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/live-streams/li2FgWk8CyBKFIGyDK1SimnL", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		expectedBody := map[string]interface{}{
			"name":     "Test live",
			"record":   true,
			"playerId": "pl4f4ferf5erfr5zed4fsdd",
		}
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expectedBody) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expectedBody)
		}
		fmt.Fprint(w, liveJSONResponses[0])
	})

	livestream, err := client.Livestreams.Update("li2FgWk8CyBKFIGyDK1SimnL", &liveRequestStruct)
	if err != nil {
		t.Errorf("Livestreams.Update error: %v", err)
	}

	expected := &liveStructs[0]
	if !reflect.DeepEqual(livestream, expected) {
		t.Errorf("Livestreams.Update\n got=%#v\nwant=%#v", livestream, expected)
	}
}

func TestLivestreams_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/live-streams/li2FgWk8CyBKFIGyDK1SimnL", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	err := client.Livestreams.Delete("li2FgWk8CyBKFIGyDK1SimnL")
	if err != nil {
		t.Errorf("Livestreams.Delete error: %v", err)
	}
}

func TestLivestreams_UploadThumbnail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/live-streams/li2FgWk8CyBKFIGyDK1SimnL/thumbnail", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, liveJSONResponses[0])
	})

	file := createTempFile("test.thumbnail", 1024*1024)
	defer os.Remove(file)

	livestream, err := client.Livestreams.UploadThumbnail("li2FgWk8CyBKFIGyDK1SimnL", file)
	if err != nil {
		t.Errorf("Livestreams.UploadThumbnail error: %v", err)
	}

	expected := &liveStructs[0]
	if !reflect.DeepEqual(livestream, expected) {
		t.Errorf("Livestreams.UploadThumbnail\n got=%#v\nwant=%#v", livestream, expected)
	}
}

func TestLivestreams_DeleteThumbnail(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/live-streams/li2FgWk8CyBKFIGyDK1SimnL/thumbnail", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		fmt.Fprint(w, liveJSONResponses[0])
	})

	livestream, err := client.Livestreams.DeleteThumbnail("li2FgWk8CyBKFIGyDK1SimnL")
	if err != nil {
		t.Errorf("Livestreams.DeleteThumbnail error: %v", err)
	}

	expected := &liveStructs[0]
	if !reflect.DeepEqual(livestream, expected) {
		t.Errorf("Livestreams.DeleteThumbnail\n got=%#v\nwant=%#v", livestream, expected)
	}
}
