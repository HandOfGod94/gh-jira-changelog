// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	jira "github.com/handofgod94/gh-jira-changelog/pkg/jira_changelog/jira"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// FetchIssue provides a mock function with given fields: issueId
func (_m *Client) FetchIssue(issueId string) (jira.Issue, error) {
	ret := _m.Called(issueId)

	var r0 jira.Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (jira.Issue, error)); ok {
		return rf(issueId)
	}
	if rf, ok := ret.Get(0).(func(string) jira.Issue); ok {
		r0 = rf(issueId)
	} else {
		r0 = ret.Get(0).(jira.Issue)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(issueId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}