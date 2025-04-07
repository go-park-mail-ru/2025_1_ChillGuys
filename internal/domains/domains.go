package domains

type ContextKey struct {
	name string
}

// Константы для ключей контекста
var (
	ReqIDKey  = ContextKey{name: "ReqId"}
	Token     = ContextKey{name: "Token"}
	UserIDKey = ContextKey{name: "UserID"}
	LoggerKey = ContextKey{name: "Logger"}
)

// Метод для приведения к строковому представлению
func (k ContextKey) String() string {
	return k.name
}
