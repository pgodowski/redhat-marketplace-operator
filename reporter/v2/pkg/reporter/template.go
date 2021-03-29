// Copyright 2021 IBM Corp.
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

package reporter

import (
	"bytes"
	"reflect"

	"text/template"

	"emperror.dev/errors"
	sprig "github.com/Masterminds/sprig/v3"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
)

type ReportTemplater struct {
	templFieldMap map[string]*template.Template
}

type ReportLabels struct {
	Label map[string]interface{}
}

func NewTemplate(promLabels *common.MeterDefPrometheusLabels) (*ReportTemplater, error) {
	if promLabels == nil {
		return nil, errors.New("metric is nil")
	}

	templater := &ReportTemplater{
		templFieldMap: make(map[string]*template.Template),
	}
	t := reflect.ValueOf(*promLabels)

	for i := 0; i < t.NumField(); i++ {
		_, ok := t.Type().Field(i).Tag.Lookup("template")

		if !ok {
			continue
		}

		v := t.Field(i).Interface()

		fieldName := t.Type().Field(i).Name
		str, ok := v.(string)

		if !ok {
			return nil, errors.NewWithDetails("template fields must be strings", "fieldName", fieldName)
		}

		templ := template.New(fieldName).Funcs(sprig.GenericFuncMap())
		templ, err := templ.Parse(str)

		if err != nil {
			return nil, errors.Wrap(err, "failed to parse template")
		}

		templater.templFieldMap[fieldName] = templ
	}

	return templater, nil
}

func (r *ReportTemplater) Execute(
	promLabels *common.MeterDefPrometheusLabels,
	values *ReportLabels) error {
	if promLabels == nil {
		return errors.New("metric is nil")
	}

	t := reflect.ValueOf(promLabels).Elem()

	for fieldName, tpl := range r.templFieldMap {
		var buff bytes.Buffer

		if !t.FieldByName(fieldName).CanSet() {
			continue
		}

		err := tpl.Execute(&buff, values)

		if err != nil {
			return errors.WrapIfWithDetails(err, "failed to execute buffer", "fieldName", fieldName)
		}

		str := buff.String()
		t.FieldByName(fieldName).SetString(str)
	}

	return nil
}
