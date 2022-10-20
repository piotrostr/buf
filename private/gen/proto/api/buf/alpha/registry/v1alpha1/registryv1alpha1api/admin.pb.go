// Copyright 2020-2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-api. DO NOT EDIT.

package registryv1alpha1api

import (
	context "context"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
)

// AdminService is the Admin service.
type AdminService interface {
	// ForceDeleteUser forces to delete a user. Resources and organizations that are
	// solely owned by the user will also be deleted.
	ForceDeleteUser(
		ctx context.Context,
		userId string,
	) (user *v1alpha1.User, organizations []*v1alpha1.Organization, repositories []*v1alpha1.Repository, plugins []*v1alpha1.Plugin, templates []*v1alpha1.Template, err error)
	// Update a user's verification status
	UpdateUserVerificationStatus(
		ctx context.Context,
		userId string,
		verificationStatus v1alpha1.VerificationStatus,
	) (err error)
	// Update a organization's verification
	UpdateOrganizationVerificationStatus(
		ctx context.Context,
		organizationId string,
		verificationStatus v1alpha1.VerificationStatus,
	) (err error)
}
