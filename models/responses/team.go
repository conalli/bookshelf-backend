package responses

// AddMemberResponse represents the result of a request to the addmember endpoint.
type AddMemberResponse struct {
	TeamID          string `json:"teamId"`
	NumMembersAdded int    `json:"numMembersAdded"`
}
