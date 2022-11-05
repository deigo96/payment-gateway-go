package helper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GoogleCloudStorage struct {
	GOOGLE_CLOUD_PROJECT_ID           string
	GOOGLE_APPLICATION_CREDENTIALS    string
	GOOGLE_CLOUD_STORAGE_BUCKET       string
	GOOGLE_CLOUD_STORAGE_PATH_PREFIX  string
	GOOGLE_CLOUD_STORAGE_API_URI      string
	GOOGLE_CLOUD_STORAGE_PATH_PROJECT string
}

type CredentialConfig struct {
	Type                        string
	Project_id                  string
	Private_key_id              string
	Private_key                 string
	Client_email                string
	Client_id                   string
	Auth_uri                    string
	Token_uri                   string
	Auth_provider_x509_cert_url string
	Client_x509_cert_url        string
}

type Bucket struct {
	BucketName string
}

type ProjectId struct {
	ProjectId string
}

func GetCredential() (s CredentialConfig) {
	jsonFile, err := ioutil.ReadFile("app.config.json")
	if err != nil {
		log.Fatal("Error load Credential File: ", err)
	}

	err = json.Unmarshal(jsonFile, &s)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return s
}

func GetGoogleCloudStorage() (g GoogleCloudStorage) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading dotenv file")
	}

	g.GOOGLE_CLOUD_PROJECT_ID = os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	g.GOOGLE_APPLICATION_CREDENTIALS = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	g.GOOGLE_CLOUD_STORAGE_BUCKET = os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET")
	g.GOOGLE_CLOUD_STORAGE_PATH_PREFIX = os.Getenv("GOOGLE_CLOUD_STORAGE_PATH_PREFIX")
	g.GOOGLE_CLOUD_STORAGE_API_URI = os.Getenv("GOOGLE_CLOUD_STORAGE_API_URI")
	g.GOOGLE_CLOUD_STORAGE_PATH_PROJECT = os.Getenv("GOOGLE_CLOUD_STORAGE_PATH_PROJECT")

	return g
}

func GetBucket() (b Bucket) {
	// credential, _ := GetCredential()
	google := GetGoogleCloudStorage()

	b.BucketName = google.GOOGLE_CLOUD_STORAGE_BUCKET
	// b.ProjectID = credential.Project_id

	return b
}
