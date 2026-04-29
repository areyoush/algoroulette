package model

type Question struct {
	ID			int		`json:"id"`
	Title		string	`json:"title"`
	Topic		string	`json:"topic"`
	Difficulty	string	`json:"difficulty"`
	Slug		string	`json:"slug"`
	Description	*string	`json:"description"`
	IsGlobal	string	`json:"is_global"`
	OwnerID		*int	`json:"owner_id,omitempty"`
	
	Status     *string `json:"status"`
	Bookmarked bool    `json:"bookmarked"`
	Notes      *string `json:"notes"`
}