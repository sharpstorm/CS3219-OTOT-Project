package auth

import (
	"errors"
	"testing"

	"backend.cs3219.comp.nus.edu.sg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIsValidToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	adapter := mocks.NewMockDatabaseApiTokenAdapter(mockCtrl)
	authenticator := &tokenAuthenticator{
		tokenAdapter: adapter,
	}

	gomock.InOrder(
		adapter.EXPECT().IsValidToken("AAA").Return(false, nil),
		adapter.EXPECT().IsValidToken("AAA").Return(true, nil),
		adapter.EXPECT().IsValidToken("AAA").Return(true, errors.New("test error")),
	)

	assert.False(t, authenticator.IsValidToken("AAA"))
	assert.True(t, authenticator.IsValidToken("AAA"))
	assert.False(t, authenticator.IsValidToken("AAA"))
}
