/*
  COLMENA-DESCRIPTION-SERVICE
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
package model

// IRepository expose the interface to be fulfilled by implementations of repositories.
type IRepository interface {

	/*
	 * GetSLAs returns the list of SLAs.
	 * The list is empty when there are no QoSDefinitions;
	 * error != nil on error
	 */
	GetSLAs() (SLAs, error)

	/*
	 * GetSLAsByName returns the list of SLAs.
	 * The list is empty when there are no QoSDefinitions;
	 * error != nil on error
	 */
	GetSLAsByName(id string) (SLAs, error)

	/*
	 * GetSLA returns the SLA identified by id.
	 * error != nil on error;
	 * error is sql.ErrNoRows if the QoSDefinition is not found
	 */
	GetSLA(id string) (*SLA, error)

	/*
	 * GetSLAsByState returns the SLAs that have one of the items in states.
	 * error != nil on error;
	 */
	GetSLAsByState(states ...State) (SLAs, error)

	/*
	 * CreateSLA stores a new SLA definitipon.
	 */
	CreateSLA(qos *SLA) (*SLA, error)

	/*
	 *UpdateSLA updates the information of an already saved instance of an agreement
	 */
	UpdateSLA(qos *SLA) (*SLA, error)

	/*
	 * DeleteSLA deletes from the repository the SLAs by id
	 */
	DeleteSLA(id string) error

	/*
	 * CreateViolation stores a new Violation.
	 * error != nil on error;
	 * error is sql.ErrNoRows if the Violation already exists
	 */
	CreateViolation(v *Violation) (*Violation, error)

	/*
	 * GetViolation returns the Violation identified by id.
	 * error != nil on error;
	 * error is sql.ErrNoRows if the Violation is not found
	 */
	GetViolation(id string) (*Violation, error)

	/*
	 * GetViolations returns the list of violations of an SLA.
	 * The list is empty when there are no violations;
	 * error != nil on error
	 */
	GetViolations(id string) (Violations, error)

	/*
	 * GetAppViolations returns the list of violations of an application.
	 * The list is empty when there are no violations;
	 * error != nil on error
	 */
	GetAppViolations(id string) (Violations, error)

	/*
	 * GetAllViolations returns the list of violations.
	 * The list is empty when there are no violations;
	 * error != nil on error
	 */
	GetAllViolations() (Violations, error)

	/*
	 * UpdateSLAState changes the state of an SLA.
	 * Returns the updated SLA; error != nil on error
	 * error is sql.ErrNoRows if the SLA does not exist
	 * Non-sentinel error is returned if not a valid transition
	 * (it is recommended to check a.IsValidTransition before UpdateSLAState)
	 */
	UpdateSLAState(id string, newState State) (*SLA, error)
}
