package gentoobug

import (
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Bug struct {
	// The unique numeric id of this bug.
	Id string
	// The summary of this bug.
	Summary string
	// The login name of the user to whom the bug is assigned.
	AssignedTo string
	// The current status of the bug.
	Status string
	// Bug url
	Url string
}

func BugInfo(url string, bugid string) Bug {

	var (
		summary, status, assigned, word string
	)

	regexSummary, _ := regexp.Compile(`(?i)<span id="short_desc_nonedit_display">(.*?)<\/span>`)
	regexStatus, _ := regexp.Compile(`(?i)<span id="static_bug_status">(.*?)\s.*?span`)
	regexAssigned, _ := regexp.Compile(`(?i)vcard.*?<span class="fn">(.*?)<\/span>`)

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
	assignedMatch := regexAssigned.FindStringSubmatch(notSoRawBody)
	if len(assignedMatch) > 1 {
		assigned = strings.TrimSpace(sanitize.HTML(assignedMatch[1]))
	}

	return Bug{bugid, summary, assigned, status, url + bugid}

}
