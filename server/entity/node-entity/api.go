package entity

//API api
type API struct {
	APIID         int
	APIName       string
	RequestURL    string
	RequestMethod string
	Protocol      string
	BalanceName   string
	IsFollow      bool
	StripPrefix   bool
	Timeout       int
	RetryCount    int
	TargetMethod  string
	TargetURL     string
	AlertValve    int
	StripSlash    bool // 是否过滤斜杠
}

//APIExtend api extend
type APIExtend struct {
	*API
	Target string
}
