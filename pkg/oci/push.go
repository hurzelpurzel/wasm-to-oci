package oci

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

// Push pushes a WASM module to an OCI registry
func Push(ref, module string, opts content.RegistryOptions) error {
	ctx, registry, store := newORASContext(opts)

	contents, err := ioutil.ReadFile(module)
	if err != nil {
		return err
	}
	desc, err := store.Add(module, ContentLayerMediaType, contents)
	if err != nil {
		return err
	}
	manifest, manifestDesc, config, configDesc, err := content.GenerateManifestAndConfig(nil, nil, desc)
	if err != nil {
		return err
	}
	store.Set(configDesc, config)
	err = store.StoreManifest(ref, manifestDesc, manifest)
	if err != nil {
		return err
	}

	log.Infof("Pushing %s to %s...\n", module, ref)

	desc, err = oras.Copy(ctx, store, ref, registry, "")
	if err != nil {
		return err
	}
	log.Infof("Size: %v", desc.Size)
	log.Infof("Pushed to %s with digest %s\n", ref, desc.Digest)

	return nil
}
