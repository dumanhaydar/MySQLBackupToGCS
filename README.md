# MySQL Backup to GCS

Build

    go build -o bin/mbg main.go


Execute

    ./bin/mbg -database *dbName* -bucket *gcsBucketName* -config *pathToMySQLconfig* -keep *keepSize*
