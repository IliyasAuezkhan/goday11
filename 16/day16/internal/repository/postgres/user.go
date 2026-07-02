package postgres
import (
	"errors"
	"fmt"
	"time"
	"gorm.io/gorm"
	"day16/internal/domain"
)

type UserDBModel struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`
	Email string `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password string `gorm:"not null;type:varchar(255)"`
	CreatedAt time.Time `gorm:"not null"`
}

func (UserDBModel) TableName() string {
	return "users"
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user *domain.User) error {
	dbModel := &UserDBModel{
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
	}

	result := r.db.Create(dbModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("user with this email already exists")
		}
		return fmt.Errorf("failed to create user in db: %w", result.Error)
	}

	user.ID = dbModel.ID
	user.CreatedAt = dbModel.CreatedAt
	return nil
}

func (r *PostgresUserRepository) FindByID(id int64) (*domain.User, error) {
	var dbModel UserDBModel
	result := r.db.First(&dbModel, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return &domain.User{
		ID:        dbModel.ID,
		Email:     dbModel.Email,
		Password:  dbModel.Password,
		CreatedAt: dbModel.CreatedAt,
	}, nil
}

func (r *PostgresUserRepository) Update(id int64, fields *domain.UpdateUserFields) error {
	updates := make(map[string]interface{})
	if fields.Email != nil {
		updates["email"] = *fields.Email
	}
	if fields.Password != nil {
		updates["password"] = *fields.Password
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.Model(&UserDBModel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("email already taken")
		}
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *PostgresUserRepository) Delete(id int64) error {
	result := r.db.Delete(&UserDBModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}