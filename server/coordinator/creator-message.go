package coordinator

type CreatorMessage struct {
	Action CreatorAction `json:"action"`
}

func (c *CreatorMessage) IsValid() bool {
	if !c.Action.IsValid() {
		return false
	}

	return true
}
