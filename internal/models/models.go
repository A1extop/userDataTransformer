package models

type XMLUsers struct {
	Users []XMLUser `xml:"user"`
}
type AgeGroupType string

const (
	AgeGroupUnder25 AgeGroupType = "До 25"
	AgeGroup25To35  AgeGroupType = "От 25 до 35"
	AgeGroupOver35  AgeGroupType = "После 35"
)

type XMLUser struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Email string `xml:"email"`
	Age   int    `xml:"age"`
}

type JSONUser struct {
	ID       string       `json:"id"`
	FullName string       `json:"full_name"`
	Email    string       `json:"email"`
	AgeGroup AgeGroupType `json:"age_group"`
}
type RetryItem struct {
	Data     *JSONUser
	Attempts int
}

func determineAgeGroup(age int) AgeGroupType {
	switch {
	case age < 25:
		return AgeGroupUnder25
	case age <= 35:
		return AgeGroup25To35
	default:
		return AgeGroupOver35
	}
}
func (u *XMLUser) ToJSONUser() *JSONUser {
	return &JSONUser{
		ID:       u.ID,
		FullName: u.Name,
		Email:    u.Email,
		AgeGroup: determineAgeGroup(u.Age),
	}
}
