# addon-installer
Installs addons from kustomize packages listed in a
configuration file of type `addons.config.x-k8s.io/AddonInstallerConfiguration`.

### usage
```shell
bin/installer --config demo/dupes.yaml
bin/installer --config demo/v1alpha1.yaml
```

### development
```shell
# fetch deps + regenerate all API's
make

# build just the binary from existing files
make only-build
```

## development on Mac OS
```shell
# fetch deps + regenerate all API's
make binary

```
