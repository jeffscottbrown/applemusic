package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jeffscottbrown/goapple/music"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSearchApple_WithWireMockContainer(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "wiremock/wiremock:3.13.0",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForListeningPort("8080/tcp").WithStartupTimeout(10 * time.Second),
	}
	wiremockC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer wiremockC.Terminate(ctx)

	host, err := wiremockC.Host(ctx)
	assert.NoError(t, err)
	port, err := wiremockC.MappedPort(ctx, "8080")
	assert.NoError(t, err)

	wiremockURL := "http://" + host + ":" + port.Port()

	music.AppleMusicAPI = wiremockURL + "/search"

	stub := map[string]any{
		"request": map[string]any{
			"method":  "GET",
			"urlPath": "/search",
			"queryParameters": map[string]any{
				"term":   map[string]string{"equalTo": "grateful"},
				"media":  map[string]string{"equalTo": "music"},
				"entity": map[string]string{"equalTo": "album"},
				"limit":  map[string]string{"equalTo": "5"},
			},
		},
		"response": map[string]any{
			"status": 200,
			"jsonBody": map[string]any{
				"results": []map[string]string{
					{
						"artistName":        "Grateful Dead",
						"collectionName":    "Wake Of The Flood",
						"collectionViewUrl": "https://en.wikipedia.org/wiki/Wake_of_the_Flood",
					},
					{
						"artistName":        "Grateful Dead",
						"collectionName":    "Shakedown Street",
						"collectionViewUrl": "https://en.wikipedia.org/wiki/Shakedown_Street",
					},
				},
			},
		},
	}

	body, _ := json.Marshal(stub)
	resp, err := http.Post(wiremockURL+"/__admin/mappings", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	result, errMsg := music.SearchApple("grateful", "5")
	assert.Empty(t, errMsg)
	assert.Len(t, result.Albums, 2)
	assert.Equal(t, "Grateful Dead", result.Albums[0].ArtistName)
	assert.Equal(t, "Wake Of The Flood", result.Albums[0].AlbumTitle)
	assert.Equal(t, "Grateful Dead", result.Albums[1].ArtistName)
	assert.Equal(t, "Shakedown Street", result.Albums[1].AlbumTitle)
}
