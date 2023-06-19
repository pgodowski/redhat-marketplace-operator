//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020 IBM Co..

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIKeyRef) DeepCopyInto(out *APIKeyRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIKeyRef.
func (in *APIKeyRef) DeepCopy() *APIKeyRef {
	if in == nil {
		return nil
	}
	out := new(APIKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccesKeyIDRef) DeepCopyInto(out *AccesKeyIDRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccesKeyIDRef.
func (in *AccesKeyIDRef) DeepCopy() *AccesKeyIDRef {
	if in == nil {
		return nil
	}
	out := new(AccesKeyIDRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAuth) DeepCopyInto(out *ClusterAuth) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAuth.
func (in *ClusterAuth) DeepCopy() *ClusterAuth {
	if in == nil {
		return nil
	}
	out := new(ClusterAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMapRef) DeepCopyInto(out *ConfigMapRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMapRef.
func (in *ConfigMapRef) DeepCopy() *ConfigMapRef {
	if in == nil {
		return nil
	}
	out := new(ConfigMapRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericMapRef) DeepCopyInto(out *GenericMapRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericMapRef.
func (in *GenericMapRef) DeepCopy() *GenericMapRef {
	if in == nil {
		return nil
	}
	out := new(GenericMapRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Git) DeepCopyInto(out *Git) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Git.
func (in *Git) DeepCopy() *Git {
	if in == nil {
		return nil
	}
	out := new(Git)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Header) DeepCopyInto(out *Header) {
	{
		in := &in
		*out = make(Header, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Header.
func (in Header) DeepCopy() Header {
	if in == nil {
		return nil
	}
	out := new(Header)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Headers) DeepCopyInto(out *Headers) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string]Header, len(*in))
		for key, val := range *in {
			var outVal map[string]string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(Header, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Headers.
func (in *Headers) DeepCopy() *Headers {
	if in == nil {
		return nil
	}
	out := new(Headers)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeadersFrom) DeepCopyInto(out *HeadersFrom) {
	*out = *in
	out.ConfigMapRef = in.ConfigMapRef
	out.SecretMapRef = in.SecretMapRef
	out.GenericMapRef = in.GenericMapRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeadersFrom.
func (in *HeadersFrom) DeepCopy() *HeadersFrom {
	if in == nil {
		return nil
	}
	out := new(HeadersFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResource) DeepCopyInto(out *RemoteResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResource.
func (in *RemoteResource) DeepCopy() *RemoteResource {
	if in == nil {
		return nil
	}
	out := new(RemoteResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RemoteResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceAPIKeyRef) DeepCopyInto(out *RemoteResourceAPIKeyRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceAPIKeyRef.
func (in *RemoteResourceAPIKeyRef) DeepCopy() *RemoteResourceAPIKeyRef {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceAPIKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceAccesKeyIDRef) DeepCopyInto(out *RemoteResourceAccesKeyIDRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceAccesKeyIDRef.
func (in *RemoteResourceAccesKeyIDRef) DeepCopy() *RemoteResourceAccesKeyIDRef {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceAccesKeyIDRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceAuth) DeepCopyInto(out *RemoteResourceAuth) {
	*out = *in
	if in.Hmac != nil {
		in, out := &in.Hmac, &out.Hmac
		*out = new(RemoteResourceHmac)
		(*in).DeepCopyInto(*out)
	}
	if in.Iam != nil {
		in, out := &in.Iam, &out.Iam
		*out = new(RemoteResourceIam)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceAuth.
func (in *RemoteResourceAuth) DeepCopy() *RemoteResourceAuth {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceHmac) DeepCopyInto(out *RemoteResourceHmac) {
	*out = *in
	in.AccesKeyIDRef.DeepCopyInto(&out.AccesKeyIDRef)
	in.SecretAccessKeyRef.DeepCopyInto(&out.SecretAccessKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceHmac.
func (in *RemoteResourceHmac) DeepCopy() *RemoteResourceHmac {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceHmac)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceIam) DeepCopyInto(out *RemoteResourceIam) {
	*out = *in
	in.APIKeyRef.DeepCopyInto(&out.APIKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceIam.
func (in *RemoteResourceIam) DeepCopy() *RemoteResourceIam {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceIam)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceList) DeepCopyInto(out *RemoteResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RemoteResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceList.
func (in *RemoteResourceList) DeepCopy() *RemoteResourceList {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RemoteResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceRequest) DeepCopyInto(out *RemoteResourceRequest) {
	*out = *in
	out.Git = in.Git
	in.Headers.DeepCopyInto(&out.Headers)
	out.HeadersFrom = in.HeadersFrom
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceRequest.
func (in *RemoteResourceRequest) DeepCopy() *RemoteResourceRequest {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceSecretAccessKeyRef) DeepCopyInto(out *RemoteResourceSecretAccessKeyRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceSecretAccessKeyRef.
func (in *RemoteResourceSecretAccessKeyRef) DeepCopy() *RemoteResourceSecretAccessKeyRef {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceSecretAccessKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceSpec) DeepCopyInto(out *RemoteResourceSpec) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
	out.ClusterAuth = in.ClusterAuth
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make([]Request, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceSpec.
func (in *RemoteResourceSpec) DeepCopy() *RemoteResourceSpec {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceStatus) DeepCopyInto(out *RemoteResourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceStatus.
func (in *RemoteResourceStatus) DeepCopy() *RemoteResourceStatus {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceValueFrom) DeepCopyInto(out *RemoteResourceValueFrom) {
	*out = *in
	in.SecretKeyRef.DeepCopyInto(&out.SecretKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceValueFrom.
func (in *RemoteResourceValueFrom) DeepCopy() *RemoteResourceValueFrom {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceValueFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Request) DeepCopyInto(out *Request) {
	*out = *in
	in.Options.DeepCopyInto(&out.Options)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Request.
func (in *Request) DeepCopy() *Request {
	if in == nil {
		return nil
	}
	out := new(Request)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3Options) DeepCopyInto(out *S3Options) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string]Header, len(*in))
		for key, val := range *in {
			var outVal map[string]string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(Header, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3Options.
func (in *S3Options) DeepCopy() *S3Options {
	if in == nil {
		return nil
	}
	out := new(S3Options)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretAccessKeyRef) DeepCopyInto(out *SecretAccessKeyRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretAccessKeyRef.
func (in *SecretAccessKeyRef) DeepCopy() *SecretAccessKeyRef {
	if in == nil {
		return nil
	}
	out := new(SecretAccessKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretMapRef) DeepCopyInto(out *SecretMapRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretMapRef.
func (in *SecretMapRef) DeepCopy() *SecretMapRef {
	if in == nil {
		return nil
	}
	out := new(SecretMapRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *URL) DeepCopyInto(out *URL) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new URL.
func (in *URL) DeepCopy() *URL {
	if in == nil {
		return nil
	}
	out := new(URL)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueFrom) DeepCopyInto(out *ValueFrom) {
	*out = *in
	in.SecretKeyRef.DeepCopyInto(&out.SecretKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueFrom.
func (in *ValueFrom) DeepCopy() *ValueFrom {
	if in == nil {
		return nil
	}
	out := new(ValueFrom)
	in.DeepCopyInto(out)
	return out
}
