package mock

import (
	"context"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/liov/hoper/go/v2/protobuf/response"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
)

func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := NewMockUserServiceServer(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	r := rand.New(rand.NewSource(1))
	req := model.NewPopulatedActiveReq(r, false)
	m.EXPECT().
		Active(context.Background(), req).
		Return(response.NewPopulatedBytesReply(r, false), nil).
		AnyTimes()
	res, _ := m.Active(context.Background(), req)
	log.Info(res)
}
