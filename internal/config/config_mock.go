package config

import "github.com/stretchr/testify/mock"

type MockConfigService struct {
	mock.Mock
}

func (mk *MockConfigService) SetUser(username string) {
	mk.Called(username)
}

func (mk *MockConfigService) GetUser() string {
	args := mk.Called()
	return args.String(0)
}

func (mk *MockConfigService) Save() error {
	return mk.Called().Error(0)
}
