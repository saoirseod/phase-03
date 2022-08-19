altering phase 02 code so that the client and server is a single application, being built using docker file. Aim is to have the application deployed so that there are two pods, one for the client and one for the server, as well as a single service that will be used for communication. Once the application is deployed within k8's, the secrets can be mounted to a pod which will replace the generated certificates

so far, to run dockerfile use 
docker build --build-arg http_proxy=http://proxy.ch.intel.com:911/ -t saoirseod/rpi-hostname:latest -f docker/Dockerfile .
or
include --build-arg https_proxy=http://proxy.ch.intel.com:911/ also
