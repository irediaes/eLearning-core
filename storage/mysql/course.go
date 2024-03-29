package storage

import (
	"github.com/biezhi/gorm-paginator/pagination"
	md "github.com/ebikode/eLearning-core/model"
	ut "github.com/ebikode/eLearning-core/utils"
)

// DBCourseStorage encapsulates DB Connection Model
type DBCourseStorage struct {
	*MDatabase
}

// NewDBCourseStorage Initialize Course Storage
func NewDBCourseStorage(db *MDatabase) *DBCourseStorage {
	return &DBCourseStorage{db}
}

// Get Fetch Single Course fron DB
func (cdb *DBCourseStorage) Get(id uint) *md.Course {
	course := md.Course{}
	// Select resource from database
	err := cdb.db.
		Preload("User").
		Where("courses.id=?", id).First(&course).Error

	if course.ID < 1 || err != nil {
		return nil
	}

	return &course
}

// GetSingleByUser Fetch Single Course fron DB
func (cdb *DBCourseStorage) GetSingleByUser(userID string, courseID uint) *md.Course {
	course := md.Course{}
	// Select resource from database
	err := cdb.db.
		Preload("User").
		Where("user_id=? AND id=?", userID, courseID).First(&course).Error

	if course.ID < 1 || err != nil {
		return nil
	}

	return &course
}

// GetAll Fetch all courses from DB
func (cdb *DBCourseStorage) GetAll(page, limit int, userType string) []*md.Course {
	var courses []*md.Course

	if userType == ut.AdminRole {
		pagination.Paging(&pagination.Param{
			DB: cdb.db.
				Preload("User").
				Order("created_at desc").
				Find(&courses),
			Page:    page,
			Limit:   limit,
			OrderBy: []string{"created_at desc"},
		}, &courses)

		return courses
	}

	pagination.Paging(&pagination.Param{
		DB: cdb.db.
			Preload("User").
			Where("status=?", ut.Approved).
			Order("created_at desc").
			Find(&courses),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &courses)

	return courses

}

// GetByUser Fetch all user' courses from DB
func (cdb *DBCourseStorage) GetByUser(userID string, page, limit int) []*md.Course {
	var courses []*md.Course

	pagination.Paging(&pagination.Param{
		DB: cdb.db.
			Preload("User").
			Where("user_id=?", userID).
			Find(&courses),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &courses)

	return courses
}

// Store Add a new course
func (cdb *DBCourseStorage) Store(p md.Course) (*md.Course, error) {

	course := p

	err := cdb.db.Create(&course).Error

	if err != nil {
		return nil, err
	}
	return cdb.Get(course.ID), nil
}

// Update a course
func (cdb *DBCourseStorage) Update(course *md.Course) (*md.Course, error) {

	err := cdb.db.Save(&course).Error

	if err != nil {
		return nil, err
	}

	return course, nil
}

// Delete a course
func (cdb *DBCourseStorage) Delete(c md.Course, isPermarnant bool) (bool, error) {

	var err error
	if isPermarnant {
		err = cdb.db.Unscoped().Delete(c).Error
	}
	if !isPermarnant {
		err = cdb.db.Delete(c).Error
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
