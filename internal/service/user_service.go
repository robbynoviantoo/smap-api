package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"smap-api/internal/config"
	"smap-api/internal/model"
	"smap-api/internal/repository"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// ambil roles dari DB
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
		Roles: roles, // kirim ke frontend
	}, nil
}

func generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(config.App.JWTExpireHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.App.JWTSecret))
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.UserResponse, error) {
	return s.repo.GetAllUsers(ctx)
}

// GetMe mengembalikan data user berdasarkan ID dari token.
func (s *UserService) GetMe(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CompleteOnboarding memproses pengisian data awal user (first login).
func (s *UserService) CompleteOnboarding(ctx context.Context, userID uint, req *model.OnboardingRequest) error {
	if req.FirstName == "" {
		return errors.New("first_name wajib diisi")
	}
	if req.NoHandphone == "" {
		return errors.New("no_handphone wajib diisi")
	}
	if req.Password == "" {
		return errors.New("password wajib diisi")
	}
	if len(req.Password) < 8 {
		return errors.New("password minimal 8 karakter")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if !user.IsFirstLogin {
		return errors.New("onboarding sudah diselesaikan sebelumnya")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password: %w", err)
	}

	// Proses foto (base64 data URL)
	var imageFilename string
	if req.PhotoCropped != "" {
		imageFilename, err = saveBase64Image(req.PhotoCropped, userID)
		if err != nil {
			return fmt.Errorf("gagal menyimpan foto: %w", err)
		}
	}

	return s.repo.UpdateOnboarding(ctx, userID, req.FirstName, req.LastName, req.NoHandphone, string(hashed), imageFilename)
}

// saveBase64Image mem-parse dan menyimpan base64 data URL sebagai file .webp di ./uploads/profile/.
func saveBase64Image(dataURL string, userID uint) (string, error) {
	// Contoh: "data:image/png;base64,iVBOR..."
	re := regexp.MustCompile(`^data:image/(\w+);base64,`)
	matches := re.FindStringSubmatch(dataURL)
	if len(matches) < 2 {
		return "", errors.New("format image tidak valid")
	}

	b64 := re.ReplaceAllString(dataURL, "")
	b64 = strings.TrimSpace(b64)

	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		// coba URL-safe encoding
		data, err = base64.URLEncoding.DecodeString(b64)
		if err != nil {
			return "", errors.New("gagal decode base64")
		}
	}

	dir := filepath.Join(".", "uploads", "profile")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Simpan dengan ekstensi asli; untuk konversi webp butuh library tambahan.
	// Kita simpan sebagai ekstensi aslinya saja untuk simplisitas.
	ext := matches[1]
	filename := fmt.Sprintf("profile_%d_%d.%s", userID, time.Now().UnixMilli(), ext)
	dst := filepath.Join(dir, filename)

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return "", err
	}
	return filename, nil
}
