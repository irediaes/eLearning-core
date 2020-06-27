package storage

import (
	"github.com/biezhi/gorm-paginator/pagination"
	md "github.com/ebikode/eLearning-core/model"
)

// DBAssessmentStorage ...
type DBAssessmentStorage struct {
	*MDatabase
}

// NewDBAssessmentStorage Initialize Assessment Storage
func NewDBAssessmentStorage(db *MDatabase) *DBAssessmentStorage {
	return &DBAssessmentStorage{db}
}

// Get assessment using application_id and assessment_id
func (pdb *DBAssessmentStorage) Get(applicationID uint, assessmentID string) *md.Assessment {
	assessment := md.Assessment{}
	// Select resource from database
	err := pdb.db.
		Preload("Application").
		Preload("Application.Course").
		Preload("Application.User").
		Preload("Question").
		Where("application_id=? AND id=?", applicationID, assessmentID).First(&assessment).Error

	if len(assessment.ID) < 1 || err != nil {
		return nil
	}

	return &assessment
}

// GetLastAssessment ...
func (pdb *DBAssessmentStorage) GetLastAssessment() *md.Assessment {
	assessment := md.Assessment{}
	// Select resource from database
	err := pdb.db.
		Preload("Application").
		Preload("Application.Course").
		Preload("Application.User").
		Preload("Question").
		Order("created_at").
		Limit(1).First(&assessment).Error

	if len(assessment.ID) < 1 || err != nil {
		return nil
	}

	return &assessment
}

// GetAll Get all assessments
func (pdb *DBAssessmentStorage) GetAll(page, limit int) []*md.Assessment {
	var assessments []*md.Assessment
	// Select resource from database
	q := pdb.db.
		Preload("Application").
		Preload("Application.Course").
		Preload("Application.User").
		Preload("Question")

	pagination.Paging(&pagination.Param{
		DB:      q.Order("created_at desc").Find(&assessments),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &assessments)

	return assessments
}

// GetByUser Get all assessments of a application  form DB
func (pdb *DBAssessmentStorage) GetByUser(applicationID string, page, limit int) []*md.Assessment {
	var assessments []*md.Assessment
	// Select resource from database
	q := pdb.db.
		Preload("Application").
		Preload("Application.Course").
		Preload("Application.User").
		Preload("Question")

	pagination.Paging(&pagination.Param{
		DB:      q.Where("application_id=?", applicationID).Order("created_at desc").Find(&assessments),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &assessments)
	return assessments
}

// GetByCourse ...
func (pdb *DBAssessmentStorage) GetByCourse(courseID int) []*md.Assessment {
	var assessments []*md.Assessment
	// Select resource from database
	pdb.db.
		Preload("Application").
		Preload("Application.Course").
		Preload("Application.User").
		Preload("Question").
		Where("course_id=?", courseID).Order("created_at desc").Find(&assessments)

	return assessments
}

// GetSingleByCourse ...
func (pdb *DBAssessmentStorage) GetSingleByCourse(courseID int) *md.Assessment {
	assessment := md.Assessment{}
	// Select resource from database
	err := pdb.db.
		Preload("Application").
		Preload("Application.Course").
		Preload("Application.User").
		Preload("Question").
		Where("course_id=?", courseID).Order("created_at desc").
		First(&assessment).Error

	if len(assessment.ID) < 1 || err != nil {
		return nil
	}

	return &assessment
}

// Store Add a new assessment
func (pdb *DBAssessmentStorage) Store(p md.Assessment) (*md.Assessment, error) {

	py := p

	err := pdb.db.Create(&py).Error

	if err != nil {
		return nil, err
	}
	return pdb.Get(py.ApplicationID, py.ID), nil
}

// Update a assessment
func (pdb *DBAssessmentStorage) Update(assessment *md.Assessment) (*md.Assessment, error) {

	err := pdb.db.Save(&assessment).Error

	if err != nil {
		return nil, err
	}

	return assessment, nil
}

// Delete a assessment
func (pdb *DBAssessmentStorage) Delete(p *md.Assessment, isPermarnant bool) (bool, error) {

	// var err error
	// if isPermarnant {
	// 	err = pdb.db.Unscoped().Delete(p).Error
	// }
	// if !isPermarnant {
	// 	err = pdb.db.Delete(p).Error
	// }

	// if err != nil {
	// 	return false, err
	// }

	return true, nil
}