package gomeassistant

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"saml.dev/gome-assistant/internal"
)

type MockState struct {
	EqualsReturn bool
	EqualsError  bool
	GetReturn    EntityState
	GetError     bool
}

func (s MockState) AfterSunrise(_ ...DurationString) bool {
	return true
}
func (s MockState) BeforeSunrise(_ ...DurationString) bool {
	return true
}
func (s MockState) AfterSunset(_ ...DurationString) bool {
	return true
}
func (s MockState) BeforeSunset(_ ...DurationString) bool {
	return true
}
func (s MockState) Get(eid string) (EntityState, error) {
	if s.GetError {
		return EntityState{}, errors.New("some error")
	}
	return s.GetReturn, nil
}
func (s MockState) List() ([]EntityState, error) {
	return []EntityState{}, nil
}
func (s MockState) Equals(eid, state string) (bool, error) {
	if s.EqualsError {
		return false, errors.New("some error")
	}
	return s.EqualsReturn, nil
}

var runOnError = internal.EnabledDisabledInfo{
	Entity:     "eid",
	State:      "state",
	RunOnError: true,
}

var dontRunOnError = internal.EnabledDisabledInfo{
	Entity:     "eid",
	State:      "state",
	RunOnError: false,
}

func list(infos ...internal.EnabledDisabledInfo) []internal.EnabledDisabledInfo {
	ret := []internal.EnabledDisabledInfo{}
	ret = append(ret, infos...)
	return ret
}

func TestEnabledEntity_StateEqual_Passes(t *testing.T) {
	state := MockState{
		EqualsReturn: true,
	}
	c := checkEnabledEntity(state, list(runOnError))
	assert.False(t, c.fail, "should pass")
}

func TestEnabledEntity_StateNotEqual_Fails(t *testing.T) {
	state := MockState{
		EqualsReturn: false,
	}
	c := checkEnabledEntity(state, list(runOnError))
	assert.True(t, c.fail, "should fail")
}

func TestEnabledEntity_NetworkError_DontRun_Fails(t *testing.T) {
	state := MockState{
		EqualsError: true,
	}
	c := checkEnabledEntity(state, list(dontRunOnError))
	assert.True(t, c.fail, "should fail")
}

func TestEnabledEntity_NetworkError_StillRun_Passes(t *testing.T) {
	state := MockState{
		EqualsError: true,
	}
	c := checkEnabledEntity(state, list(runOnError))
	assert.False(t, c.fail, "should fail")
}

func TestDisabledEntity_StateEqual_Fails(t *testing.T) {
	state := MockState{
		EqualsReturn: true,
	}
	c := checkDisabledEntity(state, list(runOnError))
	assert.True(t, c.fail, "should pass")
}

func TestDisabledEntity_StateNotEqual_Passes(t *testing.T) {
	state := MockState{
		EqualsReturn: false,
	}
	c := checkDisabledEntity(state, list(runOnError))
	assert.False(t, c.fail, "should fail")
}

func TestDisabledEntity_NetworkError_DontRun_Fails(t *testing.T) {
	state := MockState{
		EqualsError: true,
	}
	c := checkDisabledEntity(state, list(dontRunOnError))
	assert.True(t, c.fail, "should fail")
}

func TestDisabledEntity_NetworkError_StillRun_Passes(t *testing.T) {
	state := MockState{
		EqualsError: true,
	}
	c := checkDisabledEntity(state, list(runOnError))
	assert.False(t, c.fail, "should fail")
}

func TestStatesMatch(t *testing.T) {
	c := checkStatesMatch("hey", "hey")
	assert.False(t, c.fail, "should pass")
}

func TestStatesDontMatch(t *testing.T) {
	c := checkStatesMatch("hey", "bye")
	assert.True(t, c.fail, "should fail")
}
