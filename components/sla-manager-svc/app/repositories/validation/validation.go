/*
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
*/

// Usage:
//
//	repo, err := mongodb.New(config)
//	repo, _ = validation.New(repo)
package validation

import (
	"colmena/sla-management-svc/app/model"

	"bytes"
	"fmt"
)

const (
	fakeID = "_"
)

type repository struct {
	backend model.IRepository
	val     model.Validator
}

type valError struct {
	msg string
}

func (e *valError) Error() string {
	return e.msg
}

func newValError(errs []error) *valError {
	var buffer bytes.Buffer
	for _, err := range errs {
		buffer.WriteString(err.Error())
		buffer.WriteString(". ")
	}
	return &valError{msg: buffer.String()}
}

func (e *valError) IsErrValidation() bool {
	return true
}

// New returns an IRepository that performs validation before calling the actual repository.
func New(backend model.IRepository, val model.Validator) (model.IRepository, error) {
	return repository{
		backend: backend,
		val:     val,
	}, nil
}

// GetSLAs gets all SLAs.
func (r repository) GetSLAs() (model.SLAs, error) {
	return r.backend.GetSLAs()
}

// GetSLAsByName gets SLAs by Name.
func (r repository) GetSLAsByName(id string) (model.SLAs, error) {
	return r.backend.GetSLAsByName(id)
}

// GetSLA gets a SLA by id
func (r repository) GetSLA(id string) (*model.SLA, error) {
	return r.backend.GetSLA(id)
}

// GetSLAsByState returns the SLAs that have one of the items in states.
func (r repository) GetSLAsByState(states ...model.State) (model.SLAs, error) {
	return r.backend.GetSLAsByState(states...)
}

// CreateSLA validates and persists a SLA.
func (r repository) CreateSLA(agreement *model.SLA) (*model.SLA, error) {
	if errs := agreement.Validate(r.val, model.CREATE); len(errs) > 0 {
		err := newValError(errs)
		return agreement, err
	}
	return r.backend.CreateSLA(agreement)
}

// UpdateSLA validates and updates an SLA.
func (r repository) UpdateSLA(agreement *model.SLA) (*model.SLA, error) {

	/*
		It does not validate change of State.
	*/

	if errs := agreement.Validate(r.val, model.UPDATE); len(errs) > 0 {
		err := newValError(errs)
		return agreement, err
	}
	return r.backend.UpdateSLA(agreement)
}

// DeleteSLA deletes an QoSDefinition from repository.
func (r repository) DeleteSLA(id string) error {
	return r.backend.DeleteSLA(id)
}

// CreateViolation validates and persists a new Violation.
func (r repository) CreateViolation(v *model.Violation) (*model.Violation, error) {

	if errs := v.Validate(r.val, model.CREATE); len(errs) > 0 {
		err := newValError(errs)
		return v, err
	}
	return r.backend.CreateViolation(v)
}

// GetViolation returns the Violation identified by id.
func (r repository) GetViolation(id string) (*model.Violation, error) {
	return r.backend.GetViolation(id)
}

// GetViolations returns the Violations of an SLA.
func (r repository) GetViolations(id string) (model.Violations, error) {
	return r.backend.GetViolations(id)
}

// GetAppViolations returns the Violations of an application.
func (r repository) GetAppViolations(id string) (model.Violations, error) {
	return r.backend.GetAppViolations(id)
}

// GetAllViolations gets all agreements.
func (r repository) GetAllViolations() (model.Violations, error) {
	return r.backend.GetAllViolations()
}

// UpdateSLAState changes the state of a SLA.
func (r repository) UpdateSLAState(id string, newState model.State) (*model.SLA, error) {
	var err error
	newState = newState.Normalize()

	current, err := r.GetSLA(id)
	if err != nil {
		return nil, err
	}
	if !current.IsValidTransition(newState) {
		msg := fmt.Sprintf("Not valid transition from %s to %s for agreement %s",
			current.State, newState, id)
		err := &valError{msg: msg}
		return nil, err
	}
	return r.backend.UpdateSLAState(id, newState)
}
