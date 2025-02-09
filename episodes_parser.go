package go_gaku

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

type Episode struct {
	Title string
	HLS string
	Date string
}

const eps_url = "https://www.nhk.or.jp/radio-api/app/v1/web/ondemand/series?site_id=%s&corner_site_id=%s"

func GetEps(prg Program) (string, []Episode) {
	var eps []Episode
	
	res, err := http.Get(fmt.Sprintf(eps_url, prg.SiteID, prg.CornerID))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	content := string(body)

	title := gjson.Get(content, "title").String()
	result := gjson.Get(content, "episodes")
	result.ForEach(func(key, value gjson.Result) bool {
		date :=	(strings.Split(
			(strings.Split(
				value.Get("aa_contents_id").String(), ";"))[4], "T"))[0]
		eps = append(eps,
			Episode{
				value.Get("program_title").String(),
				value.Get("stream_url").String(),
				date,
			})
		return true
	})

	return title, eps
}

// Local Variables:
// tab-width: 4
// End:
