package log

type Fields map[string]interface{}

func (f Fields) ToMap() map[string]interface{} {
	return map[string]interface{}(f)
}
