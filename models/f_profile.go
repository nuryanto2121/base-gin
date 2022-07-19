package models

type ProfileForm struct {
	Name    string `json:"name" firestore:"name" valid:"Required"`
	Email   string `json:"email" firestore:"email" valid:"Required;Email"`
	PhoneNo string `json:"phone_no" firestore:"phone_no"`
	Avatar  string `json:"avatar" firestore:"avatar"`
	About   string `json:"about" firestore:"about"`
}

type ContactPhone struct {
	PhoneNo string `json:"phone_no" cql:"phone_no"`
}

type ProfileResponse struct {
	Name         string `json:"name" firestore:"name" valid:"Required"`
	Email        string `json:"email" firestore:"email" valid:"Required;Email"`
	PhoneNo      string `json:"phone_no" firestore:"phone_no"`
	Avatar       string `json:"avatar" firestore:"avatar"`
	TotalContact int    `json:"total_contact"`
	About        string `json:"about" firestore:"about"`
}
