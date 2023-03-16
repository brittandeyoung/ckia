package common

import (
	"bytes"
	"encoding/json"
)

type Check struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Criteria            string `json:"criteria"`
	RecommendedAction   string `json:"recommendedAction"`
	AdditionalResources string `json:"additionalResources"`
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
