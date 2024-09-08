package main

import (
	"context"
	"example/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api | ", log.LstdFlags) // Initializes a new Logger object that uses Standard Output, uses the name "product-api" and uses default standard flags.
	ph := handlers.NewProducts(l)
	gh := handlers.NewGoodbye(l)

	// Defines a router under the gorilla framework.
	sm := mux.NewRouter()
	// creates a subrouter of 'sm' that only monitors a GET method. This will allow us to later register handlers onto this subrouter.
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.Use(handlers.ContentTypeApplicationJsonMiddleware)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}", ph.UpdateProducts)

	// Use() means the router object will have to go through the MiddlewareValidateProduct middleware to  validate JSON and then call the next handler.
	putRouter.Use(ph.MiddlewareValidateProduct)
	putRouter.Use(handlers.ContentTypeApplicationJsonMiddleware)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)
	postRouter.Use(handlers.ContentTypeApplicationJsonMiddleware)

	sm.Handle("/goodbye", gh)
	// Why do reference a server rather than initialize it?
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,                // set the default handler.
		ErrorLog:     l,                 // set the default logg 	er.
		IdleTimeout:  120 * time.Second, // Max time given to connections using TCP Keep-Alive
		ReadTimeout:  1 * time.Second,   // Max time given to read requests from a client.
		WriteTimeout: 1 * time.Second,   // Max time given to write requests to a client.
	}

	// Executing an anonymous function concurrently in the background.
	go func() {
		err := s.ListenAndServe() // ListenAndServe takes the initialized HTTPServer object 's'.
		if err != nil {           // If no error is picked up.
			l.Fatal(err) // report a fatal error to the originally initializeed Logger object.
		}

	}()

	sigChannel := make(chan os.Signal, 1)      // Creates a channel for the server to reveive OS signal types.
	signal.Notify(sigChannel, os.Interrupt)    // Requesting to be notified of an OS interrupt signal has been sent out.
	signal.Notify(sigChannel, syscall.SIGTERM) // Requesting to be notified of an OS kill signal has been sent out.

	sig := <-sigChannel // send channel data to a variable 'sig'. The program will be blocked until either an Interrupt or a Kill signal is sent.
	l.Println("Received a terminate message, shutting down", sig)
	// context.Background() returns an empty context object. As such, this should only be used in main or top-level handler as you would want to derive child contexts instead of creating more empty context objs.
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Create a child context tc from context. tc receives a tc.Done signal after 30 seconds.
	defer cancel()                                                          // The context will be retained in memory indefinitely until program shuts down (next line), causing a memory leak.
	s.Shutdown(tc)                                                          // Shutdown HTTPserver object 's' once the child context 'tc' receives a Done signal.

}
