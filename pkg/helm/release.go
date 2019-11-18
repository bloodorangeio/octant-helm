/*
Copyright 2019 Blood Orange

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

package helm // import "github.com/bloodorangeio/octant-helm/pkg/helm"

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	rspb "helm.sh/helm/v3/pkg/release"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type (
	tempReleaseSummary struct {
		Object         *unstructured.Unstructured
		RevisionNumber int
	}

	tempReleaseSummaryMap map[string]tempReleaseSummary
)

func UnstructuredListToHelmReleaseList(ul *unstructured.UnstructuredList) []*rspb.Release {
	var releases []*rspb.Release

	latestHelmReleases := getLatestHelmReleases(ul)

	for _, summary := range latestHelmReleases {
		obj := summary.Object
		secret, err := convertUnstructuredObjectToSecret(obj)
		if err != nil {
			log.Println(err)
			continue
		}

		release, err := convertSecretToHelmRelease(secret)
		if err != nil {
			log.Println(err)
			continue
		}

		releases = append(releases, release)
	}

	return releases
}

func UnstructuredListToHelmReleaseByName(ul *unstructured.UnstructuredList, releaseName string) *rspb.Release {
	var release *rspb.Release
	latestHelmReleases := getLatestHelmReleases(ul)
	for name, summary := range latestHelmReleases {
		if name == releaseName {
			obj := summary.Object
			secret, err := convertUnstructuredObjectToSecret(obj)
			if err != nil {
				log.Println(err)
				break
			}
			release, err = convertSecretToHelmRelease(secret)
			if err != nil {
				log.Println(err)
			}
			break
		}
	}
	return release
}

// List of secrets contains all revisions,
// so just get the latest revision for each release
func getLatestHelmReleases(ul *unstructured.UnstructuredList) tempReleaseSummaryMap {
	latestHelmReleases := tempReleaseSummaryMap{}

	for _, obj := range ul.Items {
		name := obj.GetName()
		tmp := strings.Split(name, ".")
		if len(tmp) != 6 {
			// Not a Helm-related secret
			// Note: valid Helm release secrets look like: sh.helm.release.v1.myrelease.v8
			continue
		}

		releaseName := tmp[4]
		revisionNumber, err := strconv.Atoi(strings.TrimPrefix(tmp[5], "v"))
		if err != nil {
			log.Println(err)
			continue
		}

		existingEntry, ok := latestHelmReleases[releaseName]
		if !ok || existingEntry.RevisionNumber < revisionNumber {
			latestHelmReleases[releaseName] = tempReleaseSummary{
				Object:         obj.DeepCopy(),
				RevisionNumber: revisionNumber,
			}
		}
	}

	return latestHelmReleases
}

func convertUnstructuredObjectToSecret(obj *unstructured.Unstructured) (*v1.Secret, error) {
	rawJSONBytes, err := obj.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var s v1.Secret
	err = json.Unmarshal(rawJSONBytes, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func convertSecretToHelmRelease(s *v1.Secret) (*rspb.Release, error) {
	val, ok := s.Data["release"]
	if !ok {
		return nil, errors.New("secret does not contain a \"release\" field")
	}

	return decodeRelease(string(val))
}

// Rest of the code found below copied from:
// https://github.com/helm/helm/blob/9b42702a4bced339ff424a78ad68dd6be6e1a80a/pkg/storage/driver/util.go

var b64 = base64.StdEncoding
var magicGzip = []byte{0x1f, 0x8b, 0x08}

func decodeRelease(data string) (*rspb.Release, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls rspb.Release
	// unmarshal protobuf bytes
	if err := json.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}
