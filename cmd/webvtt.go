package cmd

type WebVtt struct {
	VttElements []*VTTElement `json:"vtt_elements"`
}

type VTTElement struct {
	StartTime string `json:"start_time"`
}

//AppendVttElement append VTTElement to WebVtt
func (wv *WebVtt) AppendVttElement(vtt *VTTElement) {
	wv.VttElements = append(wv.VttElements, vtt)
}
