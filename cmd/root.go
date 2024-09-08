package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func forwardRequestToOrigin(r *http.Request, origin *url.URL) (*http.Response, error) {
	req := r.Clone(r.Context())
	req.URL.Scheme = origin.Scheme
	req.URL.Host = origin.Host
	req.Host = origin.Host
	req.RequestURI = ""
	return http.DefaultClient.Do(r)
}

func writeResponse(w http.ResponseWriter, res *http.Response) {
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	if res.Body != nil {
		defer res.Body.Close()
		io.Copy(w, res.Body)
	}
}

var rootCmd = &cobra.Command{
	Use:   "caching-proxy",
	Short: "Caching Proxy Server",
	Long:  `Caching Proxy Server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetUint16("port")
		origin, _ := cmd.Flags().GetString("origin")

		originUrl, err := url.Parse(origin)
		if err != nil {
			return fmt.Errorf("error parsing origin: %s", err)
		}

		fmt.Printf("Starting caching proxy server on port %d\n", port)
		fmt.Printf("Caching proxy server will proxy requests to %v\n", originUrl)

		err = http.ListenAndServe(
			fmt.Sprintf("localhost:%d", port),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				res, err := forwardRequestToOrigin(r, originUrl)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				writeResponse(w, res)
			}))

		if err != nil {
			return fmt.Errorf("error starting caching proxy server: %s", err)
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Uint16P("port", "p", 3000, "Port to listen on")
	rootCmd.Flags().StringP("origin", "o", "", "Origin to proxy requests to")
	rootCmd.MarkFlagRequired("origin")
}
