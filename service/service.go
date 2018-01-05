package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/hoangbeatsb3/ttcourse/config"
	"github.com/hoangbeatsb3/ttcourse/model"
	// "github.com/hoangbeatsb3/ttcourse/repository"
	"../repository"
	"github.com/prometheus/common/log"
)

var cfg = config.LoadEnvConfig()
var repo = repository.NewRepo()

func FindAllCourses(w http.ResponseWriter, r *http.Request) {

	courses, err := repo.FindAllCourses()
	HandleError(err)

	if err := json.NewEncoder(w).Encode(courses); err != nil {
		panic(err)
	}
}

func FindCoursesByName(w http.ResponseWriter, r *http.Request) {

	parm := chi.URLParam(r, "name")

	courses, err := repo.FindCourseByName(strings.ToLower(parm))
	HandleError(err)
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		panic(err)
	}
}

func FindCourseByAlias(w http.ResponseWriter, r *http.Request) {

	parm := chi.URLParam(r, "alias")

	course, err := repo.FindCourseByAlias(strings.ToLower(parm))
	HandleError(err)
	if err := json.NewEncoder(w).Encode(course); err != nil {
		panic(err)
	}
}

func FindHighestVote(w http.ResponseWriter, r *http.Request) {

	courses, err := repo.FindAllCourses()
	HandleError(err)
	course := courses[0]
	flag := 1

	for flag < len(courses) {
		if courses[flag].Vote > course.Vote {
			course = courses[flag]
		}
		flag = flag + 1
	}

	if err := json.NewEncoder(w).Encode(course); err != nil {
		panic(err)
	}
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course model.Course
	body, err := getParams(r)
	HandleError(err)

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &course); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	course.Alias = strings.ToLower(course.Name)
	course.Vote = 0
	if course.Participant != nil {
		course.Vote = 1
	}
	repo.CreateCourse(course)

}

func VoteCourse(w http.ResponseWriter, r *http.Request) {

	var course model.Course
	body, err := ioutil.ReadAll(r.Body)
	HandleError(err)

	if err := json.Unmarshal(body, &course); err != nil {
		panic(err)
	}

	courseExists, err := repo.CheckIfExists(course)
	// HandleError(err)

	if err != nil {
		repo.CreateCourse(course)
		log.Info("Create new course: ", course)
	} else {

		log.Info("Vote for course: ", course)
		participant := model.Participant{
			Id:    len(courseExists.Participant),
			Name:  course.Participant[0].Name,
			Email: course.Participant[0].Email,
		}

		courseExists.Participant = append(courseExists.Participant, participant)
		courseExists.Vote += 1

		repo.Vote(courseExists)
	}
}

func getParams(r *http.Request) ([]byte, error) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
