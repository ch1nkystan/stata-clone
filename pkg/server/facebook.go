package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

type FacebookEvent struct {
	EventName string                 `json:"event_name"`
	EventID   string                 `json:"event_id"`
	EventTime string                 `json:"event_time"`
	UserData  map[string]interface{} `json:"user_data"`
}

type FacebookEventsRequest struct {
	Data []FacebookEvent `json:"data"`
}

func sendFacebookEvent(pl *types.PixelLink) {
	url := "https://graph.facebook.com/v19.0/" + string(pl.FBPixelID) + "/events?access_token=" + pl.FBAccessMarker

	reqBody := FacebookEventsRequest{
		Data: []FacebookEvent{
			{
				EventName: "Lead",
				EventTime: string(time.Now().Unix()),
				UserData: map[string]interface{}{
					"fbp": pl.FBP,
					"fbc": pl.FBC,
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
