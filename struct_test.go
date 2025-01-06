package jobj

// KeyFactResponseV2 is the requested json response schema from the key fact extractor
type KeyFactResponseV2Ex struct {
	Facts []struct {
		// The fact extracted from the company overview, inclusive of necessary context
		Fact string `json:"fact" xml:"fact"`
		// The type of fact, only one of the following: business, product, management, regulation, other
		FactType string `json:"fact_type" xml:"fact_type"`
	} `xml:"Facts>Fact"`
}

type KeyFactV2 struct {
	// LLM response request
	Fact     string `json:"fact"`
	FactType string `json:"fact_type"`
}
