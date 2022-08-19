package schema

//user struct

const DATABASE string = "auth"
const COLLECTION string = "users"

type User struct { //Payload
	Email            string `bson:"email,omitempty"`
	Username         string `bson:"username,omitempty"`
	Password         string `bson:"password,omitempty"`
	InitialTimestamp string `bson:"intitial_timestamp,omitempty"`
	RecentTimestamp  string `bson:"recent_timestamp,omitempty"`
	Role             int    `bson:"role,omitempty"`
}

func (user User) Database() string {
	return DATABASE
}

func (user User) Collection() string {
	return COLLECTION
}

func (user User) VerifyPassword(candidePassword string) bool {
	return user.Password == candidePassword //Decrypt
}
