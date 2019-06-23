package main

import (
	"fmt"
	"github.com/n-inja/gomniauth-traq"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	// NOTE: Don't change this, the auth settings on the providers
	// are coded to this path for this example.
	Address string = ":8000"
)

func write(ctx context.Context, output string) {
	ctx.HttpResponseWriter().Write([]byte(output))
}

func writeHeader(ctx context.Context) {
	write(ctx, "Gomniauth - Example web app")
}

func respondWithError(ctx context.Context, errorMessage string) error {
	writeHeader(ctx)
	write(ctx, fmt.Sprintf("Error: %s", errorMessage))
	return nil
}

func main() {

	// setup the providers
	gomniauth.SetSecurityKey("zN06SFqw6t6HTKBKqa71FT9mOBYRr2D2UuFNcT579bEp8mgl1iF41Wm5xGD1ioCG")

	clientID := os.Getenv("TRAQ_CLIENT_ID")
	clientSecret := os.Getenv("TRAQ_CLIENT_SECRET")
	if clientID == "" {
		log.Fatal("you should set TRAQ_CLIENT_ID")
	}
	if clientSecret == "" {
		log.Fatal("you should set TRAQ_CLIENT_SECRET")
	}

	gomniauth.WithProviders(
		gomniauth_traq.New(clientID, clientSecret, "http://localhost:8000/auth/traq/callback"))

	goweb.Map("/", func(ctx context.Context) error {

		return goweb.Respond.With(ctx, http.StatusOK, []byte(`
      <html>
        <body>
          <h2>Log in with...</h2>
          <ul>
            <li>
              <a href="auth/traq/login">traQ</a>
            </li>
          </ul>
        </body>
      </html>
    `))

	})

	/*
	   GET /auth/{provider}/login
	   Redirects them to the fmtin page for the specified provider.
	*/
	goweb.Map("auth/{provider}/login", func(ctx context.Context) error {

		provider, err := gomniauth.Provider(ctx.PathValue("provider"))

		if err != nil {
			return err
		}

		state := gomniauth.NewState("after", "success")

		// if you want to request additional scopes from the provider,
		// pass them as login?scope=scope1,scope2
		//options := objx.MSI("scope", ctx.QueryValue("scope"))

		authUrl, err := provider.GetBeginAuthURL(state, nil)

		if err != nil {
			return err
		}

		// redirect
		return goweb.Respond.WithRedirect(ctx, authUrl)

	})

	goweb.Map("auth/{provider}/callback", func(ctx context.Context) error {

		provider, err := gomniauth.Provider(ctx.PathValue("provider"))

		if err != nil {
			return err
		}

		creds, err := provider.CompleteAuth(ctx.QueryParams())

		if err != nil {
			return err
		}

		/*
			// get the state
			state, stateErr := gomniauth.StateFromParam(ctx.QueryValue("state"))
			if stateErr != nil {
				return stateErr
			}
			// redirect to the 'after' URL
			afterUrl := state.GetStringOrDefault("after", "error?e=No after parameter was set in the state")
		*/

		// load the user
		user, userErr := provider.GetUser(creds)

		fmt.Println(user.Name())
		fmt.Println(user.Nickname())
		fmt.Println(user.AuthCode())

		if userErr != nil {
			return userErr
		}

		return goweb.API.RespondWithData(ctx, user)

		// redirect
		//return goweb.Respond.WithRedirect(ctx, afterUrl)

	})

	/*
	   ----------------------------------------------------------------
	   START OF WEB SERVER CODE
	   ----------------------------------------------------------------
	*/

	log.Println("Starting...")
	fmt.Print("Gomniauth - Example web app\n")
	fmt.Print("by Mat Ryer and Tyler Bunnell\n")
	fmt.Print(" \n")
	fmt.Print("Starting Goweb powered server...\n")

	// make a http server using the goweb.DefaultHttpHandler()
	s := &http.Server{
		Addr:           Address,
		Handler:        goweb.DefaultHttpHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	listener, listenErr := net.Listen("tcp", Address)

	fmt.Printf("  visit: %s\n", Address)

	if listenErr != nil {
		log.Fatalf("Could not listen: %s", listenErr)
	}

	fmt.Println("\n")
	fmt.Println("Try some of these routes:\n")
	fmt.Printf("%s", goweb.DefaultHttpHandler())
	fmt.Println("\n\n")

	go func() {
		for _ = range c {

			// sig is a ^C, handle it

			// stop the HTTP server
			fmt.Print("Stopping the server...\n")
			listener.Close()

			/*
			   Tidy up and tear down
			*/
			fmt.Print("Tearing down...\n")

			// TODO: tidy code up here

			log.Fatal("Finished - bye bye.  ;-)\n")

		}
	}()

	// begin the server
	log.Fatalf("Error in Serve: %s\n", s.Serve(listener))

	/*
	   ----------------------------------------------------------------
	   END OF WEB SERVER CODE
	   ----------------------------------------------------------------
	*/

}
