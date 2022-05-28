package jsonoutput

type PulumiJSONOutput struct {
	Steps         []PulumiJSONSteps       `json:"steps"`
	Diagnostics   []PulumiJSONDiagnostics `json:"diagnostics"`
	Duration      int64                   `json:"duration"`
	ChangeSummary PulumiJSONChangeSummary `json:"changeSummary"`
}

type PulumiJSONSteps struct {
	Op       string `json:"op"`
	Urn      string `json:"urn"`
	Provider string `json:"provider,omitempty"`

	DiffReasons []string `json:"diffReasons,omitempty"`
}
type PulumiJSONDiagnostics struct {
	Urn      string `json:"urn,omitempty"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}
type PulumiJSONChangeSummary struct {
	Create  int `json:"create"`
	Delete  int `json:"delete"`
	Replace int `json:"replace"`
	Same    int `json:"same"`
	Update  int `json:"update"`
}
