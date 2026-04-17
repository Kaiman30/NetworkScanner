package models

// Results структура для хранения результатов проверки
type Results struct {
	Passed []string `json:"passed"`
	Failed []string `json:"failed"`
}

// NewResults создает новый экземпляр Results
func NewResults() *Results {
	return &Results{
		Passed: make([]string, 0),
		Failed: make([]string, 0),
	}
}

// AddPassed добавляет пройденную проверку
func (r *Results) AddPassed(names ...string) {
	r.Passed = append(r.Passed, names...)
}

// AddFailed добавляет проваленную проверку
func (r *Results) AddFailed(names ...string) {
	r.Failed = append(r.Failed, names...)
}
