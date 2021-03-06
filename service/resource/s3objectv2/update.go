package s3objectv2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	updateBucketState, err := toBucketObjectState(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for key, bucketObject := range updateBucketState {
		if bucketObject.Key != "" {
			s3PutInput, err := toPutObjectInput(bucketObject)
			if err != nil {
				return microerror.Mask(err)
			}

			_, err = r.awsClients.S3.PutObject(&s3PutInput)
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.LogCtx(ctx, "debug", fmt.Sprintf("updating S3 object '%s': updated", key))
		} else {
			r.logger.LogCtx(ctx, "debug", fmt.Sprintf("updating S3 object '%s': already updated", key))
		}
	}

	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	create, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()
	patch.SetCreateChange(create)
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentBucketState, err := toBucketObjectState(currentState)
	if err != nil {
		return s3.PutObjectInput{}, microerror.Mask(err)
	}

	desiredBucketState, err := toBucketObjectState(desiredState)
	if err != nil {
		return s3.PutObjectInput{}, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "finding out if the s3 objects should be updated")

	updateState := map[string]BucketObjectState{}

	for key, bucketObject := range desiredBucketState {
		if _, ok := currentBucketState[key]; !ok {
			updateState[key] = BucketObjectState{}
		}

		currentObject := currentBucketState[key]
		if currentObject.Body != "" && bucketObject.Body != currentObject.Body {
			updateState[key] = bucketObject
		} else {
			updateState[key] = BucketObjectState{}
		}
	}

	return updateState, nil
}
