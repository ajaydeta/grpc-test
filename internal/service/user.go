package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"tablelink/domain"
	"tablelink/internal/repository"
	pb "tablelink/proto/pb/proto"
	"tablelink/util"
	"time"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

const secret = "AllYourBase"

func (i *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := i.repo.FindUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword())); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	rr, err := i.repo.FindRoleRightByRoleId(ctx, user.RoleID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"role_id":  user.RoleID,
		"r_create": rr.RCreate,
		"r_read":   rr.RRead,
		"r_update": rr.RUpdate,
		"r_delete": rr.RDelete,
	})

	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	return &pb.LoginResponse{
		AccessToken: ss,
	}, nil
}

func (i *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {

	user, err := i.repo.FindUserById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	role, err := i.repo.FindRoleById(ctx, user.RoleID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	res := pb.User{
		RoleId:   user.RoleID,
		RoleName: role.Name,
		Name:     user.Name,
		Email:    user.Email,
	}

	if user.LastAccess != nil {
		res.LastAccess = user.LastAccess.Format(time.RFC3339)
	}

	return &res, nil
}

func (i *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.SuccessStatusResponse, error) {
	user, err := i.repo.FindUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if user != nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already exist")
	}

	userInsert := domain.User{
		RoleID: req.GetRoleId(),
		Email:  req.GetEmail(),
		Name:   req.GetName(),
	}

	userInsert.Password, err = util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	err = i.repo.CreateUser(ctx, &userInsert)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	return &pb.SuccessStatusResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (i *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.SuccessStatusResponse, error) {
	user, err := i.repo.FindUserById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	userUpdate := domain.User{
		ID:   req.GetId(),
		Name: req.GetName(),
	}

	err = i.repo.UpdateNameUser(ctx, &userUpdate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	return &pb.SuccessStatusResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (i *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.SuccessStatusResponse, error) {
	user, err := i.repo.FindUserById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	err = i.repo.DeleteUserById(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	return &pb.SuccessStatusResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func validateJwt(ctx context.Context) (map[string]int, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	ls := md.Get("x-link-service")
	if len(ls) == 0 {
		return nil, errors.New("empty x-link-service")
	}

	tokenRaw := md.Get("Authorization")
	if len(tokenRaw) == 0 {
		return nil, errors.New("empty authorization")
	}

	token :=

	token, err := jwt.Parse(, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
//switch {
//case token.Valid:
//fmt.Println("You look nice today")
//case errors.Is(err, jwt.ErrTokenMalformed):
//fmt.Println("That's not even a token")
//case errors.Is(err, jwt.ErrTokenSignatureInvalid): // Invalid signature     fmt. Println("Invalid signature") case errors. Is(err, jwt. ErrTokenExpired) || errors. Is(err, jwt. ErrTokenNotValidYet):     // Token is either expired or not active yet     fmt. Println("Timing is everything") default:     fmt. Println("Couldn't handle this token:", err) }

}
