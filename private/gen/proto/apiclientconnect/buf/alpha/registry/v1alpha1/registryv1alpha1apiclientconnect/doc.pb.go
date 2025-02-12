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

// Code generated by protoc-gen-go-apiclientconnect. DO NOT EDIT.

package registryv1alpha1apiclientconnect

import (
	context "context"
	registryv1alpha1connect "github.com/bufbuild/buf/private/gen/proto/connect/buf/alpha/registry/v1alpha1/registryv1alpha1connect"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	connect_go "github.com/bufbuild/connect-go"
	zap "go.uber.org/zap"
)

type docServiceClient struct {
	logger *zap.Logger
	client registryv1alpha1connect.DocServiceClient
}

// GetSourceDirectoryInfo retrieves the directory and file structure for the
// given owner, repository and reference.
//
// The purpose of this is to get a representation of the file tree for a given
// module to enable exploring the module by navigating through its contents.
func (s *docServiceClient) GetSourceDirectoryInfo(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
) (root *v1alpha1.FileInfo, _ error) {
	response, err := s.client.GetSourceDirectoryInfo(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.GetSourceDirectoryInfoRequest{
				Owner:      owner,
				Repository: repository,
				Reference:  reference,
			}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg.Root, nil
}

// GetSourceFile retrieves the source contents for the given owner, repository,
// reference, and path.
func (s *docServiceClient) GetSourceFile(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
	path string,
) (content []byte, _ error) {
	response, err := s.client.GetSourceFile(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.GetSourceFileRequest{
				Owner:      owner,
				Repository: repository,
				Reference:  reference,
				Path:       path,
			}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg.Content, nil
}

// GetModulePackages retrieves the list of packages for the module based on the given
// owner, repository, and reference.
func (s *docServiceClient) GetModulePackages(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
) (name string, modulePackages []*v1alpha1.ModulePackage, _ error) {
	response, err := s.client.GetModulePackages(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.GetModulePackagesRequest{
				Owner:      owner,
				Repository: repository,
				Reference:  reference,
			}),
	)
	if err != nil {
		return "", nil, err
	}
	return response.Msg.Name, response.Msg.ModulePackages, nil
}

// GetModuleDocumentation retrieves the documentation for module based on the given
// owner, repository, and reference.
func (s *docServiceClient) GetModuleDocumentation(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
) (moduleDocumentation *v1alpha1.ModuleDocumentation, _ error) {
	response, err := s.client.GetModuleDocumentation(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.GetModuleDocumentationRequest{
				Owner:      owner,
				Repository: repository,
				Reference:  reference,
			}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg.ModuleDocumentation, nil
}

// GetPackageDocumentation retrieves a a slice of documentation structures
// for the given owner, repository, reference, and package name.
func (s *docServiceClient) GetPackageDocumentation(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
	packageName string,
) (packageDocumentation *v1alpha1.PackageDocumentation, _ error) {
	response, err := s.client.GetPackageDocumentation(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.GetPackageDocumentationRequest{
				Owner:       owner,
				Repository:  repository,
				Reference:   reference,
				PackageName: packageName,
			}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg.PackageDocumentation, nil
}
