package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dghubble/oauth1"
	habconfig "github.com/habiliai/alice/api/pkg/config"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

const (
	habiliBotId = "alice_habili"
)

type TweetResponse struct {
	Link       string `json:"link"`
	Content    string `json:"content"`
	UploadedAt string `json:"uploaded_at"`
}

func Tweet(ctx context.Context, conf habconfig.TwitterConfig, message string) (*TweetResponse, error) {
	type (
		Payload struct {
			Text string `json:"text"`
		}
		ResponseBody struct {
			Data struct {
				Id   string `json:"id"`
				Text string `json:"text"`
			} `json:"data"`
			Errors []struct {
				Detail string `json:"detail"`
				Title  string `json:"title"`
				Type   string `json:"type"`
				Status int    `json:"status"`
			} `json:"errors"`
		}
	)
	const (
		tweetUrl = "https://api.x.com/2/tweets"
	)

	config := oauth1.NewConfig(conf.ConsumerKey, conf.ConsumerSecret)
	token := oauth1.NewToken(conf.AccessToken, conf.AccessTokenSecret)
	httpClient := config.Client(ctx, token)

	payload, err := json.Marshal(&Payload{
		Text: message,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload")
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", tweetUrl, bytes.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to post tweet")
	}
	defer res.Body.Close()

	responseBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read response body")
	}
	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("failed to post tweet, status code: %d, body: '%s'", res.StatusCode, string(responseBodyBytes))
	}

	var respBody ResponseBody
	if err := json.Unmarshal(responseBodyBytes, &respBody); err != nil {
		return nil, errors.Wrapf(err, "failed to decode response body")
	}

	if len(respBody.Errors) > 0 {
		return nil, errors.Errorf("failed to post tweet, errors: %v", respBody.Errors)
	}

	tweetId := respBody.Data.Id

	return &TweetResponse{
		Link:       "https://twitter.com/" + habiliBotId + "/status/" + tweetId,
		Content:    respBody.Data.Text,
		UploadedAt: time.Now().Format(time.RFC1123),
	}, nil
}

func SendPostToTwitter(ctx *Context, args []byte) (any, error) {
	logger.Info("call post to twitter callback")
	request := struct {
		Text string `json:"text"`
	}{}

	if err := json.Unmarshal(args, &request); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal args")
	}

	resp, err := Tweet(ctx, ctx.config.Twitter, request.Text)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to tweet")
	}

	logger.Info("post sent to twitter")

	if err := ctx.UpdateMemory(map[string]any{
		"tweet": resp,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

func init() {
	dispatchFunctions["post_to_twitter"] = SendPostToTwitter
	dispatchFunctions["post_to_x"] = SendPostToTwitter
}
