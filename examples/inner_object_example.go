package examples

import obj "github.com/mhpenta/jobj"

type TranscriptCorrectionsResponse struct {
	obj.Schema
}

// NewTranscriptCorrectionsResponse returns a schema for the QuartrTranscriptSpeakersResponse schema.
func NewTranscriptCorrectionsResponse() *TranscriptCorrectionsResponse {
	h := &TranscriptCorrectionsResponse{}
	h.Name = "TranscriptCorrectionsResponse"
	h.Description = "TranscriptCorrectionsResponse is the requested json response schema for our transcription corrector, which looks at transcript paragraphs for errors and suggests corrections."

	correctionTypes := obj.AnyOf("correction_type", []obj.ConstDescription{
		{
			Const:       "merge",
			Description: "merge two paragraphs into one with a single speaker",
		},
		{
			Const:       "speaker_correction",
			Description: "correct the speaker number",
		},
	}).
		Desc("Type of correction to be made").
		Required()

	correction := obj.Object("correction",
		[]*obj.Field{
			correctionTypes,
			obj.Int("new_speaker_number").Required().Desc("New speaker number to be assigned to the paragraph"),
			obj.Int("paragraph_for_new_speaker_number").Required().Desc("Paragraph number to be assigned to the new speaker number"),
		}).Required().Desc("Details of the correction to be made to the transcript")

	h.Fields = []*obj.Field{
		obj.Array("corrections", []*obj.Field{correction}).Required(),
	}

	return h
}
