package main

import (
	"github.com/gophergala/ImgurGo/imagestore"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
)

type Factory struct {
	config *Configuration
}

func (this *Factory) NewImageStores() []imagestore.ImageStore {
	stores := []imagestore.ImageStore{}

	for _, configWrapper := range this.config.Stores {
		switch configWrapper.Type {
		case "S3StoreConfig":
			config, _ := configWrapper.Config.(S3StoreConfig)
			store := this.NewS3ImageStore(config)
			stores = append(stores, store)
		case "LocalStoreConfig":
			config, _ := configWrapper.Config.(LocalStoreConfig)
			store := this.NewLocalImageStore(config)
			stores = append(stores, store)
		default:
			log.Fatal("Unsupported store %s", configWrapper.Type)
		}
	}

	return stores
}

func (this *Factory) NewS3ImageStore(config S3StoreConfig) imagestore.ImageStore {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.Regions[config.Region])
	mapper := imagestore.NewNamePathMapper(config.NamePathRegex, config.NamePathMap)

	return imagestore.NewS3ImageStore(
		config.BucketName,
		config.StoreRoot,
		client,
		mapper,
	)
}

func (this *Factory) NewLocalImageStore(config LocalStoreConfig) imagestore.ImageStore {
	mapper := imagestore.NewNamePathMapper(config.NamePathRegex, config.NamePathMap)
	return imagestore.NewLocalImageStore(config.StoreRoot, mapper)
}

func (this *Factory) NewStoreObject(name string, mime string, imgType string) *imagestore.StoreObject {
	return &imagestore.StoreObject{
		Name:     name,
		MimeType: mime,
		Type:     imgType,
	}
}

func (this *Factory) NewHashGenerator(store imagestore.ImageStore) *HashGenerator {
	hashGen := &HashGenerator{
		make(chan string),
		this.config.HashLength,
		store,
	}

	hashGen.init()
	return hashGen
}
