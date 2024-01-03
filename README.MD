# Postgress Data Import Service

This services got triggert by RDS when a new order is created and exports that as JSON and puts it on S3 in tyhe /process 'folder'. 

## Installation

Service is deployed with Terraform

```bash
terraform apply
```
Uninstall

```bash
terraform destroy
```

## Usage

Compile the Go code and zip it so Terraform can deploy the AWS Lambda

```bash
./compile.sh