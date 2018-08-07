package main

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/images/oci"
	"github.com/containerd/containerd/namespaces"
)

func main() {
	if err := run(namespaces.WithNamespace(context.Background(), "default"), os.Args[1]); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, path string) error {
	if len(os.Args) != 3 {
		return fmt.Errorf("expect 2 arguments <tar> <name>")
	}
	client, err := containerd.New(defaults.DefaultAddress)
	if err != nil {
		return err
	}
	defer client.Close()
	image, err := importImage(ctx, client, path)
	if err != nil {
		return err
	}
	return client.Install(ctx, image)
}

func importImage(ctx context.Context, client *containerd.Client, path string) (containerd.Image, error) {
	name := os.Args[2]
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	images, err := client.Import(
		ctx,
		&oci.V1Importer{
			ImageName: name,
		},
		f,
	)
	if err != nil {
		return nil, err
	}
	return images[0], nil
}
