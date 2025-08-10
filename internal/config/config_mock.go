package config

import "github.com/stretchr/testify/mock"

type MockConfigService struct {
	mock.Mock
}

func (mk *MockConfigService) SetUser(username string) {
	mk.Called(username)
}

func (mk *MockConfigService) Save() error {
	return mk.Called().Error(0)
}