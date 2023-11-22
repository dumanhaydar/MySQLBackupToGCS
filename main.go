package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"time"
)

var mysqldumpPath string
var bzip2Path string

func main() {
	var err error
	dbName := flag.String("database", "", "Nom de la base de données")
	bucketName := flag.String("bucket", "", "Nom du bucket GCS")
	configPath := flag.String("config", "./.my.cnf", "Chemin vers le fichier de configuration MySQL")
	keepSize := flag.Int("keep", 5, "Nombre de sauvegardes à conserver")

	flag.Parse()

	mysqldumpPath, err = findExecutablePath("mysqldump")
	if err != nil {
		fmt.Println("mysqldump non trouvé:", err)
		return
	}

	bzip2Path, err = findExecutablePath("bzip2")
	if err != nil {
		fmt.Println("bzip2 non trouvé:", err)
		return
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFile := "backup-" + *dbName + "-" + timestamp + ".sql.bz2"

	if err := backupAndUploadToGCS(*dbName, *configPath, *bucketName, backupFile); err != nil {
		log.Fatal(err)
	}

	if err := cleanupOldBackups(*bucketName, *keepSize); err != nil {
		log.Fatal(err)
	}
}

func backupAndUploadToGCS(dbName, configPath, bucketName, gcsFileName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(gcsFileName)
	wc := object.NewWriter(ctx)
	defer wc.Close()

	mysqldumpCmd := exec.Command(
		"mysqldump",
		"--defaults-extra-file="+configPath,
		"--opt",
		"--complete-insert",
		"--single-transaction",
		"--max_allowed_packet=32M",
		"--max_allowed_packet=32M",
		"--databases", dbName,
	)
	bzipCmd := exec.Command("bzip2", "-9")

	bzipCmd.Stdin, _ = mysqldumpCmd.StdoutPipe()
	bzipCmd.Stdout = wc
	bzipCmd.Stderr = os.Stderr
	mysqldumpCmd.Stderr = os.Stderr

	if err := bzipCmd.Start(); err != nil {
		return err
	}

	if err := mysqldumpCmd.Run(); err != nil {
		return err
	}

	if err := bzipCmd.Wait(); err != nil {
		return err
	}

	return nil
}

func cleanupOldBackups(bucketName string, keep int) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	it := bucket.Objects(ctx, nil)

	var backups []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		backups = append(backups, attrs.Name)
	}

	sort.Strings(backups)

	for i := 0; i < len(backups)-keep; i++ {
		object := bucket.Object(backups[i])
		if err := object.Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}

func findExecutablePath(executableName string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", executableName)
	} else {
		cmd = exec.Command("which", executableName)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	return out.String(), err
}
