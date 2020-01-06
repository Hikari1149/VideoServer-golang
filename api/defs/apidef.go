package defs



//request
type UserCredential struct{
	UserName string `json:"user_name"`
	Pwd string `json:"pwd"`
}
//Data model
type VideoInfo struct{
	Id string
	AuthorId int
	Name string
	DisplayCtime string
}

type Comment struct{
	Id string
	VideoId string
	Author string
	Content string
}

type SimpleSession struct{
	Username string
	TTL int64
}

//response
type SignedUp struct{
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}