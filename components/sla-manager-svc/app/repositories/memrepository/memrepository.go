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

/*
Package memrepository is a simple implementation of a model.IRepository intended for developing purposes.
*/
package memrepository

import (
	"fmt"
	"time"

	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
)

// path used in logs
const pathLOG string = "SLA > Repository > Memory >  "

// MemRepository is a repository in memory
type MemRepository struct {
	agreements map[string]model.SLA
	violations map[string]model.Violation
}

// NewMemRepository creates a MemRepository with an initial state set by the parameters
func NewMemRepository(agreements map[string]model.SLA, violations map[string]model.Violation) MemRepository {
	var r MemRepository

	if agreements == nil {
		agreements = make(map[string]model.SLA)
	}
	if violations == nil {
		violations = make(map[string]model.Violation)
	}

	r = MemRepository{
		agreements: agreements,
		violations: violations,
	}
	return r
}

// New creates a new instance of MemRepository
func New() (MemRepository, error) {
	return NewMemRepository(nil, nil), nil
}

///////////////////////////////////////////////////////////////////////////////

/*
GetAllQoSDefinitions returns the list of QoSDefinitions.

The list is empty when there are no QoSDefinitions;
error != nil on error
*/
func (r MemRepository) GetSLAs() (model.SLAs, error) {
	result := make(model.SLAs, 0, len(r.agreements))

	for _, value := range r.agreements {
		result = append(result, value)
	}
	return result, nil
}

// GetSLAsByName gets SLAs by Name.
func (r MemRepository) GetSLAsByName(id string) (model.SLAs, error) {
	result := make(model.SLAs, 0, len(r.agreements))

	for _, value := range r.agreements {
		if value.Name == id {
			result = append(result, value)
		}
	}
	return result, nil
}

/*
GetQoSDefinitionsByState returns the QoSDefinitions that match any of the items in states.

error != nil on error
*/
func (r MemRepository) GetSLAsByState(states ...model.State) (model.SLAs, error) {
	result := make(model.SLAs, 0)

	for _, a := range r.agreements {
		for _, state := range states {
			if a.State == state {
				result = append(result, a)
			}
		}
	}
	return result, nil
}

/*
GetQoSDefinition returns the QoSDefinition identified by id.

error != nil on error;
error is sql.ErrNoRows if the QoSDefinition is not found
*/
func (r MemRepository) GetSLA(id string) (*model.SLA, error) {
	var err error

	item, ok := r.agreements[id]

	if ok {
		err = nil
	} else {
		err = model.ErrNotFound
	}
	return &item, err
}

/*
CreateQoSDefinition stores a new QoSDefinition.

error != nil on error;
error is sql.ErrNoRows if the QoSDefinition already exists
*/
func (r MemRepository) CreateSLA(agreement *model.SLA) (*model.SLA, error) {
	logs.GetLogger().Info(pathLOG + "[CreateAgreement] Adding NEW agreement to default (memory) repository ...")
	agreementstr := fmt.Sprintf("%#v", agreement)
	logs.GetLogger().Debug(pathLOG + "[CreateAgreement] Agreement: " + agreementstr)

	var err error

	id := agreement.Id
	_, ok := r.agreements[id]

	if ok {
		err = model.ErrAlreadyExist
	} else {
		agreement.Creation = time.Now()

		//var texp *time.Time
		texp := new(time.Time)
		*texp = time.Now().AddDate(1, 0, 0)

		agreement.Expiration = texp

		r.agreements[id] = *agreement
	}
	return agreement, err
}

/*
UpdateQoSDefinition updates the information of an already saved instance of a QoSDefinition
*/
func (r MemRepository) UpdateSLA(agreement *model.SLA) (*model.SLA, error) {
	var err error

	id := agreement.Id
	_, ok := r.agreements[id]

	if !ok {
		err = model.ErrNotFound
	} else {
		r.agreements[id] = *agreement
	}
	return agreement, err
}

/*
DeleteQoSDefinition deletes from the repository the QoSDefinition whose id is provider.Id.

error != nil on error;
error is sql.ErrNoRows if the Agreement does not exist.
*/
func (r MemRepository) DeleteSLA(id string) error {
	var err error

	_, ok := r.agreements[id]
	if ok {
		delete(r.agreements, id)
	} else {
		err = model.ErrNotFound
	}
	return err
}

/*
CreateViolation stores a new Violation.

error != nil on error;
error is sql.ErrNoRows if the Violation already exists
*/
func (r MemRepository) CreateViolation(v *model.Violation) (*model.Violation, error) {
	logs.GetLogger().Info(pathLOG + "[CreateViolation] Adding Violation to default (memory) repository ...")
	vstr := fmt.Sprintf("%#v", v)
	logs.GetLogger().Debug(pathLOG + "[CreateViolation] Agreement: " + vstr)

	var err error

	id := v.Id

	if _, ok := r.violations[id]; ok {
		err = model.ErrAlreadyExist
	} else {
		r.violations[id] = *v
	}
	return v, err
}

/*
GetViolation returns the Violation identified by id.

error != nil on error;
error is sql.ErrNoRows if the Violation is not found
*/
func (r MemRepository) GetViolation(id string) (*model.Violation, error) {
	var err error

	item, ok := r.violations[id]

	if ok {
		err = nil
	} else {
		err = model.ErrNotFound
	}
	return &item, err
}

/*
GetViolations returns the Violations of an SLA.

The list is empty when there are no violations;
error != nil on error
*/
func (r MemRepository) GetViolations(id string) (model.Violations, error) {
	result := make(model.Violations, 0, len(r.violations))

	for _, value := range r.violations {
		if value.AgreementId == id {
			result = append(result, value)
		}
	}
	return result, nil
}

/*
GetAppViolations returns the Violations of an application.

The list is empty when there are no violations;
error != nil on error
*/
func (r MemRepository) GetAppViolations(id string) (model.Violations, error) {
	result := make(model.Violations, 0, len(r.violations))

	for _, value := range r.violations {
		if value.AppId == id {
			result = append(result, value)
		}
	}
	return result, nil
}

/*
GetAllViolations returns the list of violations.

The list is empty when there are no violations;
error != nil on error
*/
func (r MemRepository) GetAllViolations() (model.Violations, error) {
	result := make(model.Violations, 0, len(r.violations))

	for _, value := range r.violations {
		result = append(result, value)
	}
	return result, nil
}

/*
UpdateQoSDefinitionState transits the state of the QoSDefinition
*/
func (r MemRepository) UpdateSLAState(id string, newState model.State) (*model.SLA, error) {

	var ok bool
	var err error
	var current model.SLA
	var result *model.SLA

	current, ok = r.agreements[id]

	if !ok {
		err = model.ErrNotFound
	} else {
		current.State = newState
		r.agreements[id] = current
		result = &current
	}
	return result, err
}
