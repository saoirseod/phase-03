#Script to generate certificate authority and server/client certificates 
#cA are needed to generate a cert, it's also used to verify the cert owner
#deletes all previous files
rm *.pem
touch client-ext.conf
touch server-ext.conf

echo "subjectAltName=IP:0.0.0.0" > client-ext.conf
echo "subjectAltName=IP:0.0.0.0" > server-ext.conf

#SERVER SIDE
#Creates certificates authority's private key and self-signed certificate
#req creates and processes certificate requests in PKCS#10 format and can create self-signed certs for root cert authority
#because it is a self signed certificate, the Issuer and Subject will be the same
#-subj to automatically input identity info
#-nodes makes sure that the private key will not be encrypted 

#---CA private key and cert---
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=IR/ST=Clare/L=Shannon/O=Intel/OU=RandD/CN=root/emailAddress=s.od@intel.com"
#x509 means that the cert is self signed 

echo "The certificate authority's self-signed certificate is:"
#-noout to make sure that the original encoded value is not outputted
openssl x509 -in ca-cert.pem -noout -text

#Similar to first openssl command but we dont self-sign, hence, x509 and day's are removed, since no certificate is generated
#We only want a CSR
#Creates the web server's private key and certificate signing request

#---Server private key and CSR---
openssl req -newkey rsa:4096 -nodes -keyout  server-key.pem -out server-req.pem -subj "/C=IR/ST=Clare/L=Shannon/O=Intel/OU=WebServer/CN=root/emailAddress=w.s@intel.com"

#---Signing the Server Certificate Request---

openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf

#Using CA's private key to sign web servers CSR and retrieve signed certificate
#-req option to show that we will pass in a certificate request, -in to pass in the request file that follows
#The -CAcreateserial parameter is used so that the Certificate authority can make sure that each cert that is signed has a unique serial number, and a file will be generated if it does
#-extfile tells openssl that extra options were added

echo "The server's signed certificate is:"

openssl x509 -in server-cert.pem -noout -text

echo "Verifying the certificate.."

openssl verify -CAfile ca-cert.pem server-cert.pem

echo "Verification finshed."

#CLIENT SIDE

#---command to generate the client's private key and CSR---

openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=IR/ST=Clare/L=Shannon/O=Intel/OU=RandD/CN=root/emailAddress=example.client@intel.com"

#---signing the client's CSR---
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.conf

echo "The client's signed certificate:"
openssl x509 -in client-cert.pem -noout -text
