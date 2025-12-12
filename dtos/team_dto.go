package dtos

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required"`
}

type TeamResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
