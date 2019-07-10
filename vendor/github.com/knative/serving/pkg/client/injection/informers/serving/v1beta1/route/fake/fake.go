/*
Copyright 2019 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package fake

import (
	"context"

	fake "github.com/knative/serving/pkg/client/injection/informers/serving/factory/fake"
	route "github.com/knative/serving/pkg/client/injection/informers/serving/v1beta1/route"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
)

var Get = route.Get

func init() {
	injection.Fake.RegisterInformer(withInformer)
}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := fake.Get(ctx)
	inf := f.Serving().V1beta1().Routes()
	return context.WithValue(ctx, route.Key{}, inf), inf.Informer()
}
