// Package example implements the whole example of the pacakges
// note: pie-rum-sdk examples are written by me while rest of the
// example is written by the flash 2.5 😅
package example

// Please Do Provide the Feedback and all recommended
// To improve the pie-rum model more

type Resp struct {
	// Stored []e.Object `json:"stored"`
	Info string `json:"info"`
}

type Req struct {
	ID      *string `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Query   *string `json:"query,omitempty"`
	Profile string  `json:"profile"`
}
