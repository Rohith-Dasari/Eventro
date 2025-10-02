package authorisation

import (
	"context"
	"errors"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"reflect"
	"testing"

	"eventro2/mocks"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		userRepo userrepository.UserRepository
	}
	tests := []struct {
		name string
		args args
		want *AuthService
	}{
		{
			name: "valid constructor",
			args: args{
				userRepo: mocks.NewMockUserRepository(ctrl),
			},
			want: &AuthService{
				UserRepo: mocks.NewMockUserRepository(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.userRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_ValidateLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securePass123"), bcrypt.DefaultCost)
	validUser := &models.User{
		UserID:    "user-1",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		IsBlocked: false,
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      models.User
		wantErr   bool
	}{
		{
			name: "successful login",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "securePass123",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("test@example.com").
					Return(validUser, nil)
			},
			want:    *validUser,
			wantErr: false,
		},
		{
			name: "user not found",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:      context.Background(),
				email:    "unknown@example.com",
				password: "irrelevant",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("unknown@example.com").
					Return(nil, errors.New("not found"))
			},
			want:    models.User{},
			wantErr: true,
		},
		{
			name: "incorrect password",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "wrongPassword",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("test@example.com").
					Return(validUser, nil)
			},
			want:    models.User{},
			wantErr: true,
		},
		{
			name: "blocked user",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:      context.Background(),
				email:    "blocked@example.com",
				password: "securePass123",
			},
			mockSetup: func() {
				blockedUser := *validUser
				blockedUser.Email = "blocked@example.com"
				blockedUser.IsBlocked = true
				mockUserRepo.EXPECT().
					GetByEmail("blocked@example.com").
					Return(&blockedUser, nil)
			},
			want:    models.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			a := &AuthService{
				UserRepo: tt.fields.UserRepo,
			}
			got, err := a.ValidateLogin(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateLogin() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestAuthService_HashPassword(t *testing.T) {
	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successfully hashes password",
			fields: fields{
				UserRepo: nil,
			},
			args: args{
				password: "mySecurePassword123!",
			},
			wantErr: false,
		},
		{
			name: "empty password",
			fields: fields{
				UserRepo: nil,
			},
			args: args{
				password: "",
			},
			wantErr: false, // bcrypt allows empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				UserRepo: tt.fields.UserRepo,
			}
			got, err := a.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate that returned hash matches the original password
			if err := bcrypt.CompareHashAndPassword([]byte(got), []byte(tt.args.password)); err != nil {
				t.Errorf("Hashed password does not match original: %v", err)
			}
		})
	}
}

func TestAuthService_IsValidEmail(t *testing.T) {
	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		email string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid email",
			args: args{
				email: "user@example.com",
			},
			want: true,
		},
		{
			name: "missing @ symbol",
			args: args{
				email: "userexample.com",
			},
			want: false,
		},
		{
			name: "missing domain",
			args: args{
				email: "user@",
			},
			want: false,
		},
		{
			name: "empty email",
			args: args{
				email: "",
			},
			want: false,
		},
		{
			name: "email with spaces",
			args: args{
				email: "user @example.com",
			},
			want: false,
		},
		{
			name: "email with subdomain",
			args: args{
				email: "user@mail.example.com",
			},
			want: true,
		},
		{
			name: "email with plus",
			args: args{
				email: "user+tag@example.com",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				UserRepo: tt.fields.UserRepo,
			}
			if got := a.IsValidEmail(tt.args.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_IsValidPhoneNumber(t *testing.T) {
	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		phone string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid phone number with plus",
			args: args{phone: "+1234567890"},
			want: true,
		},
		{
			name: "valid without plus",
			args: args{phone: "1234567890"},
			want: true,
		},
		{
			name: "too short",
			args: args{phone: "12"},
			want: false,
		},
		{
			name: "starts with zero",
			args: args{phone: "+0123456789"},
			want: false,
		},
		{
			name: "contains letters",
			args: args{phone: "+123ABC456"},
			want: false,
		},
		{
			name: "contains symbols",
			args: args{phone: "+123-456-7890"},
			want: false,
		},
		{
			name: "empty string",
			args: args{phone: ""},
			want: false,
		},
		{
			name: "too long",
			args: args{phone: "+12345678901234567"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				UserRepo: tt.fields.UserRepo,
			}
			if got := a.IsValidPhoneNumber(tt.args.phone); got != tt.want {
				t.Errorf("IsValidPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_IsValidPassword(t *testing.T) {
	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		args   args
		want   bool
		fields fields
	}{
		{
			name: "valid password with all requirements",
			args: args{password: "StrongPass123!"},
			want: true,
		},
		{
			name: "too short",
			args: args{password: "Shrt1!"},
			want: false,
		},
		{
			name: "missing number",
			args: args{password: "StrongPass!!"},
			want: false,
		},
		{
			name: "missing uppercase",
			args: args{password: "weakpass123!"},
			want: false,
		},
		{
			name: "missing lowercase",
			args: args{password: "STRONG123!"},
			want: false,
		},
		{
			name: "missing symbol",
			args: args{password: "StrongPass123"},
			want: false,
		},
		{
			name: "empty string",
			args: args{password: ""},
			want: false,
		},
		{
			name: "only special characters",
			args: args{password: "!@#$%^&*()-+"},
			want: false,
		},
		{
			name: "just long enough and valid",
			args: args{password: "Aabc1234!@#$"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				UserRepo: tt.fields.UserRepo,
			}
			if got := a.IsValidPassword(tt.args.password); got != tt.want {
				t.Errorf("IsValidPassword(%q) = %v, want %v", tt.args.password, got, tt.want)
			}
		})
	}
}

func TestAuthService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	type fields struct {
		UserRepo userrepository.UserRepository
	}
	type args struct {
		ctx         context.Context
		username    string
		email       string
		phoneNumber string
		password    string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful signup",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:         context.Background(),
				username:    "john_doe",
				email:       "john@example.com",
				phoneNumber: "+12345678901",
				password:    "SecurePass123!",
			},
			mockSetup: func() {
				// Email doesn't exist
				mockUserRepo.EXPECT().
					GetByEmail("john@example.com").
					Return(nil, errors.New("not found"))

				mockUserRepo.EXPECT().
					Create(gomock.Any()).
					DoAndReturn(func(user *models.User) error {
						if user.Email != "john@example.com" {
							return errors.New("email mismatch")
						}
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "invalid email format",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:         context.Background(),
				username:    "user1",
				email:       "invalid-email",
				phoneNumber: "+12345678901",
				password:    "ValidPass123!",
			},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name: "email already exists",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:         context.Background(),
				username:    "user1",
				email:       "exists@example.com",
				phoneNumber: "+12345678901",
				password:    "ValidPass123!",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("exists@example.com").
					Return(&models.User{Email: "exists@example.com"}, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:         context.Background(),
				username:    "user2",
				email:       "new@example.com",
				phoneNumber: "+12345678901",
				password:    "weak",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("new@example.com").
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "invalid phone number",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:         context.Background(),
				username:    "user3",
				email:       "user3@example.com",
				phoneNumber: "123abc",
				password:    "ValidPass123!",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("user3@example.com").
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "repo create fails",
			fields: fields{
				UserRepo: mockUserRepo,
			},
			args: args{
				ctx:         context.Background(),
				username:    "failuser",
				email:       "fail@example.com",
				phoneNumber: "+12345678901",
				password:    "StrongPass123!",
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetByEmail("fail@example.com").
					Return(nil, errors.New("not found"))
				mockUserRepo.EXPECT().
					Create(gomock.Any()).
					Return(errors.New("db failure"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			a := &AuthService{
				UserRepo: tt.fields.UserRepo,
			}
			got, err := a.Signup(tt.args.ctx, tt.args.username, tt.args.email, tt.args.phoneNumber, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.Email != tt.args.email {
				t.Errorf("Signup() got = %+v, want email %v", got, tt.args.email)
			}
		})
	}
}
