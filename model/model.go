package model

type (
	Course struct {
		Name        string       `json:"name"`
		Alias       string       `json:"alias"`
		Vote        int          `json:"vote"`
		Participant Participants `json:"participant"`
	}

	Participant struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

type (
	Courses      []Course
	Participants []Participant
)
