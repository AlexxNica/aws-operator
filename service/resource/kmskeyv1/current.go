package kmskeyv1

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/aws-operator/service/keyv1"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	currentState := KMSKeyState{}

	customObject, err := keyv1.ToCustomObject(obj)
	if err != nil {
		return currentState, err
	}

	clusterID := keyv1.ClusterID(customObject)
	alias := toAlias(clusterID)
	input := &kms.DescribeKeyInput{
		KeyId: aws.String(alias),
	}

	output, err := r.awsClients.KMS.DescribeKey(input)
	if err != nil {
		return currentState, microerror.Mask(err)
	}

	currentState.KeyID = *output.KeyMetadata.KeyId
	currentState.KeyAlias = alias

	return currentState, nil
}