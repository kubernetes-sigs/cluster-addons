This cloudbuild job is invoked by prow after every commit to master,
and will push images to the staging image repository.

It can also be run against your own GCP project, which is helpful when
testing or adding to the cloudbuild job.  An example command is:

```
gcloud builds submit --config=dev/staging/cloudbuild.yaml . --substitutions=_DOCKER_IMAGE_PREFIX=$(gcloud config get-value project)/
```

