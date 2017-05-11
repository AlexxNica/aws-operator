package aws

import (
	"fmt"

	"github.com/juju/errgo"
)

var (
	notImplementedMethodError = errgo.New("not implemented")

	noBucketInBucketObjectError = errgo.New("Object needs to belong to some bucket")

	routeTableFindError    = errgo.New("Couldn't find route table")
	securityGroupFindError = errgo.New("Couldn't find security group")
	subnetFindError        = errgo.New("Couldn't find subnet")
	vpcFindError           = errgo.New("Couldn't find VPC")
	instanceFindError      = errgo.New("Couldn't find EC2 instance")

	resourceDeleteError       = errgo.New("Couldn't delete resource, it lacks the necessary data (ID)")
	clientNotInitializedError = errgo.New("The client has not been initialized")

	kmsKeyAliasEmptyError = errgo.New("the KMS key alias cannot be empty")
)

type DomainNamedResourceNotFoundError struct {
	Domain string
}

func (e DomainNamedResourceNotFoundError) Error() string {
	return fmt.Sprintf("No Hosted Zones found for domain %s", e.Domain)
}

type NamedResourceNotFoundError struct {
	Name string
}

func (e NamedResourceNotFoundError) Error() string {
	return fmt.Sprintf("The resource was not found: %s", e.Name)
}

var gatewayNotFoundError = errgo.New("gateway not found")

// IsGatewayNotFound asserts gatewayNotFoundError.
func IsGatewayNotFound(err error) bool {
	return errgo.Cause(err) == gatewayNotFoundError
}

func IsInstanceFindError(err error) bool {
	return errgo.Cause(err) == instanceFindError
}

var routeNotFoundError = errgo.New("route not found")

// IsRouteNotFoundError asserts routeNotFoundError.
func IsRouteNotFoundError(err error) bool {
	return errgo.Cause(err) == routeNotFoundError
}

func IsSecurityGroupFind(err error) bool {
	return errgo.Cause(err) == securityGroupFindError
}

func IsSubnetFind(err error) bool {
	return errgo.Cause(err) == subnetFindError
}

func IsVpcFindError(err error) bool {
	return errgo.Cause(err) == vpcFindError
}

func IsKMSKeyAliasEmpty(err error) bool {
	return errgo.Cause(err) == kmsKeyAliasEmptyError
}
