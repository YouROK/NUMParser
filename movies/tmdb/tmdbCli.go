package tmdb

import (
	"NUMParser/client"
	"errors"
	"net/url"
	"time"
)

func readPageTmdb(path string, params map[string]string, results interface{}) error {

	link := tmdbEndpoint + path
	querys := url.Values{}
	for key, value := range params {
		querys.Set(key, value)
	}

	link += "?" + querys.Encode()

	retryCodes := []int{
		429,
		500, 501, 502, 503, 504,
	}

	resp, _, errs := client.Get(link).
		AppendHeader("accept", "application/json").
		AppendHeader("Authorization", TMDBAuthKey).
		Retry(5, time.Second*20, retryCodes...).
		Timeout(time.Second * 120).
		EndStruct(results)

	if len(errs) > 0 {
		return errs[0]
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}
