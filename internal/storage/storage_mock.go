package storage

// MockStorage — мок для тестов (реализует URLStorage).
type MockStorage struct {
	SaveURLFunc   func(urlToSave string, alias string) (int64, error)
	GetURLFunc    func(alias string) (string, error)
	DeleteURLFunc func(alias string) error
}

func (m *MockStorage) SaveURL(urlToSave string, alias string) (int64, error) {
	if m.SaveURLFunc != nil {
		return m.SaveURLFunc(urlToSave, alias)
	}
	return 0, nil
}

func (m *MockStorage) GetURL(alias string) (string, error) {
	if m.GetURLFunc != nil {
		return m.GetURLFunc(alias)
	}
	return "", ErrURLNotFound
}

func (m *MockStorage) DeleteURL(alias string) error {
	if m.DeleteURLFunc != nil {
		return m.DeleteURLFunc(alias)
	}
	return nil
}
