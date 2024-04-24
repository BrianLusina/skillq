package usersvc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	mockuserrepo "github.com/BrianLusina/skillq/server/app/internal/domain/mocks/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	mockmessagepublisher "github.com/BrianLusina/skillq/server/infra/messaging/mocks"
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
		mockCtrl      *gomock.Controller
		mockUserRepo  *mockuserrepo.MockUserRepoPort
		mockPublisher *mockmessagepublisher.MockPublisher
		userSvc       userService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(t)
		mockUserRepo = mockuserrepo.NewMockUserRepoPort(mockCtrl)
		mockPublisher = mockmessagepublisher.NewMockPublisher(mockCtrl)
		userSvc = userService{userRepo: mockUserRepo, messagePublisher: mockPublisher}

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

				// no query to fetch user by UUID was triggered
				mockUserRepo.EXPECT().GetUserByUUID(ctx, gomock.Any()).Times(0)

				// no user verification was created
				mockUserRepo.EXPECT().CreateUserVerification(ctx, gomock.Any()).Times(0)

				// no message was published
				mockPublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Times(0)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})

			It("should return error when repo fails to query user by UUID when creating email verification", func() {
				defer mockCtrl.Finish()

				mockRepoError := errors.New("failed to fetch user by UUID")

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    []byte{},
					JobTitle: "The Boss",
				}

				createdUser, err := user.New(user.UserParams{
					EntityParams: entity.EntityParams{
						EntityIDParams: entity.EntityIDParams{
							UUID:  id.NewUUID(),
							KeyID: id.NewKeyID(),
							XID:   id.NewXid(),
						},
						EntityTimestampParams: entity.EntityTimestampParams{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Metadata: map[string]any{},
					},
					Name:      request.Name,
					Email:     request.Email,
					ImageData: request.Image,
					Skills:    request.Skills,
					JobTitle:  request.JobTitle,
					Password:  "hashedPassword",
				})

				assert.NoError(t, err)

				// no error when creating user
				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				// query to fetch user by UUID was triggered & failed
				mockUserRepo.EXPECT().GetUserByUUID(ctx, gomock.Any()).Return(nil, mockRepoError).Times(1)

				// no user verification was created
				mockUserRepo.EXPECT().CreateUserVerification(ctx, gomock.Any()).Times(0)

				// no message was published
				mockPublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Times(0)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})

			It("should return error when repo fails to create user verification when creating email verification", func() {
				defer mockCtrl.Finish()

				mockRepoError := errors.New("failed to save user verification")

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    []byte{},
					JobTitle: "The Boss",
				}

				createdUser, err := user.New(user.UserParams{
					EntityParams: entity.EntityParams{
						EntityIDParams: entity.EntityIDParams{
							UUID:  id.NewUUID(),
							KeyID: id.NewKeyID(),
							XID:   id.NewXid(),
						},
						EntityTimestampParams: entity.EntityTimestampParams{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Metadata: map[string]any{},
					},
					Name:      request.Name,
					Email:     request.Email,
					ImageData: request.Image,
					Skills:    request.Skills,
					JobTitle:  request.JobTitle,
					Password:  "hashedPassword",
				})

				assert.NoError(t, err)

				// no error when creating user
				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				// query to fetch user by UUID was triggered & failed
				mockUserRepo.EXPECT().GetUserByUUID(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				// no user verification was created
				mockUserRepo.EXPECT().CreateUserVerification(ctx, gomock.Any()).Return(nil, mockRepoError).Times(1)

				// no message was published
				mockPublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Times(0)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})

			It("should return error when there is a failure to publish user email verification when creating email verification", func() {
				defer mockCtrl.Finish()

				mockError := errors.New("failed to publish user verification")

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    []byte{},
					JobTitle: "The Boss",
				}

				createdUser, err := user.New(user.UserParams{
					EntityParams: entity.EntityParams{
						EntityIDParams: entity.EntityIDParams{
							UUID:  id.NewUUID(),
							KeyID: id.NewKeyID(),
							XID:   id.NewXid(),
						},
						EntityTimestampParams: entity.EntityTimestampParams{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Metadata: map[string]any{},
					},
					Name:      request.Name,
					Email:     request.Email,
					ImageData: request.Image,
					Skills:    request.Skills,
					JobTitle:  request.JobTitle,
					Password:  "hashedPassword",
				})

				assert.NoError(t, err)

				// no error when creating user
				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				// query to fetch user by UUID was triggered & failed
				mockUserRepo.EXPECT().GetUserByUUID(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				userVerification := user.NewVerification(user.UserVerificationParams{
					ID:         id.NewUUID(),
					UserId:     createdUser.UUID(),
					Code:       "code",
					IsVerified: false,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				})

				// no user verification was created
				mockUserRepo.EXPECT().CreateUserVerification(ctx, gomock.Any()).Return(&userVerification, nil).Times(1)

				// no message was published
				mockPublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Return(mockError).Times(1)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})

		})
	})
})
