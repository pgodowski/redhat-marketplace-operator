package reconcileutils

// Code generated by github.com/launchdarkly/go-options.  DO NOT EDIT.

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplyDeleteActionOptionFunc func(c *deleteActionOptions) error

func (f ApplyDeleteActionOptionFunc) apply(c *deleteActionOptions) error {
	return f(c)
}

func newDeleteActionOptions(options ...DeleteActionOption) (deleteActionOptions, error) {
	var c deleteActionOptions
	err := applyDeleteActionOptionsOptions(&c, options...)
	return c, err
}

func applyDeleteActionOptionsOptions(c *deleteActionOptions, options ...DeleteActionOption) error {
	for _, o := range options {
		if err := o.apply(c); err != nil {
			return err
		}
	}
	return nil
}

type DeleteActionOption interface {
	apply(*deleteActionOptions) error
}

func DeleteWithStatusCondition(o UpdateStatusConditionFunc) ApplyDeleteActionOptionFunc {
	return func(c *deleteActionOptions) error {
		c.WithStatusCondition = o
		return nil
	}
}

func DeleteWithDeleteOptions(o ...client.DeleteOption) ApplyDeleteActionOptionFunc {
	return func(c *deleteActionOptions) error {
		c.WithDeleteOptions = o
		return nil
	}
}
