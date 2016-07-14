package test

import (
	"testing"

	"github.com/coralproject/pillar/pkg/model"
)

func TestUser(t *testing.T) {
	u1 := model.User{Name: "", Avatar: ""}
	//no error is a problem for incomplete models
	if err := u1.Validate(); err == nil {
		t.Fail()
	}

	u2 := model.User{Name: "", Avatar: "coralavatar"}
	//no error is a problem for incomplete models
	if err := u2.Validate(); err == nil {
		t.Fail()
	}
}

func TestComment(t *testing.T) {
	c1 := model.Comment{Body: ""}
	//no error is a problem for incomplete models
	if err := c1.Validate(); err == nil {
		t.Fail()
	}
}

func TestAsset(t *testing.T) {
	c1 := model.Asset{URL: ""}
	//no error is a problem for incomplete models
	if err := c1.Validate(); err == nil {
		t.Fail()
	}
}
