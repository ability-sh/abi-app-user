package srv

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Nick     string `json:"nick,omitempty"`
	Password string `json:"password,omitempty"`
	Ctime    int64  `json:"ctime"`
}

type UserCreateTask struct {
	Name     string `json:"name"`
	Nick     string `json:"nick,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserGetTask struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Nick        string `json:"nick,omitempty"`
	AutoCreated bool   `json:"autoCreated"`
}

type UserSetTask struct {
	Id       string  `json:"id"`
	Name     *string `json:"name"`
	Nick     *string `json:"nick,omitempty"`
	Password *string `json:"password,omitempty"`
}

type InfoSetTask struct {
	Id   string      `json:"id"`
	Key  string      `json:"key"`
	Info interface{} `json:"info"`
}

type InfoGetTask struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type LoginTask struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

type UserBatchGetTask struct {
	Ids string `json:"ids"`
}

type InfoBatchGetTask struct {
	Ids string `json:"ids"`
	Key string `json:"key"`
}
