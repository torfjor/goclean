package main

import (
	"context"
	"flag"
	"fmt"
	pb "goclean/gen"
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
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if flags.NArg() < 2 {
		return fmt.Errorf("usage: %s address name\n", flags.Name())
	}

	conn, err := grpc.Dial(flags.Arg(0), grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewGreetingServiceClient(conn)
	res, err := client.Greet(ctx, &pb.GreetRequest{Name: flags.Arg(1)})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, res.GetGreeting())
	return nil
}
