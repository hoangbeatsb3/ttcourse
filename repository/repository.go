package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hoangbeatsb3/ttcourse/model"
	"github.com/garyburd/redigo/redis"
	"github.com/prometheus/common/log"
)

type Repo struct {
	c *redis.Conn
}

func NewRepo(port string) *Repo {
	c, err := redis.Dial("tcp", port)
	HandleError(err)

	r := &Repo{
		c: &c,
	}
	return r
}

func (r *Repo) FindAllCourses() model.Courses {

	c := *r.c
	keys, err := c.Do("KEYS", "trainingtoolcourse:*")

	HandleError(err)

	var courses model.Courses
	for _, k := range keys.([]interface{}) {

		var course model.Course
		reply, err := c.Do("GET", k.([]byte))
		HandleError(err)
		if err := json.Unmarshal(reply.([]byte), &course); err != nil {
			panic(err)
		}
		courses = append(courses, course)
	}

	return courses
}

func (r *Repo) FindCourseByName(name string) model.Courses {

	c := *r.c

	var courses model.Courses
	keys, err := r.GetKeys("trainingtoolcourse:*" + name + "*")
	if err != nil {
		panic(err)
	}

	if len(keys) <= 0 {
		return courses
	}
	for _, v := range keys {
		v = strings.Replace(v, "trainingtoolcourse:", "", -1)
		var course model.Course
		reply, err := c.Do("GET", "trainingtoolcourse:"+v)
		HandleError(err)
		log.Info("GET OK")
		if err = json.Unmarshal(reply.([]byte), &course); err != nil {
			panic(err)
		}
		courses = append(courses, course)
	}

	return courses
}

func (r *Repo) GetKeys(pattern string) ([]string, error) {

	c := *r.c

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(c.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (r *Repo) FindCourseByAlias(name string) model.Course {

	c := *r.c
	var course model.Course

	reply, err := c.Do("GET", "trainingtoolcourse:"+name)
	HandleError(err)

	if err = json.Unmarshal(reply.([]byte), &course); err != nil {
		panic(err)
	}

	log.Info("Reply OK")
	return course
}

func (r *Repo) CreateCourse(p model.Course) {

	c := *r.c

	p.Alias = strings.ToLower(p.Name)
	p.Vote = 0
	if p.Participant != nil {
		p.Vote = 1
	}

	b, err := json.Marshal(p)
	HandleError(err)
	reply, err := c.Do("SET", "trainingtoolcourse:"+p.Alias, b)
	HandleError(err)
	log.Info("Reply ", reply)
}

func (r *Repo) CheckIfExists(p model.Course) (bool, model.Course) {
	c := *r.c
	reply, err := c.Do("GET", "trainingtoolcourse:"+strings.ToLower(p.Name))
	if err != nil {
		panic(err)
	}

	if reply == nil {
		return false, p // false means doesn't exist
	}

	var course model.Course
	if err = json.Unmarshal(reply.([]byte), &course); err != nil {
		panic(err)
	}
	return true, course
}

func (r *Repo) Vote(p model.Course) {

	c := *r.c

	b, err := json.Marshal(p)
	reply, err := c.Do("SET", "trainingtoolcourse:"+strings.ToLower(p.Name), b)
	log.Info("Vote Status: ", reply)
	HandleError(err)
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
