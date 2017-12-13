package cloudformationv2

import (
	"context"
	"testing"

	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/aws-operator/service/resource/cloudformationv2/adapter"
)

func Test_Resource_Cloudformation_newCreate(t *testing.T) {
	clusterTpo := &v1alpha1.AWSConfig{
		Spec: v1alpha1.AWSConfigSpec{
			Cluster: v1alpha1.Cluster{
				ID: "test-cluster",
				Kubernetes: v1alpha1.ClusterKubernetes{
					IngressController: v1alpha1.ClusterKubernetesIngressController{
						Domain: "mysubdomain.mydomain.com",
					},
				},
			},
			AWS: v1alpha1.AWSConfigSpecAWS{
				AZ: "myaz",
				Workers: []v1alpha1.AWSConfigSpecAWSNode{
					v1alpha1.AWSConfigSpecAWSNode{
						ImageID: "myimageid",
					},
				},
			},
		},
	}

	testCases := []struct {
		obj               interface{}
		currentState      interface{}
		desiredState      interface{}
		expectedStackName string
		description       string
	}{
		{
			description:       "current and desired state empty, expected empty",
			obj:               clusterTpo,
			currentState:      StackState{},
			desiredState:      StackState{},
			expectedStackName: "",
		},
		{
			description:  "current state empty, desired state not empty, expected desired state",
			obj:          clusterTpo,
			currentState: StackState{},
			desiredState: StackState{
				Name: "desired",
			},
			expectedStackName: "desired",
		},
		{
			description: "current state not empty, desired state not empty but different, expected desired state",
			obj:         clusterTpo,
			currentState: StackState{
				Name: "current",
			},
			desiredState: StackState{
				Name: "desired",
			},
			expectedStackName: "desired",
		},
	}

	var err error
	var newResource *Resource
	{
		resourceConfig := DefaultConfig()
		resourceConfig.Clients = adapter.Clients{
			EC2: &adapter.EC2ClientMock{},
			IAM: &adapter.IAMClientMock{},
			KMS: &adapter.KMSClientMock{},
		}
		resourceConfig.Logger = microloggertest.New()
		newResource, err = New(resourceConfig)
		if err != nil {
			t.Error("expected", nil, "got", err)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := newResource.newCreateChange(context.TODO(), tc.obj, tc.currentState, tc.desiredState)
			if err != nil {
				t.Errorf("expected '%v' got '%#v'", nil, err)
			}
			createChange, ok := result.(awscloudformation.CreateStackInput)
			if !ok {
				t.Errorf("expected '%T', got '%T'", createChange, result)
			}
			if createChange.StackName != nil && *createChange.StackName != tc.expectedStackName {
				t.Errorf("expected %s, got %s", tc.expectedStackName, createChange.StackName)
			}
		})
	}
}
