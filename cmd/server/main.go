package main

import (
	"context"
	"flag"
	"fmt"
	"goclean"
	"goclean/authorization"
	pb "goclean/gen"
	"goclean/inmem"
	"net"
	"net/http"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		httpFlag = flags.Bool("http", false, "run http server")
		grpcFlag = flags.Bool("grpc", true, "run grpc server")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	greeter := goclean.NewAuthorizedGreetingFunc(
		goclean.NewGreeterFunc(&inmem.Store{Templates: []string{"Hello, %s!", "Bonjour, %s!"}}),
		goclean.Greet,
	)

	done := make(chan struct{})
	if *grpcFlag {
		ln, err := net.Listen("tcp", "0.0.0.0:8080")
		if err != nil {
			return err
		}

		go func() {
			grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(authorization.GRPCAuthorizer()))
			pb.RegisterGreetingServiceServer(grpcSrv, struct {
				goclean.GreeterFunc
				pb.UnsafeGreetingServiceServer
			}{
				greeter,
				pb.UnimplementedGreetingServiceServer{},
			})

			go func() {
				grpcSrv.Serve(ln)
			}()

			<-done
			grpcSrv.GracefulStop()
		}()
	}

	if *httpFlag {
		go func() {
			srv := http.Server{Addr: "0.0.0.0:8081", Handler: authorization.HTTPAuthorizer(greeter)}
			go func() {
				srv.ListenAndServe()
			}()

			<-done
			srv.Shutdown(context.Background())
		}()
	}

	<-ctx.Done()
	close(done)

	return nil
}
