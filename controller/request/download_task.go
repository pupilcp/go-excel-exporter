package request

type CreateTaskReq struct {
	UserId        int         `json:"user_id"`
	FileName      string      `json:"file_name"`
	RequestUrl    string      `json:"request_url"`
	RequestParams interface{} `json:"request_params"`
	RequestMethod string      `json:"request_method"`
}
