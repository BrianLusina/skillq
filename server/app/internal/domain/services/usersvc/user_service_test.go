package usersvc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	mockuserrepo "github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories/mocks"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	mockamqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher/mocks"
	mockstorageclient "github.com/BrianLusina/skillq/server/infra/storage/mocks"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")
}

var _ = Describe("User Service", func() {
	t := GinkgoT()

	var (
		mockCtrl             *gomock.Controller
		mockUserRepo         *mockuserrepo.MockUserRepoPort
		mockMessagePublisher *mockamqppublisher.MockAmqpEventPublisher
		mockStorageClient    *mockstorageclient.MockStorageClient
		userSvc              userService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(t)
		mockUserRepo = mockuserrepo.NewMockUserRepoPort(mockCtrl)
		mockMessagePublisher = mockamqppublisher.NewMockAmqpEventPublisher(mockCtrl)
		mockStorageClient = mockstorageclient.NewMockStorageClient(mockCtrl)
		userSvc = userService{
			userRepo:         mockUserRepo,
			messagePublisher: mockMessagePublisher,
			storageClient:    mockStorageClient,
		}

		assert.NotNil(t, userSvc)
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
					Image:    inbound.UserImageRequest{},
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
				imageRequest := inbound.UserImageRequest{
					Type:    "image/png",
					Content: "",
				}

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    imageRequest,
					JobTitle: "The Boss",
				}

				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, mockRepoError)

				// no query to fetch user by UUID was triggered
				mockUserRepo.EXPECT().GetUserByUUID(ctx, gomock.Any()).Times(0)

				// no message was published
				mockMessagePublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Times(0)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})

			It("should return error when publisher fails to publish verification event after creating user", func() {
				defer mockCtrl.Finish()

				mockPublisherError := errors.New("failed to publish user verification event")

				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    inbound.UserImageRequest{},
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
					Name:     request.Name,
					Email:    request.Email,
					Skills:   request.Skills,
					JobTitle: request.JobTitle,
					Password: "hashedPassword",
				})

				assert.NoError(t, err)

				// no error when creating user
				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				// message failed to publish
				mockMessagePublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Return(mockPublisherError).Times(1)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.Nil(t, actualUser)
				assert.Error(t, actualErr)
			})

			It("should return nil error when there is success creating a user, publishing user email verification event and store image task", func() {
				defer mockCtrl.Finish()
				imageRequest := inbound.UserImageRequest{
					Type:    "image/png",
					Content: "",
				}
				request := inbound.UserRequest{
					Name:     "John Doe",
					Email:    "fake@example.com",
					Password: "password",
					Skills:   []string{},
					Image:    imageRequest,
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
					Name:     request.Name,
					Email:    request.Email,
					Skills:   request.Skills,
					JobTitle: request.JobTitle,
					Password: "hashedPassword",
				})

				assert.NoError(t, err)

				// no error when creating user
				mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&createdUser, nil).Times(1)

				// message failed to be published
				mockMessagePublisher.EXPECT().Publish(ctx, gomock.Any(), gomock.Any()).Return(nil).Times(2)

				actualUser, actualErr := userSvc.CreateUser(context.Background(), request)
				assert.NotNil(t, actualUser)
				assert.NoError(t, actualErr)
			})
		})
	})

	Context("Fetching an existing user", func() {
		When("by UUID", func() {
			It("should return error if there is a failure retrieving from repository", func() {
				defer mockCtrl.Finish()

				expectedErr := errors.New("failed to retrieve user")
				userUUID := id.NewUUID()

				mockUserRepo.EXPECT().GetUserByUUID(ctx, userUUID).Return(nil, expectedErr).Times(1)

				actualUser, actualErr := userSvc.GetUserByUUID(ctx, userUUID.String())
				Expect(actualUser).To(BeNil())
				assert.Error(t, actualErr)
			})

			It("should return user response if there is success retrieving from repository", func() {
				defer mockCtrl.Finish()
				name := faker.FirstName()
				email := faker.Email()
				imageUrl := faker.URL()
				skill := faker.Word()
				jobTitle := faker.Word()
				password := faker.Password()

				existingUser, err := user.New(user.UserParams{
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
					Name:     name,
					Email:    email,
					ImageUrl: imageUrl,
					Skills:   []string{skill},
					JobTitle: jobTitle,
					Password: password,
				})
				assert.NoError(t, err)

				userUUID := id.NewUUID()

				mockUserRepo.EXPECT().GetUserByUUID(ctx, userUUID).Return(&existingUser, nil).Times(1)

				actualUser, actualErr := userSvc.GetUserByUUID(ctx, userUUID.String())
				assert.Nil(t, actualErr)
				Expect(actualUser).NotTo(BeNil())

			})
		})
	})
})
