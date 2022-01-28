package responses

// AddMemberRes represents the result of a request to the addmember endpoint.
type AddMemberRes struct {
	TeamID          string `json:"teamId"`
	NumMembersAdded int    `json:"numMembersAdded"`
}
