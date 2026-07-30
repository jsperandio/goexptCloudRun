package httpclient

// nao exporta a chave no log pq n vai no header
type NoopLogger struct{}

func (nl NoopLogger) Errorf(string, ...any) {}

func (nl NoopLogger) Warnf(string, ...any) {}

func (nl NoopLogger) Debugf(string, ...any) {}
