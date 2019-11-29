# grpc-service


Simple go [gRPC](https://grpc.io/) service for Cloud Run.

> Note, to keep this readme short, I will be asking you to execute scripts rather than listing here complete commands. You should really review each one of these scripts for content, and, to understand the individual commands so you can use them in the future.

## Pre-requirements

### Setup

If you don't have one already, start by creating new project and configuring your [Google Cloud SDK](https://cloud.google.com/sdk/docs/). In addition, you will also need to install the `alpha` component.

```shell
gcloud components install alpha
```

Also, if you have not done so already, you will have [set up Cloud Run](https://cloud.google.com/run/docs/setup).


### Config

All the variables used in this sample are defined in the [bin/config](bin/config) file. You can edit these to your preferred values.

* `SERVICE` - (default "grpc-service") - this is the name of the service. It's used to build image, create service account, and in Cloud Run deployment. If you do need to change that, edit it before any other step.
* `SERVICE_REGION` - (default "us-central1") - this is the GCP region to which you want to deploy this service. For complete list of regions where Cloud Run service is currently offered see https://cloud.google.com/about/locations/
* `SERVICE_VERSION` - (default "v1") - this is the version of the service which is used to tag the container image as well as the Cloud Run service

## Deployment

```shell
bin/setup
```

```shell
bin/db
```

Returns

```shell
Creating instance...done.
Creating database...done.
Schema updating...done.
```

Next, build the server container image which will be used to deploy Cloud Run service using the [bin/image](bin/image) script:

```shell
bin/image
```

Returns

```shell
gcr.io/PROJECT/distributed-echo:0.5.6  SUCCESS
```

Once the container image and service account are ready, you can now deploy the new service using [bin/deploy](bin/deploy) script:

```shell
bin/deploy
```






