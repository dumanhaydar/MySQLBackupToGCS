
# MySQL Backup to Google Cloud Storage (GCS)

## Description

This project provides an automated solution for backing up MySQL databases to Google Cloud Storage. It is designed to be simple, efficient, and easy to configure.

## Features

-   Complete backup of MySQL databases.
-   Secure storage in Google Cloud Storage.
-   Customizable configuration for different environments.
-   Backup size management with the automatic deletion of old backups.

## Prerequisites

-   GoLang installed on the machine.
-   Access to a MySQL instance.
-   A configured account and bucket on Google Cloud Storage.
-   `bzip2` and `mysqldump` installed for database compression and dumping.

## Installation and Configuration

Clone the repository and build the executable:

    git clone https://github.com/dumanhaydar/MySQLBackupToGCS.git
    cd MySQLBackupToGCS
    go build -o bin/mbg main.go

Ensure `bzip2` and `mysqldump` are installed:

bashCopy code

    sudo apt-get install bzip2 
    sudo apt-get install mysql-client

## Usage

To execute the backup script, use the following command:

    ./bin/mbg -database *dbName* -bucket *gcsBucketName* -config *pathToMySQLconfig* -keep *keepSize*

Where:

-   `*dbName*` is the name of the MySQL database.
-   `*gcsBucketName*` is the name of the Google Cloud Storage bucket.
-   `*pathToMySQLconfig*` is the path to the MySQL configuration file _(exemple .my.cnf)_
-   `*keepSize*` is the number of backups to keep.

## Security

Ensure credentials and configurations are managed securely.

## Contribution

Contributions are welcome. Please submit your pull requests on GitHub.
