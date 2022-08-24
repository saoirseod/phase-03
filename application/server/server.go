//combined application with server and client, using flags and assigned port in struct to start application, file should be a single executable, not two seperate client and server files
//alterations to be made, client working?
//NOTE: pods not in running state until this compiles, reading in of old cert files to be removed once secrets are mounted!

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const defaultPort = 9000

var clientFlag = false

type server struct {
	pb.UnimplementedGreeterServer
	port int //server port
}

//%v is the value in a default format when printing structs
//SayHello method implements helloworlds' GreeterServer
//built off of golang's grpc hello world example
func (s *server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Name received: %v", request.GetName())
	return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {

	//Define the flag.
	//Bind the flag.
	//Parse the flag.
	args := server{}

	flag.IntVar(&args.port, "dport", defaultPort, "Specify application port")
	flag.BoolVar(&clientFlag, "boolFlag", false, "A boolean flag")

	flag.Parse()

	fmt.Println("If Boolean Flag is equal to false then client is in use, otherwise it is the server.")
	fmt.Println("Boolean Flag is ", clientFlag)

	if !clientFlag {
		serverFunc(args.port)

	} else {
		clientFunc(args.port)
	}

}

func serverFunc(port int) {
	//assigns listening port
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("Oops! Failed to listen: %v", err)
	}

	//here, loading in certificate authorities certificate. The process followed is that the client sends a cert to the server,
	//the server then uses the certificate authorities key to make sure that the cert sent is valid
	caPem, err := ioutil.ReadFile("../certs/ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	//cert pool understanding reference https://golang.hotexamples.com/examples/crypto.x509/-/NewCertPool/golang-newcertpool-function-examples.html
	//cert pool made that adds certificate authorities cert, x509 is the type of cert
	//x509 certs generally contain a public key
	certPool := x509.NewCertPool()
	//.AppendCertsFromPEM is a golang function, it parses PEM encoded certificates and appends them (adds) to the pool
	certPool.AppendCertsFromPEM(caPem)

	//here the .LoadX509KeyPair method is used to load the servers cert and key
	serverCert, err := tls.LoadX509KeyPair("../certs/server-cert.pem", "../certs/server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// configuration of the certificate what we want to
	configuration := &tls.Config{
		//pulls in fetched server cert from the method above
		Certificates: []tls.Certificate{serverCert},
		//.RequireAndVerifyClientCert indicates that a client's cert should be requested upon handshake
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
	}

	//creates and holds all of the tls credentials
	tls := credentials.NewTLS(configuration)

	//initiates the grpc server with the stated grpc tls credentials
	grpcServerSetup := grpc.NewServer(grpc.Creds(tls))

	//link back to code at the top, same as in hello-world go example
	//the service is now registered in the server
	pb.RegisterGreeterServer(grpcServerSetup, &server{})

	//Print out the port being listened at so that we know its working
	log.Printf("Currently listening at port %v", lis.Addr())
	if err := grpcServerSetup.Serve(lis); err != nil {
		log.Fatalf("Uh oh! There's an error with the grpc server: %v", err)
	}
}

func clientFunc(port int) {
	//here, loading in certificate authorities certificate. The process followed is that the client sends a cert to the server,
	//the server then uses the certificate authorities key to make sure that the cert sent is valid
	//this is used from the clients side this time to amke sure that the server is verified
	caCert, err := ioutil.ReadFile("../certs/ca-cert.pem")
	if err != nil {
		log.Fatal(caCert)
	}

	//cert pool understanding reference https://golang.hotexamples.com/examples/crypto.x509/-/NewCertPool/golang-newcertpool-function-examples.html
	//append is used to add
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	//reads in the client cert
	clientCert, err := tls.LoadX509KeyPair("../certs/client-cert.pem", "../certs/client-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	configuration := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	tls := credentials.NewTLS(configuration)

	//establishes the client connection, only works because of the server and client -ext.conf files that specify an IP
	connection, err := grpc.Dial(
		"0.0.0.0:9000",
		grpc.WithTransportCredentials(tls),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	client := pb.NewGreeterClient(connection)
	log.Printf("Connection created!")

	//This is where the server is communicated with and response is printed if all worked
	//using context here to get the time to show it below
	//for best practice with using the cancel function, it must get called at some point
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//called here
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Saoirse O'Donovan"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Communication secured using mTLS...")
	log.Printf("Greeting: %s", resp.GetMessage())

}
