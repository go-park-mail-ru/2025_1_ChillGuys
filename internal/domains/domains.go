package domains

type Key string

var (
	ReqIDKey  Key = "ReqId"
	Token     Key = "Token"
	UserIDKey Key = "UserID"
	LoggerKey Key = "Logger"
)

func (k Key) String() string {
	return string(k)
}
