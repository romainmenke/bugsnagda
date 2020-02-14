package bugsnagda

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type APIError struct {
	Code        int      `json:"code"`
	CodeMeaning string   `json:"-"`
	Errors      []string `json:"errors"`
}

func (x APIError) String() string {
	return fmt.Sprintf("code : %d - %s, messages : %s", x.Code, x.CodeMeaning, strings.Join(x.Errors, ", "))
}

func (x APIError) Error() string {
	return x.String()
}

func errorFromResponse(resp *http.Response) error {
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	apiErr := &APIError{}
	err := decoder.Decode(apiErr)
	if err != nil {
		log.Println("Failed to decode an API error, this message is not the actual API error. Please open an issue if you suspect a package error at https://github.com/romainmenke/bugsnagda")
		apiErr := &APIError{
			Code:        resp.StatusCode,
			CodeMeaning: resp.Status,
			Errors:      []string{},
		}

		if apiErr.CodeMeaning == "" {
			apiErr.CodeMeaning = http.StatusText(resp.StatusCode)
		}

		return apiErr
	}

	switch apiErr.Code {
	case 30000:
		apiErr.CodeMeaning = "API access restricted due to lapsed payment"
	case 31000:
		apiErr.CodeMeaning = "API access restricted due to expired trial or invalid subscription"
	case 32000:
		apiErr.CodeMeaning = "API access restricted for this user due to lack of available seats"
	case 60000:
		apiErr.CodeMeaning = "Results limited due to unusually high number of Errors"
	}

	return apiErr
}
