package examples

import (
	"github.com/mhpenta/jobj"
)

type HeadlinesResponse struct {
	jobj.Schema
}

func NewHeadlineResponse() *HeadlinesResponse {
	h := &HeadlinesResponse{}
	h.Name = "HeadlinesResponse"
	h.Description = "HeadlinesResponse is the requested json response schema from the press release headline extractor"
	h.Fields = []*jobj.Field{
		jobj.Text("headline").
			Desc("The exact headline from the press release (in proper case)").Required(),
		jobj.Text("headline_without_company_name").
			Desc("The headline from the press release modified to remove the company name (in proper case)").Required(),
		jobj.Float("confidence").
			Desc("Confidence in the headlines extracted").Required(),
	}
	return h
}
