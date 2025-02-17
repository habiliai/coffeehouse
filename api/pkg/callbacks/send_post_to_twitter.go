package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dghubble/oauth1"
	habconfig "github.com/habiliai/habiliai/api/pkg/config"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func Tweet(ctx context.Context, conf habconfig.TwitterConfig, message string) error {
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
		return errors.Wrapf(err, "failed to marshal payload")
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", tweetUrl, bytes.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to post tweet")
	}
	defer res.Body.Close()

	responseBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response body")
	}
	if res.StatusCode != http.StatusCreated {
		return errors.Errorf("failed to post tweet, status code: %d, body: '%s'", res.StatusCode, string(responseBodyBytes))
	}

	var respBody ResponseBody
	if err := json.Unmarshal(responseBodyBytes, &respBody); err != nil {
		return errors.Wrapf(err, "failed to decode response body")
	}

	if len(respBody.Errors) > 0 {
		return errors.Errorf("failed to post tweet, errors: %v", respBody.Errors)
	}

	return nil
}

func SendPostToTwitter(s *service, ctx context.Context, args []byte, metadata Metadata) (any, error) {
	logger.Info("call post to twitter callback")
	request := struct {
		Message   string   `json:"message"`
		Tags      []string `json:"tags"`
		MediaUrls []string `json:"media_urls"`
	}{}

	if err := json.Unmarshal(args, &request); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal args")
	}

	if err := Tweet(ctx, s.config.Twitter, request.Message); err != nil {
		return nil, errors.Wrapf(err, "failed to tweet")
	}

	logger.Info("post sent to twitter")
	return nil, nil
}

func init() {
	dispatchFunctions["post_to_twitter"] = SendPostToTwitter
}
