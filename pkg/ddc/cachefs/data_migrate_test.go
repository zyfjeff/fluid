/*
  Copyright 2023 The Fluid Authors.

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

package cachefs

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brahma-adshonor/gohook"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	"github.com/fluid-cloudnative/fluid/api/v1alpha1"
	cdatamigrate "github.com/fluid-cloudnative/fluid/pkg/datamigrate"
	cruntime "github.com/fluid-cloudnative/fluid/pkg/runtime"
	"github.com/fluid-cloudnative/fluid/pkg/utils/fake"
	"github.com/fluid-cloudnative/fluid/pkg/utils/helm"
)

func TestCacheFSEngine_CreateDataMigrateJob(t *testing.T) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-dataset-cachefs-values",
			Namespace: "fluid",
		},
		Data: map[string]string{
			"data": valuesConfigMapData,
		},
	}

	mockExecCheckReleaseCommon := func(name string, namespace string) (exist bool, err error) {
		return false, nil
	}
	mockExecCheckReleaseErr := func(name string, namespace string) (exist bool, err error) {
		return false, errors.New("fail to check release")
	}
	mockExecInstallReleaseCommon := func(name string, namespace string, valueFile string, chartName string) error {
		return nil
	}
	mockExecInstallReleaseErr := func(name string, namespace string, valueFile string, chartName string) error {
		return errors.New("fail to install datamigrate chart")
	}

	wrappedUnhookCheckRelease := func() {
		err := gohook.UnHook(helm.CheckRelease)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	wrappedUnhookInstallRelease := func() {
		err := gohook.UnHook(helm.InstallRelease)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	targetDataMigrate := v1alpha1.DataMigrate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "fluid",
		},
		Spec: v1alpha1.DataMigrateSpec{
			From: v1alpha1.DataToMigrate{
				DataSet: &v1alpha1.DatasetToMigrate{
					Name:      "test-dataset",
					Namespace: "fluid",
				},
			},
			To: v1alpha1.DataToMigrate{
				ExternalStorage: &v1alpha1.ExternalStorage{
					URI: "minio://test/test",
				},
			},
		},
	}
	datasetInputs := []v1alpha1.Dataset{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dataset",
				Namespace: "fluid",
			},
		},
	}
	podListInputs := []corev1.PodList{{
		Items: []corev1.Pod{{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"a": "b"},
			},
		}},
	}}
	testObjs := []runtime.Object{}
	testObjs = append(testObjs, configMap)
	for _, datasetInput := range datasetInputs {
		testObjs = append(testObjs, datasetInput.DeepCopy())
	}
	for _, podInput := range podListInputs {
		testObjs = append(testObjs, podInput.DeepCopy())
	}
	testScheme.AddKnownTypes(corev1.SchemeGroupVersion, configMap)
	client := fake.NewFakeClientWithScheme(testScheme, testObjs...)
	engine := &CacheFSEngine{
		name:      "cachefs",
		namespace: "fluid",
		Client:    client,
		Log:       fake.NullLogger(),
	}
	ctx := cruntime.ReconcileRequestContext{
		Log:      fake.NullLogger(),
		Client:   client,
		Recorder: record.NewFakeRecorder(1),
	}

	err := gohook.Hook(helm.CheckRelease, mockExecCheckReleaseErr, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = engine.CreateDataMigrateJob(ctx, targetDataMigrate)
	if err == nil {
		t.Errorf("fail to catch the error: %v", err)
	}
	wrappedUnhookCheckRelease()

	err = gohook.Hook(helm.CheckRelease, mockExecCheckReleaseCommon, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = gohook.Hook(helm.InstallRelease, mockExecInstallReleaseErr, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = engine.CreateDataMigrateJob(ctx, targetDataMigrate)
	if err == nil {
		t.Errorf("fail to catch the error: %v", err)
	}
	wrappedUnhookInstallRelease()

	err = gohook.Hook(helm.InstallRelease, mockExecInstallReleaseCommon, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = engine.CreateDataMigrateJob(ctx, targetDataMigrate)
	if err != nil {
		t.Errorf("fail to exec the function: %v", err)
	}
	wrappedUnhookCheckRelease()
}

func TestCacheFSEngine_generateDataMigrateValueFile(t *testing.T) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-dataset-cachefs-values",
			Namespace: "fluid",
		},
		Data: map[string]string{
			"data": ``,
		},
	}

	datasetInputs := []v1alpha1.Dataset{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dataset",
				Namespace: "fluid",
			},
			Spec: v1alpha1.DatasetSpec{},
		},
	}

	testObjs := []runtime.Object{}
	testObjs = append(testObjs, configMap)
	for _, datasetInput := range datasetInputs {
		testObjs = append(testObjs, datasetInput.DeepCopy())
	}
	client := fake.NewFakeClientWithScheme(testScheme, testObjs...)

	context := cruntime.ReconcileRequestContext{
		Client: client,
	}

	dataMigrateNoTarget := v1alpha1.DataMigrate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-datamigrate",
			Namespace: "fluid",
		},
		Spec: v1alpha1.DataMigrateSpec{
			From: v1alpha1.DataToMigrate{
				DataSet: &v1alpha1.DatasetToMigrate{
					Name:      "test-dataset",
					Namespace: "fluid",
				},
			},
			To: v1alpha1.DataToMigrate{
				ExternalStorage: &v1alpha1.ExternalStorage{
					URI: "minio://test/test",
				},
			},
		},
	}
	dataMigrateWithTarget := v1alpha1.DataMigrate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-datamigrate",
			Namespace: "fluid",
		},
		Spec: v1alpha1.DataMigrateSpec{
			From: v1alpha1.DataToMigrate{
				DataSet: &v1alpha1.DatasetToMigrate{
					Name:      "test-dataset",
					Namespace: "fluid",
					Path:      "/test/",
				},
			},
			To: v1alpha1.DataToMigrate{
				ExternalStorage: &v1alpha1.ExternalStorage{
					URI: "minio://test/test",
				},
			},
			Options: map[string]string{
				"exclude": "4.png",
			},
		},
	}

	var testCases = []struct {
		dataMigrate    v1alpha1.DataMigrate
		expectFileName string
	}{
		{
			dataMigrate:    dataMigrateNoTarget,
			expectFileName: filepath.Join(os.TempDir(), "fluid-test-datamigrate-migrate-values.yaml"),
		},
		{
			dataMigrate:    dataMigrateWithTarget,
			expectFileName: filepath.Join(os.TempDir(), "fluid-test-datamigrate-migrate-values.yaml"),
		},
	}

	for _, test := range testCases {
		engine := CacheFSEngine{
			name:      "cachefs",
			namespace: "fluid",
			Client:    client,
			Log:       fake.NullLogger(),
		}
		fileName, err := engine.generateDataMigrateValueFile(context, test.dataMigrate)
		if err != nil {
			t.Errorf("fail to generate the datamigrate value file: %v", err)
		}
		if !strings.Contains(fileName, test.expectFileName) {
			t.Errorf("got value: %v, want value: %v", fileName, test.expectFileName)
		}
	}
}

func TestCacheFSEngine_genDataUrl(t *testing.T) {
	cacheFSSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cachefs-secret",
		},
		Data: map[string][]byte{
			"access-key": []byte(base64.StdEncoding.EncodeToString([]byte("test"))),
			"secret-key": []byte(base64.StdEncoding.EncodeToString([]byte("test"))),
			"metaurl":    []byte(base64.StdEncoding.EncodeToString([]byte("test"))),
		},
	}
	cacheFSConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cachefs-values",
			Namespace: "default",
		},
		Data: map[string]string{
			"data": valuesConfigMapData,
		},
	}
	testObjs := []runtime.Object{}
	testObjs = append(testObjs, (*cacheFSSecret).DeepCopy(), cacheFSConfigMap)
	client := fake.NewFakeClientWithScheme(testScheme, testObjs...)

	type args struct {
		data v1alpha1.DataToMigrate
		info *cdatamigrate.DataMigrateInfo
	}
	tests := []struct {
		name        string
		args        args
		wantDataUrl string
		wantErr     bool
	}{
		{
			name: "test-external",
			args: args{
				data: v1alpha1.DataToMigrate{ExternalStorage: &v1alpha1.ExternalStorage{
					URI: "http://minio/",
					EncryptOptions: []v1alpha1.EncryptOption{
						{
							Name: "access-key",
							ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
								Name: "cachefs-secret",
								Key:  "access-key",
							}},
						},
						{
							Name: "secret-key",
							ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
								Name: "cachefs-secret",
								Key:  "secret-key",
							}},
						},
						{
							Name: "token",
							ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
								Name: "cachefs-secret",
								Key:  "token",
							}},
						},
					},
				}},
				info: &cdatamigrate.DataMigrateInfo{
					EncryptOptions: []v1alpha1.EncryptOption{
						{
							Name: "access-key",
							ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
								Name: "cachefs-secret",
								Key:  "access-key",
							}},
						},
						{
							Name: "secret-key",
							ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
								Name: "cachefs-secret",
								Key:  "secret-key",
							}},
						},
						{
							Name: "token",
							ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
								Name: "cachefs-secret",
								Key:  "token",
							}},
						},
					},
					Options: map[string]string{},
				},
			},
			wantDataUrl: "http://${ACCESS_KEY}:${SECRET_KEY}:${TOKEN}@minio/",
			wantErr:     false,
		},
		{
			name: "test-external-subpath",
			args: args{
				data: v1alpha1.DataToMigrate{ExternalStorage: &v1alpha1.ExternalStorage{
					URI: "http://minio/test/",
					EncryptOptions: []v1alpha1.EncryptOption{{
						Name: "access-key",
						ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
							Name: "cachefs-secret",
							Key:  "access-key",
						}},
					}},
				}},
				info: &cdatamigrate.DataMigrateInfo{
					EncryptOptions: []v1alpha1.EncryptOption{{
						Name: "access-key",
						ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
							Name: "cachefs-secret",
							Key:  "access-key",
						}},
					}},
					Options: map[string]string{},
				},
			},
			wantDataUrl: "http://${ACCESS_KEY}:@minio/test/",
			wantErr:     false,
		},
		{
			name: "test-external-subpath-file",
			args: args{
				data: v1alpha1.DataToMigrate{ExternalStorage: &v1alpha1.ExternalStorage{
					URI: "http://minio/test",
					EncryptOptions: []v1alpha1.EncryptOption{{
						Name: "access-key",
						ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
							Name: "cachefs-secret",
							Key:  "access-key",
						}},
					}},
				}},
				info: &cdatamigrate.DataMigrateInfo{
					EncryptOptions: []v1alpha1.EncryptOption{{
						Name: "access-key",
						ValueFrom: v1alpha1.EncryptOptionSource{SecretKeyRef: v1alpha1.SecretKeySelector{
							Name: "cachefs-secret",
							Key:  "access-key",
						}},
					}},
					Options: map[string]string{},
				},
			},
			wantDataUrl: "http://${ACCESS_KEY}:@minio/test",
			wantErr:     false,
		},
		{
			name: "test-dataset",
			args: args{
				data: v1alpha1.DataToMigrate{
					DataSet: &v1alpha1.DatasetToMigrate{
						Name:      "test",
						Namespace: "default",
						Path:      "/subpath/",
					},
				},
				info: &cdatamigrate.DataMigrateInfo{
					EncryptOptions: []v1alpha1.EncryptOption{},
					Options:        map[string]string{},
				},
			},
			wantDataUrl: "jfs://FLUID_METAURL/subpath/",
			wantErr:     false,
		},
		{
			name: "test-dataset-no-path",
			args: args{
				data: v1alpha1.DataToMigrate{
					DataSet: &v1alpha1.DatasetToMigrate{
						Name:      "test",
						Namespace: "default",
					},
				},
				info: &cdatamigrate.DataMigrateInfo{
					EncryptOptions: []v1alpha1.EncryptOption{},
					Options:        map[string]string{},
				},
			},
			wantDataUrl: "jfs://FLUID_METAURL/",
			wantErr:     false,
		},
		{
			name: "test-dataset-subpath-file",
			args: args{
				data: v1alpha1.DataToMigrate{
					DataSet: &v1alpha1.DatasetToMigrate{
						Name:      "test",
						Namespace: "default",
						Path:      "/subpath",
					},
				},
				info: &cdatamigrate.DataMigrateInfo{
					EncryptOptions: []v1alpha1.EncryptOption{},
					Options:        map[string]string{},
				},
			},
			wantDataUrl: "jfs://FLUID_METAURL/subpath",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &CacheFSEngine{
				Client: client,
				Log:    fake.NullLogger(),
			}
			gotDataUrl, err := j.genDataUrl(tt.args.data, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("genDataUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDataUrl != tt.wantDataUrl {
				t.Errorf("genDataUrl() gotDataUrl = %v, want %v", gotDataUrl, tt.wantDataUrl)
			}
		})
	}
}
