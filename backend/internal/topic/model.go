package topic

type Topic struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

func getTopicKey(id string) string {
	return "topic:" + id
}
