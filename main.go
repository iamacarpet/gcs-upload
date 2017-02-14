package main

import (
    "os"
    "io"
    "mime"
    "path"
    "archive/tar"
    "path/filepath"
    "golang.org/x/net/context"
    "github.com/klauspost/pgzip"
    "cloud.google.com/go/storage"
    "github.com/iamacarpet/gcs-upload/logger"
)

// Requires several environment variables:
//  GCS_BUCKET for the name of the GCS bucket to upload to, omitting the gs://
//  GOOGLE_APPLICATION_CREDENTIALS as the path to the credentials JSON file.
//  https_proxy for the proxy address, if required.
// And the first argument is the path to the directory to archive and upload.

func main(){
    logger.Init("gcs-upload")

    if len(os.Args) != 2 {
        logger.Fatalf("Invalid Number of Arguments - Unable to determine directory to upload.")
    }
    root_directory := os.Args[1]
    if ok := check_directory(root_directory); ! ok {
        logger.Fatalf("Invalid directory argument (%s), that isn't a directory.", root_directory)
    }
    fileName := path.Base(root_directory) + ".tar.gz"
    gsBucketPath := os.Getenv("GCS_BUCKET")

    logger.Infof("Starting upload of directory \"%s\" to GCS Bucket \"%s\" as compressed archive \"%s\".", root_directory, gsBucketPath, fileName)

    ctx := context.Background()

    client, err := storage.NewClient(ctx)
    if err != nil {
        logger.Fatalf("Failed to connect to GCS: %s", err)
    }

    bkt := client.Bucket(gsBucketPath)

    obj := bkt.Object(fileName)
    w := obj.NewWriter(ctx)
    if err := create_and_upload(root_directory, w); err != nil {
        logger.Fatalf("Error while archiving and uploading directory: %s", err)
    }
    if err := w.Close(); err != nil {
        logger.Fatalf("Failed to close file write handle for GCS Object (%s): %s", fileName, err)
    }

    attrs := storage.ObjectAttrsToUpdate{
        ContentType:    mime.TypeByExtension("gz"),
    }
    if _, err := obj.Update(ctx, attrs); err != nil {
		logger.Fatalf("Failed to update object attributes for file (%s): %s", fileName, err)
	}

    logger.Infof("Completed upload of directory \"%s\" to GCS Bucket \"%s\" as compressed archive \"%s\".", root_directory, gsBucketPath, fileName)
}

func check_directory(root_directory string) bool {
    fi, err := os.Stat(root_directory)
    if err != nil {
        return false
    }
    mode := fi.Mode()
    if mode.IsDir() {
        return true
    } else {
        return false
    }
}

func create_and_upload(directory string, obj io.Writer) error {
    gw := pgzip.NewWriter(obj)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

    directory, _ = filepath.Abs(directory)
    root_directory, _ := filepath.Abs(path.Dir(directory))

    walkFn := func(path string, info os.FileInfo, err error) error {
        if info.Mode().IsDir() {
            return nil
        }
        // Because of scoping we can reference the external root_directory variable

        new_path := path[len(root_directory)+1:]
        if len(new_path) == 0 {
            return nil
        }
        fr, err := os.Open(path)
        if err != nil {
            return err
        }
        defer fr.Close()

        if h, err := tar.FileInfoHeader(info, new_path); err != nil {
            logger.Fatalf("Failed to generate tar header for file \"%s\" to archive: %s", new_path, err)
        } else {
            h.Name = new_path
            if err = tw.WriteHeader(h); err != nil {
                logger.Fatalf("Failed to write tar header for file \"%s\" to archive: %s", new_path, err)
            }
        }
        if length, err := io.Copy( tw, fr ); err != nil {
            logger.Fatalf("Failed to copy file \"%s\" to archive: %s", new_path, err)
        } else {
            logger.Infof("Successfully added file \"%s\" to archive, length: %d", new_path, length)
        }
        return nil
    }

    if err := filepath.Walk(directory, walkFn); err != nil {
        return err
    }

    return nil
}
