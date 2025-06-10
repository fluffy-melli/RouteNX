package cache

type Block struct {
	OriginIP  string `json:"origin_ip"`
	ForwardIP string `json:"forward_ip"`
	Host      string `json:"host"`
	Time      string `json:"time"`
}

type Error struct {
	Error string `json:"error"`
	Time  string `json:"time"`
}

type Logger struct {
	Block []Block `json:"block"`
	Error []Error `json:"error"`
}

func (l *Logger) AddBlockLog(block Block) {
	l.Block = append(l.Block, block)
}

func (l *Logger) AddErrorLog(err Error) {
	l.Error = append(l.Error, err)
}

func NewLogger() *Logger {
	return &Logger{
		Block: make([]Block, 0),
		Error: make([]Error, 0),
	}
}
