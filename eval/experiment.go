package main

type Experiment struct {
	NumBoxes        uint64    `json:"num_boxes`
	MLength         uint64    `json:"message_length"`
	Construction1MS [][]int64 `json:"construction1_ms"`
	ExpressMS       []int64   `json:"express_ms"`
}
