/*

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

package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/addon-operators/installer/pkg/apis/config"
	"sigs.k8s.io/addon-operators/installer/pkg/apis/config/scheme"
	"sigs.k8s.io/addon-operators/installer/pkg/apis/config/v1alpha1"
)

func readConfig(f *flags) (*config.AddonInstallerConfiguration, error) {
	// The internal config object to be populated and written to STDOUT
	cfg := &config.AddonInstallerConfiguration{}
	// If the config flag was specified, try to deserialize the file
	var err error
	if f.configFileChanged {
		err = decodeFileInto(*f.configFile, cfg)
	}
	if err != nil {
		return cfg, err
	}
	// If the dry-run flag was specified, override the config
	if f.dryRunChanged {
		cfg.DryRun = *f.dryRun
	}
	return cfg, nil
}

// decodeFileInto reads a file and decodes the it into an internal type
func decodeFileInto(filePath string, obj runtime.Object) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Regardless of if the bytes are of any external version,
	// it will be read successfully and converted into the internal version
	return runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), content, obj)
}

func fprintConfig(w io.Writer, cfg *config.AddonInstallerConfiguration) error {
	cfgbytes, err := marshalYAML(cfg, v1alpha1.SchemeGroupVersion)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s", cfgbytes)
	return nil
}

// marshalYAML marshals any ComponentConfig object registered in the scheme for the specific version
func marshalYAML(obj runtime.Object, groupVersion schema.GroupVersion) ([]byte, error) {
	// yamlEncoder is a generic-purpose encoder to YAML for this scheme
	yamlEncoder := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	// versionSpecificEncoder writes out YAML bytes for exactly this v1alpha1 version
	versionSpecificEncoder := scheme.Codecs.EncoderForVersion(yamlEncoder, groupVersion)
	// Encode the object to YAML for the given version
	return runtime.Encode(versionSpecificEncoder, obj)
}
