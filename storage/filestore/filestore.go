package filestore

import (
	"app/config"
	"app/pkg/logs"
	appStorage "app/storage"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"net/url"
	"os"
)

var (
	ErrCouldNotOpenFile  = fmt.Errorf("could_not_open_file")
	ErrCouldNotCopyFile  = fmt.Errorf("could_not_copy_file")
	ErrCouldNotCloseFile = fmt.Errorf("could_not_close_file")
)

type filestoreKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type Filestore struct {
	bucket *storage.BucketHandle
	log    logs.LoggerInterface
}

const (
	FolderCategory = appStorage.Folder("category")
	FolderProduct  = appStorage.Folder("product")
)

func NewFilestore(cfg *config.FileStorageConfig, log logs.LoggerInterface) appStorage.FileStorageInterface {
	ctx := context.Background()
	fileKey := filestoreKey{
		Type:                    os.Getenv("FIREBASE_TYPE"),
		ProjectID:               os.Getenv("FIREBASE_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		PrivateKey:              "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDVH/SeJAifMuOb\n7V4dtL0L+YxK3AJjndTQ0eLEWdyFfnVRGtcAawNwJfOXVMTSjJwXNXgMLP9/W3+C\nvvsqxXXTalVgo6kWn1ICz9ZfmfF+o/QPHiQeoujqkslkbd2LYKEg2WReSKq5naet\n7nLuCp7bD94xkatvKns9H09vbjCDNpXZkemdqoYllxnGloQntv9lAfDvPYCeSAyc\nylRVeCNEF+2xV5BENJz9tnshzfNHD0bbjAutZ2STfjH1++DXHtjcZmUCU8hGkYrY\nU1EGfPVc/7QOxT/ZQchFDaWyhNvdBoodHPQeA2jYd0R244w9DK75VRqLc58YsCdR\nQua7Hx/pAgMBAAECggEAHQprLdymyMae6YGJMw63OeP/Crj69wcBLOkvTLvFRvLJ\n2qDdI3icMdJsubYP0cmI1XLFkRHVSM3iIfc9yKtlGEdOK0pxwRw5iRo3/F6WNRWd\n19eEFXaeZU8AyUIH5Qo9j85/lNyzAsg+79OPmUHHjqhqtEzEqZ5bck3gvZmBW8eQ\nXcos99EU3EU5vPdnV0V/rs6+zh0Qeo3mh7DwLsp9W2TIm0Sxu8cOFCVwU+YJrk9z\nOgjydfcQl664rJY6ITlQkSZiQKWy4TTWvcroJt+KriBu7X29hVdAJkkvQOlEpNu/\n5BLTuU+wXz2IG78gWi5DpzM1D1/05/G5nifX4C94PwKBgQD6iqYpsNr2x69NEBY4\nESZGwUY39E9i1XJyx5QzjzbfSSdGs0T7xPB+Sgo6vLevTfN4g19XmxfRrJGsVbxy\nwIzN/nb29wSpVnDKiBVsR+Dk8Yg5cbi6n7X9c+698SflVQ84RMi1mRlrz9AJaiPc\nbXK5HebMa5A+NflL+W+bu9sCYwKBgQDZxJ8BgDDIj72I7ekA9PBa8oudx4UUCMxl\nfAjS2e4IJDdg1p5K9m5Jyr1Y6OMFbSJRL23h6Fi4wRzkUqxaah8mzwBRRDqCMyUT\nntKvjzovMFj7Cdlm+WCFFqNMH/fOOgTlWJGip0YVxnzyUtF3ieKPKfLc6DshGFdr\nECeNtrSAQwKBgEjOo+z3pRoT+2B0rVBLw4jKP8Kg77Tz/FdYoju9gZ+vnYdRL1nO\n6Gh60bAyCVsbVwaNftZxjqFy+b5QB/x88i4mpaGtNSCUqyBgHYGi/brqacDvyFQL\nd5KY7ycpfoOJjWu3qXAEdru632TtAFDdSXp8Mwbyty8s9i5a5VEnbUSrAoGAbquI\nC2E0aZjzP9V4pq3UQMQmxCaTsRzPk3u3mEB8wdJ1+lbX10zpu8K2+6pPRYCzAgNS\nmo5UGIC7yCVjxgdMkZJ9nM9J1MVdQF1kwSfO8BBoCBx3SefOb5STpKpSa5H8zvl1\n+e18prBa62O/ZDrE0vEEpdO3yRfvxU9OaqzBirUCgYAFK88d2UKLrDTQlU0Zka6P\nolewqVkT5kQxksggXRoRqfY3krb5qvN4h+jJSqf/BrUZ0pmGHXPoue8Xon/vYVUs\nVdILJa5YO6qCb2O6rokmK6l5gJ259l8n2M8a1oN7LeYtZR8RnGg02vfkFB9JQR+A\nE2hlnO4kWabRKKMAYEMl1A==\n-----END PRIVATE KEY-----\n",
		ClientEmail:             os.Getenv("FIREBASE_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FIREBASE_CLIENT_ID"),
		AuthURI:                 os.Getenv("FIREBASE_AUTH_URI"),
		TokenURI:                os.Getenv("FIREBASE_TOKEN_URI"),
		AuthProviderX509CertUrl: os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertUrl:       os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
		UniverseDomain:          os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
	}
	fileBytes, err := json.Marshal(&fileKey)
	if err != nil {
		log.Error("could not marshal key", logs.Error(err))
		panic(err)
	}

	opt := option.WithCredentialsJSON(fileBytes)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Error("could not create firebase app", logs.Error(err))
		panic(err)
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Error("could not get client for firebase storage", logs.Error(err))
		panic(err)
	}
	bucket, err := storageClient.Bucket(cfg.URL)
	if err != nil {
		log.Error("could not get the bucket", logs.Error(err))
		panic(err)
	}
	return &Filestore{
		bucket: bucket,
		log:    log,
	}
}

func (fs *Filestore) Create(model *multipart.FileHeader, imageFolder appStorage.Folder, id string) (string, error) {
	file, err := model.Open()
	if err != nil {
		fs.log.Error("could not open file", logs.Error(err))
		return "", ErrCouldNotOpenFile
	}
	defer file.Close()

	idPath := string(imageFolder) + "/" + id
	if os.Getenv("ENV") == config.LocalMode {
		idPath = "test/" + idPath
	}

	fmt.Println(idPath)
	wc := fs.bucket.Object(idPath).NewWriter(context.Background())

	wc.ObjectAttrs.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": idPath,
	}

	if _, err := io.Copy(wc, file); err != nil {
		fs.log.Error("could not copy image file", logs.Error(err))
		return "", ErrCouldNotCopyFile
	}
	if err := wc.Close(); err != nil {
		fs.log.Error("could not close writer", logs.Error(err))
		return "", ErrCouldNotCloseFile
	}

	return idPath, nil
}

func (fs *Filestore) GetURL(path string) string {
	name := url.PathEscape(path)
	link := fmt.Sprintf(
		`https://firebasestorage.googleapis.com/v0/b/ulab-market.appspot.com/o/%s?alt=media&token=%s`,
		name, name,
	)
	return link
}

func (fs *Filestore) DeleteFile(path string) error {
	ctx := context.Background()

	file := fs.bucket.Object(path)

	if err := file.Delete(ctx); err != nil {
		return err
	}
	return nil
}
