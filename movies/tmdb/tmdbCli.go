package tmdb

import (
	"github.com/jmcvetta/napping"
	"log"
	"time"
)

func readPageTmdb(path string, params napping.Params, results interface{}) error {
	var err error
	var resp *napping.Response
	for i := 0; i < 5; i++ {
		urlParams := params.AsUrlValues()
		resp, err = napping.Get(
			tmdbEndpoint+path,
			&urlParams,
			&results,
			nil,
		)

		if err != nil {
			log.Println(err)
		} else if resp.Status() == 304 {
			time.Sleep(time.Second * 5)
			continue
		} else if resp.Status() != 200 {
			log.Printf("Bad status %s: %d\n", path, resp.Status())
		}
		break
	}
	return err
}
