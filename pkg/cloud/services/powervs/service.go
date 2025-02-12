/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package powervs

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_images"
	"github.com/IBM-Cloud/power-go-client/power/models"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
)

var _ PowerVS = &Service{}

// Service holds the PowerVS Service specific information.
type Service struct {
	session        *ibmpisession.IBMPISession
	instanceClient *instance.IBMPIInstanceClient
	networkClient  *instance.IBMPINetworkClient
	imageClient    *instance.IBMPIImageClient
	jobClient      *instance.IBMPIJobClient
}

// ServiceOptions holds the PowerVS Service Options specific information.
type ServiceOptions struct {
	*ibmpisession.IBMPIOptions

	CloudInstanceID string
}

// CreateInstance creates the virtual machine in the Power VS service instance.
func (s *Service) CreateInstance(body *models.PVMInstanceCreate) (*models.PVMInstanceList, error) {
	return s.instanceClient.Create(body)
}

// DeleteInstance deletes the virtual machine in the Power VS service instance.
func (s *Service) DeleteInstance(id string) error {
	return s.instanceClient.Delete(id)
}

// GetAllInstance returns all the virtual machine in the Power VS service instance.
func (s *Service) GetAllInstance() (*models.PVMInstances, error) {
	return s.instanceClient.GetAll()
}

// GetInstance returns the virtual machine in the Power VS service instance.
func (s *Service) GetInstance(id string) (*models.PVMInstance, error) {
	return s.instanceClient.Get(id)
}

// GetImage returns the image in the Power VS service instance.
func (s *Service) GetImage(id string) (*models.Image, error) {
	return s.imageClient.Get(id)
}

// GetAllImage returns all the images in the Power VS service instance.
func (s *Service) GetAllImage() (*models.Images, error) {
	return s.imageClient.GetAll()
}

// DeleteImage deletes the image in the Power VS service instance.
func (s *Service) DeleteImage(id string) error {
	return s.imageClient.Delete(id)
}

// CreateCosImage creates a import job to import the image in the Power VS service instance.
func (s *Service) CreateCosImage(body *models.CreateCosImageImportJob) (*models.JobReference, error) {
	return s.imageClient.CreateCosImage(body)
}

// GetCosImages returns the last import job in the Power VS service instance.
func (s *Service) GetCosImages(id string) (*models.Job, error) {
	params := p_cloud_images.NewPcloudV1CloudinstancesCosimagesGetParams().WithCloudInstanceID(id)
	resp, err := s.session.Power.PCloudImages.PcloudV1CloudinstancesCosimagesGet(params, s.session.AuthInfo(id))
	if err != nil || resp.Payload == nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GetJob returns the import job to in the Power VS service instance.
func (s *Service) GetJob(id string) (*models.Job, error) {
	return s.jobClient.Get(id)
}

// DeleteJob deletes the image import job in the Power VS service instance.
func (s *Service) DeleteJob(id string) error {
	return s.jobClient.Delete(id)
}

// GetAllNetwork returns all the networks in the Power VS service instance.
func (s *Service) GetAllNetwork() (*models.Networks, error) {
	return s.networkClient.GetAll()
}

// NewService returns a new service for the Power VS api client.
func NewService(options ServiceOptions) (PowerVS, error) {
	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, err
	}
	options.Authenticator = auth
	session, err := ibmpisession.NewIBMPISession(options.IBMPIOptions)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	return &Service{
		session:        session,
		instanceClient: instance.NewIBMPIInstanceClient(ctx, session, options.CloudInstanceID),
		networkClient:  instance.NewIBMPINetworkClient(ctx, session, options.CloudInstanceID),
		imageClient:    instance.NewIBMPIImageClient(ctx, session, options.CloudInstanceID),
		jobClient:      instance.NewIBMPIJobClient(ctx, session, options.CloudInstanceID),
	}, nil
}
