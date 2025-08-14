package config

import "github.com/stretchr/testify/mock"

type MockConfigService struct {
	mock.Mock
}

func (mk *MockConfigService) SetCurrentUsername(username string) {
	mk.Called(username)
}

func (mk *MockConfigService) GetCurrentUsername() string {
	args := mk.Called()
	return args.String(0)
}

func (mk *MockConfigService) Save() error {
	return mk.Called().Error(0)
}
