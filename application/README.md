Altering phase 02 code so that the client and server is a single application, being built using docker file. Aim is to have the application deployed so that there are two pods, one for the client and one for the server, as well as a single service that will be used for communication. Once the application is deployed within k8's, the secrets can be mounted to a pod which will replace the generated certificates.

Just using a random rpi image here.

So far, to run dockerfile use: 

docker build --build-arg http_proxy=http://proxy.ch.intel.com:911/ -t saoirseod/rpi-hostname:v1 -f docker/Dockerfile .

or

include '--build-arg https_proxy=http://proxy.ch.intel.com:911/' also

Then to retag the image:

docker tag saoirseod/rpi-hostname:latest saoirseod/rpi-hostname:v1

followed by:

docker image save -o /tmp/image.tar saoirseod/rpi-hostname:v1

k3s ctr images import /tmp/image.tar


golang-app .yaml files still to be edited, error with pods running
