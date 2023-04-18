package handler

import (
	"errors"
	"github.com/GarnBarn/common-go/httpserver"
	"github.com/GarnBarn/gb-tag-service/model"
	"github.com/GarnBarn/gb-tag-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Tag struct {
	validate   validator.Validate
	tagService service.Tag
}

var (
	ErrGinBadRequestBody = gin.H{"message": "bad request body."}
)

func NewTagHandler(validate validator.Validate, tagService service.Tag) Tag {
	return Tag{
		validate:   validate,
		tagService: tagService,
	}
}

func (t *Tag) GetAllTag(c *gin.Context) {

	uid := c.Param(httpserver.UserUidKey)

	tags, err := t.tagService.GetAllTag(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	tagsPublic := []model.TagPublic{}

	for _, item := range tags {
		tagsPublic = append(tagsPublic, model.ToTagPublic(item, true))
	}

	c.JSON(http.StatusOK, model.BulkResponse[model.TagPublic]{
		Count:    len(tagsPublic),
		Previous: nil,
		Next:     nil,
		Results:  tagsPublic,
	})

}

func (t *Tag) CreateTag(c *gin.Context) {

	var tagRequest model.CreateTagRequest

	err := c.ShouldBind(&tagRequest)
	if err != nil {
		logrus.Warn("Requets Body binding error: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	err = t.validate.Struct(tagRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tag := tagRequest.ToTag(c.Param(httpserver.UserUidKey))

	err = t.tagService.CreateTag(&tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}

	tagPublic := model.ToTagPublic(tag, false)
	c.JSON(http.StatusOK, tagPublic)
}

func (t *Tag) UpdateTag(c *gin.Context) {
	tagIdString, ok := c.Params.Get("tagId")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	// Check if the tagId is int parsable
	tagId, err := strconv.Atoi(tagIdString)
	if err != nil {
		logrus.Warn("Can't convert tagId to int: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	// Bind the request body.
	var updateTagRequest model.UpdateTagRequest
	err = c.ShouldBind(&updateTagRequest)
	if err != nil {
		logrus.Warn("Can't bind request body to model: ", err)
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	err = t.validate.Struct(updateTagRequest)
	if err != nil {
		logrus.Warn("Struct validation failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tag, err := t.tagService.UpdateTag(tagId, &updateTagRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something happen in the server"})
		return
	}

	tagPublic := model.ToTagPublic(*tag, true)
	c.JSON(http.StatusOK, tagPublic)
}

func (t *Tag) GetTagById(c *gin.Context) {
	tagIdStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}

	publicTag, err := t.tagService.GetTagById(tagId, true)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "tag id not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrGinBadRequestBody)
		return
	}

	c.JSON(http.StatusOK, publicTag)
}

func (t *Tag) DeleteTag(c *gin.Context) {
	tagIdString, ok := c.Params.Get("tagId")
	if !ok {
		c.JSON(http.StatusBadRequest, ErrGinBadRequestBody)
		return
	}
	tagId, err := strconv.Atoi(tagIdString)
	err = t.tagService.DeleteTag(tagId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
