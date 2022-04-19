package main

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/mail"
	"reflect"
)

var DB = &Storage{}

func SetupStorage() error {
	dsn := "host=postgres user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	var err error
	DB.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

type GenderType string

const (
	MALE   GenderType = "m"
	FEMALE GenderType = "f"
)

type UserProfile struct {
	Username   *string     `json:"username,omitempty" gorm:"username"`
	Email      *string     `json:"email,omitempty" gorm:"email"`
	Gender     *GenderType `json:"gender,omitempty" gorm:"gender"`
	WinCount   *int        `json:"win_count,omitempty" gorm:"win_count"`
	LoseCount  *int        `json:"lose_count,omitempty" gorm:"lose_count"`
	TimeInGame *int        `json:"time_in_game,omitempty" gorm:"time_in_game"`
}

func (profile *UserProfile) validateCreate() error {
	if profile.Username == nil {
		return errors.New("username should not be empty")
	}
	if profile.Email == nil {
		return errors.New("email should not be empty")
	} else if !validateEmail(*profile.Email) {
		return errors.New("not valid email")
	}
	if profile.Gender == nil {
		return errors.New("gender should not be empty")
	} else if *profile.Gender != MALE && *profile.Gender != FEMALE {
		return errors.New("not valid gender")
	}
	return nil
}

func (profile *UserProfile) getNotNilFields() []string {
	var res []string
	profileReflect := reflect.ValueOf(*profile)
	for i := 0; i < profileReflect.NumField(); i++ {
		field := profileReflect.Field(i)
		if !field.IsNil() {
			res = append(res, profileReflect.Type().Field(i).Tag.Get("gorm"))
		}
	}
	return res
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) GetUsers() ([]UserProfile, error) {
	var users []UserProfile
	res := s.db.Find(&users)
	return users, res.Error
}

func (s *Storage) CreateUser(user *UserProfile) error {
	if err := user.validateCreate(); err != nil {
		return err
	}
	columnsToSave := user.getNotNilFields()
	res := s.db.Select(columnsToSave).Create(user)
	return res.Error
}

func (s *Storage) GetUser(username string) (*UserProfile, error) {
	user := &UserProfile{}
	res := s.db.Where("username = ?", username).Find(user)
	if res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, nil
	}
	return user, nil
}

func (s *Storage) UpdateUser(username string, patch *UserProfile) (int64, error) {
	columnsToSave := patch.getNotNilFields()
	res := s.db.Select(columnsToSave).Where("username = ?", username).Updates(patch)
	return res.RowsAffected, res.Error
}

func (s *Storage) DeleteUser(username string) (int64, error) {
	res := s.db.Where("username = ?", username).Delete(&UserProfile{})
	return res.RowsAffected, res.Error
}

func (s *Storage) CreateJob() (uint64, error) {
	pdf := &PDFJob{}
	res := s.db.Create(pdf)
	return pdf.ID, res.Error
}

func (s *Storage) UpdateJob(pdf PDFJob) (int64, error) {
	res := s.db.Select("file_path").Where("id = ?", pdf.ID).Updates(pdf)
	return res.RowsAffected, res.Error
}

func (s *Storage) GetJob(id int) (*PDFJob, error) {
	pdf := &PDFJob{}
	res := s.db.Where("id = ?", id).Find(&pdf)
	if res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, nil
	}
	return pdf, nil
}
