package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	ase "github.com/ebikode/eLearning-core/domain/assessment"
	que "github.com/ebikode/eLearning-core/domain/question"
	usr "github.com/ebikode/eLearning-core/domain/user"
	md "github.com/ebikode/eLearning-core/model"
	tr "github.com/ebikode/eLearning-core/translation"
	ut "github.com/ebikode/eLearning-core/utils"
)

// GetAssessmentEndpoint fetch a single assessment
// func GetAssessmentEndpoint(asr ase.AssessmentService) http.HandlerFunc {

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Get User Token Data
// 		tokenData := r.Context().Value("tokenData").(*md.UserTokenData)
// 		userID := string(tokenData.UserID)

// 		assessmentID := chi.URLParam(r, "assessmentID")
// 		applicationID := chi.URLParam(r, "applicationID")

// 		var assessment *md.Assessment
// 		assessment = asr.GetAssessment(applicationID, assessmentID)

// 		resp := ut.Message(true, "")
// 		resp["assessment"] = assessment
// 		ut.Respond(w, r, resp)
// 	}
// }

// GetAdminAssessmentsEndpoint fetch a single assessment
func GetAdminAssessmentsEndpoint(asr ase.AssessmentService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page, limit := ut.PaginationParams(r)

		assessments := asr.GetAssessments(page, limit)

		var nextPage int
		if len(assessments) == limit {
			nextPage = page + 1
		}

		resp := ut.Message(true, "")
		resp["current_page"] = page
		resp["next_page"] = nextPage
		resp["limit"] = limit
		resp["assessments"] = assessments
		ut.Respond(w, r, resp)
	}

}

// CreateAssessmentEndpoint ...
func CreateAssessmentEndpoint(asr ase.AssessmentService, uss usr.UserService, qs que.QuestionService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get User Token Data
		tokenData := r.Context().Value("tokenData").(*md.UserTokenData)
		userID := string(tokenData.UserID)
		payloads := []ase.Payload{}
		err := json.NewDecoder(r.Body).Decode(&payloads)
		fmt.Println("second Error check", err)

		tParam := tr.TParam{
			Key:          "error.request_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		if err != nil {
			// Respond with an errortra nslated

			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		//  if the length of the payloads uploaded is less than 1
		if len(payloads) < 1 {
			tParam = tr.TParam{
				Key:          "validation.lesser",
				TemplateData: map[string]interface{}{"Min": 1},
				PluralCount:  nil,
			}
			// Respond with an error translated
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		if len(payloads) > 50 {
			tParam = tr.TParam{
				Key:          "validation.greater",
				TemplateData: map[string]interface{}{"Max": 50},
				PluralCount:  nil,
			}
			// Respond with an error translated
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}
		// assessments created that will be returned
		var createdAssessments []*md.Assessment

		checkUser := uss.GetUser(userID)

		if checkUser == nil {
			tParam = tr.TParam{
				Key:          "error.user_not_found",
				TemplateData: nil,
				PluralCount:  nil,
			}
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		// loop through the payloads and create them
		for _, v := range payloads {

			// Validate assessment input
			err = ase.Validate(v, r)
			if err != nil {
				validationFields := ase.ValidationFields{}
				fmt.Println("third Error check", validationFields)
				b, _ := json.Marshal(err)
				// Respond with an errortranslated
				resp := ut.Message(false, ut.Translate(tParam, r))
				json.Unmarshal(b, &validationFields)
				resp["error"] = validationFields
				ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
				return

			}

			question := qs.GetQuestion(v.QuestionID)

			isCorrect := false

			if question.Answer == v.SelectedAnswer {
				isCorrect = true
			}

			assessment := md.Assessment{
				ApplicationID:  v.ApplicationID,
				QuestionID:     v.QuestionID,
				SelectedAnswer: v.SelectedAnswer,
				CorrectAnswer:  question.Answer,
				IsCorrect:      isCorrect,
			}

			// Create a assessment
			newAssessment, errParam, err := asr.CreateAssessment(assessment)
			if err != nil {
				// Check if the error is dupliassessmention error
				cErr := ut.CheckUniqueError(r, err)
				if cErr != nil {
					ut.ErrorRespond(http.StatusBadRequest, w, r, ut.Message(false, cErr.Error()))
					return
				}
				// Respond with an errortranslated
				ut.ErrorRespond(http.StatusBadRequest, w, r, ut.Message(false, ut.Translate(errParam, r)))
				return
			}
			// add the created category to the slice
			createdAssessments = append(createdAssessments, newAssessment)
		}

		tParam = tr.TParam{
			Key:          "general.resource_created",
			TemplateData: nil,
			PluralCount:  nil,
		}

		resp := ut.Message(true, ut.Translate(tParam, r))
		resp["assessments"] = createdAssessments
		ut.Respond(w, r, resp)

	}

}