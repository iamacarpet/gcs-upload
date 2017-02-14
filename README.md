# gcs-upload
Tool to take a directory, transparently add it to a tar.gz archive (using parallel gzip) and upload to GCS all in one shot.

```
Usage:
    For this example, we assume your GCS bucket is "gs://backup-storage",
    Our JSON service account credentials file is located at "./creds.json",
    And the directory we want to compress and upload is "LOCAL_DIRECTORY", which will be compressed to "LOCAL_DIRECTORY.tar.gz".

    $ GCS_BUCKET=backup-storage GOOGLE_APPLICATION_CREDENTIALS="./creds.json" ./gcs-upload "LOCAL_DIRECTORY"

    Proxy can be configured using the standard "https_proxy" ENV variable (thanks Go!).
```

The process is logged to syslog on the local machine with an app name of "gcs-upload".
