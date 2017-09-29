package progresstracker

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/kennygrant/sanitize"
)

type Bug struct {
	// The unique numeric id of this bug.
	Id string
	// The summary of this bug.
	Summary string
	// The current status of the bug.
	Status string
	// Bug url
	Url string
}

func BugInfo(url string, bugid string) Bug {

	var (
		summary, status, word string
	)

	regexSummary, _ := regexp.Compile(`(?i)<title>(.*?)<\/title>`)
	regexStatus, _ := regexp.Compile(`(?i)<td class="status">(.*?)</td>`)

	msgArray := strings.Split(url, " ")

	for _, word = range msgArray {
		if strings.Contains(word, "http") {
			url = word
			break
		}
		if !strings.Contains(word, "http") && strings.Contains(word, "www") {
			url = "http://" + word
			break
		}
	}

	resp, err := http.Get(url + bugid)

	if err != nil {
		return Bug{}
	}

	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Bug{}
	}

	body := string(rawBody)
	noNewLines := strings.Replace(body, "\n", "", -1)
	noCarriageReturns := strings.Replace(noNewLines, "\r", "", -1)
	notSoRawBody := noCarriageReturns

	titleMatch := regexSummary.FindStringSubmatch(notSoRawBody)
	if len(titleMatch) > 1 {
		summary = strings.TrimSpace(sanitize.HTML(titleMatch[1]))
	}

	statusMatch := regexStatus.FindStringSubmatch(notSoRawBody)
	if len(statusMatch) > 1 {
		status = strings.TrimSpace(sanitize.HTML(statusMatch[1]))
	}

	return Bug{bugid, summary, status, url + bugid}

}
