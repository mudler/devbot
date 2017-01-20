package spamdetect

import (
	"fmt"
	"github.com/dghubble/sling"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

const baseURL = "https://api.uclassify.com/v1/"

type UclassifyError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e UclassifyError) Error() string {
	return fmt.Sprintf("uclassify: %v %v", strconv.Itoa(e.StatusCode), e.Message)
}

type ClassifyRequest struct {
	Texts []string `json:"texts"`
}

type ClassifyResponse struct {
	TextCoverage   float64 `json:"textCoverage"`
	Classification []struct {
		ClassName string  `json:"className"`
		P         float64 `json:"p"`
	} `json:"classification"`
}

// IssueService provides methods for creating and reading issues.
type ClassifyService struct {
	sling *sling.Sling
}

// Client is a tiny Github client
type Client struct {
	ClassifyService *ClassifyService
	// other service endpoints...
}

// NewIssueService returns a new IssueService.
func NewClassifyService(httpClient *http.Client) *ClassifyService {
	return &ClassifyService{
		sling: sling.New().Client(httpClient).Base(baseURL),
	}
}

// NewClient returns a new Client
func NewClient(key string) *Client {
	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: key}
	token.TokenType = "Token"
	httpClient := config.Client(oauth2.NoContext, token)
	return &Client{
		ClassifyService: NewClassifyService(httpClient),
	}
}

// Create creates a new issue on the specified repository.
func (s *ClassifyService) Classify(username, classifier string, texts []string) ([]ClassifyResponse, *http.Response, error) {

	request := &ClassifyRequest{
		Texts: texts,
	}

	response := new([]ClassifyResponse)
	classifyError := new(UclassifyError)
	path := fmt.Sprintf("%s/%s/classify", username, classifier)

	resp, err := s.sling.New().Post(path).BodyJSON(request).Set("Content-Type", "application/json").Receive(response, classifyError)

	if err == nil {
		err = classifyError
	}
	return *response, resp, err
}
