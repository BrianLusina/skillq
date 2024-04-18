package usersvc

import (
	"context"
	"errors"
	"testing"

	mockuserrepo "github.com/BrianLusina/skillq/server/app/internal/domain/mocks/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")
}

var _ = Describe("User Service", func() {
	t := GinkgoT()

	var (
		mockCtrl     *gomock.Controller
		mockUserRepo *mockuserrepo.MockUserRepoPort
		userSvc      userService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(t)
		mockUserRepo = mockuserrepo.NewMockUserRepoPort(mockCtrl)
		userSvc = userService{
			userRepo: mockUserRepo,
		}

		assert.NotNil(GinkgoT(), userSvc)
	})

	ctx := context.Background()

	Context("Creating a new user", func() {

		When("with invalid user input", func() {
			It("invalid email address", func() {
				defer mockCtrl.Finish()

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake",
					Password: "password",
					Skills:   []string{},
					Image:    []byte{},
					JobTitle: "The Boss",
				}

				actualUser, actualErr := userSvc.CreateUser(ctx, request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})
		})

		When("persisting user data", func() {
			It("should return error when repo fails to persist user data", func() {
				defer mockCtrl.Finish()

				mockRepoError := errors.New("failed to persist user information")

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    []byte{},
					JobTitle: "The Boss",
				}

				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, mockRepoError)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})
		})
	})
})
