package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/lib/pq"
)

// OrderDetail struct representeert de structuur van je queryresultaten.
type OrderDetail struct {
	InstanceID                          string    `json:"instance_id"`
	AppID                               string    `json:"app_id"`
	AppdID                              string    `json:"appd_id"`
	Hyperscaler                         string    `json:"hyperscaler"`
	Environment                         string    `json:"environment"`
	InstanceClass                       string    `json:"instance_class"`
	SizeStorage                         string    `json:"size_storage"`
	DbName                              string    `json:"db_name"`
	PgMajorVersion                      string    `json:"pg_major_version"`
	Collation                           string    `json:"collation"`
	Encoding                            string    `json:"encoding"`
	MaintainanceWindowStartUtcDay       string    `json:"maintainance_window_start_utc_day"`
	MaintainanceWindowStartUtcStarttime string    `json:"maintainance_window_start_utc_starttime"`
	MaintainanceWindowStartUtcDuration  time.Time `json:"maintainance_window_start_utc_duration"`
	BackupRetentionDays                 int       `json:"backup_retention_days"`
}

func HandleRequest(ctx context.Context) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("AWS_REGION")

	// Verbindingsstring
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Maak verbinding met de database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Controleer de databaseverbinding
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Voer een query uit
	sqlStatement := `SELECT instance_id, app_id, appd_id, hyperscaler, environment, instance_class, 
                     size_storage, db_name, pg_major_version, "collation", encoding, 
                     maintainance_window_start_utc_day, maintainance_window_start_utc_starttime, 
                     maintainance_window_start_utc_duration, backup_retention_days
                     FROM s_cloudinv.v_order_details 
                     WHERE instance_id = 'cpds000004';`

	var detail OrderDetail
	row := db.QueryRow(sqlStatement)
	err = row.Scan(&detail.InstanceID, &detail.AppID, &detail.AppdID, &detail.Hyperscaler, &detail.Environment, &detail.InstanceClass,
		&detail.SizeStorage, &detail.DbName, &detail.PgMajorVersion, &detail.Collation, &detail.Encoding,
		&detail.MaintainanceWindowStartUtcDay, &detail.MaintainanceWindowStartUtcStarttime,
		&detail.MaintainanceWindowStartUtcDuration, &detail.BackupRetentionDays)
	if err != nil {
		log.Fatal(err)
	}

	// Converteer de struct naar JSON
	jsonData, err := json.Marshal(detail)
	if err != nil {
		log.Fatal(err)
	}

	// Upload naar S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Fatal(err)
	}

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("order-detail.json"),
		Body:   bytes.NewReader(jsonData),
	})
	if err != nil {
		log.Fatal("Failed to upload:", err)
	}

	fmt.Println("Successfully uploaded to S3")
}

func main() {
	lambda.Start(HandleRequest)
}
