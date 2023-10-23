// Code generated by mockery v1.0.0. DO NOT EDIT.

package arimocks

import (
	ari "github.com/PolyAI-LDN/ari/v6"
	mock "github.com/stretchr/testify/mock"
)

// DTMFSender is an autogenerated mock type for the DTMFSender type
type DTMFSender struct {
	mock.Mock
}

// SendDTMF provides a mock function with given fields: dtmf, opts
func (_m *DTMFSender) SendDTMF(dtmf string, opts *ari.DTMFOptions) {
	_m.Called(dtmf, opts)
}
