package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/hoangbeatsb3/ttcourse/model"
	"github.com/lnquy/fugu/config"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	c *redis.Conn
}

var courseKey = "ttcourse:"

func NewRepo(cfg *config.Config) (*Repo, error) {
	c, err := redis.Dial("tcp", cfg.Server.GetRedisPort())
	if err != nil {
		return nil, err
	}

	return &Repo{
		c: &c,
	}, nil
}

func (r *Repo) FindAllCourses() (model.Courses, error) {

	c := *r.c
	keys, err := c.Do("KEYS", courseKey+"*")

	if err != nil {
		return nil, err
	}

	var courses model.Courses
	for _, k := range keys.([]interface{}) {

		var course model.Course
		reply, err := c.Do("GET", k.([]byte))

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(reply.([]byte), &course); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r *Repo) FindCourseByName(name string) (model.Courses, error) {

	c := *r.c

	var courses model.Courses
	keys, err := r.GetKeys(courseKey + "*" + name + "*")
	if err != nil {
		return nil, err
	}

	for _, v := range keys {
		v = strings.Replace(v, courseKey, "", -1)
		var course model.Course
		reply, err := c.Do("GET", courseKey+v)
		if err != nil {
			return nil, err
		}
		logrus.Info("GET OK")
		if err = json.Unmarshal(reply.([]byte), &course); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
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

		iter, err := redis.Int(arr[0], nil)
		if err != nil {
			return nil, err
		}
		k, err := redis.Strings(arr[1], nil)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (r *Repo) FindCourseByAlias(name string) (model.Course, error) {

	c := *r.c
	var course model.Course

	reply, err := c.Do("GET", courseKey+name)
	if err != nil {
		return course, err
	}

	if err = json.Unmarshal(reply.([]byte), &course); err != nil {
		return course, err
	}

	logrus.Info("Reply OK")
	return course, nil
}

func (r *Repo) CreateCourse(p model.Course) (model.Course, error) {

	c := *r.c

	b, err := json.Marshal(p)
	if err != nil {
		return p, err
	}
	reply, err := c.Do("SET", courseKey+p.Alias, b)
	if err != nil {
		return p, err
	}
	logrus.Info("Reply ", reply)
	return p, nil
}

func (r *Repo) CheckIfExists(p model.Course) (model.Course, error) {
	c := *r.c

	var course model.Course
	reply, err := c.Do("GET", courseKey+strings.ToLower(p.Name))
	if err != nil {
		return course, err
	}

	if reply == nil {
		return course, err
	}

	if err = json.Unmarshal(reply.([]byte), &course); err != nil {
		return course, err
	}
	return course, nil
}

func (r *Repo) Vote(p model.Course) (model.Course, error) {

	c := *r.c

	b, err := json.Marshal(p)
	reply, err := c.Do("SET", courseKey+strings.ToLower(p.Name), b)
	logrus.Info("Vote Status: ", reply)
	if err != nil {
		return p, err
	}

	return p, nil
}
