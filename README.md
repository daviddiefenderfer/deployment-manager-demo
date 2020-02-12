### Summary
A basic deployment script written in Go that will
deploy a given config file. This guide outlines deploying
a docker image with the sample config.yaml provided.

The purpose of this is to provide a easily reproducible
program in hopes to aid in a feature request to get a 
deployment `Outputs` property back from a successful deploy
through google-api-go-client. In the current state deploying 
the config file with the `gcloud` program will successfully
return the `Outputs` property. Deploying from the google-api-go-client, will not.

### Setup
You will need to have gcloud installed on your system and a project with billing attached.
Once the project has successfully been created
replace the {{projectName}} values in config.yaml file with this
same project name. From there you can enable the needed apis.

`gcloud services enable deploymentmanager.googleapis.com --project {{projectName}}`

`gcloud services enable compute.googleapis.com --project {{projectName}}`

Now that there is a project and a valid config we're
ready to begin with deploying our app via `gcloud`.


###Deploying with `gcloud`cli:

For this part we see the expected `Outputs` behavior per the [documentation](https://cloud.google.com/deployment-manager/docs/configuration/expose-information-outputs) by deploying via `gcloud`.
 
`gcloud deployment-manager deployments create hello-world --config config.yaml --project sample-dm-demo`

If there were no issues with your setup you should successfully of deployed this config file and 
received a response with the outputs defined in the config.yaml file.

```bash
NAME  TYPE                 STATE      ERRORS  INTENT
vm    compute.v1.instance  COMPLETED  []
OUTPUTS    VALUE
IPAddress  35.192.141.3

```

We can delete this now to re-deploy the same config from the api-go-client

`gcloud deployment-manager deployments delete hello-world`

###Deploying with google-api-go-client
The go program written here is very basic and prints out the json response from a successful
deployment. As you will see, this deployment will not return any "FinalValues" or "Outputs" properties


For this one you will need credentials. You can use `gcloud auth application-default login` if you feel comfortable
installing the credentials in your local user home. Else you will need to download your service-account json key
and set it with the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

`go run main.go --config config.yaml --project {{projectName}}`

You will see here the Resource response aligns well with the expected response per the [API](https://cloud.google.com/deployment-manager/docs/reference/latest/resources) docs. 
We don't see anything regarding the Outputs value is any of the API models. 

###Conclusion
The response from deploying via CLI provides a desirable "Outputs" property that is either unavailable
in the api-go-client library or is not in an intuitive location to access.