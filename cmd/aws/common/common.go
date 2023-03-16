package common

type Check struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Criteria            string `json:"criteria"`
	RecommendedAction   string `json:"recommendedAction"`
	AdditionalResources string `json:"additionalResources"`
}
