package cmd

type WebVtt struct {
	VttElements []*VTTElement `json:"vtt_elements"`
}

type VTTElement struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Position  string `json:"position"`
	Line      string `json:"line"`
	Text      string `json:"text"`
}

//Scanner scan and bind one block of vtt file.
func (wv *WebVtt) Scanner() *VTTElement {
	return nil
}

//AppendVttElement append VTTElement to WebVtt
func (wv *WebVtt) AppendVttElement(vtt *VTTElement) {
	wv.VttElements = append(wv.VttElements, vtt)
}

func (wv *WebVtt) NewVttElement() *VTTElement {
	return &VTTElement{}
}
