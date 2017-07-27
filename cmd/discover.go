package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/djboris9/go-upnp/discovery"
	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover UPnP devices on your network",
	Long:  `Discover is used to find UPnP devices on your network`,

	Run: func(cmd *cobra.Command, args []string) {
		t, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			log.Fatal(err)
		}

		devs, err := discovery.Discover(time.Duration(t) * time.Second)
		if err != nil {
			log.Fatal(err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		for _, d := range devs {
			fmt.Fprintf(w, "FriendlyName\t%v\n", d.FriendlyName)
			fmt.Fprintf(w, "Manufacturer\t%v\n", d.Manufacturer)
			fmt.Fprintf(w, "ModelName\t%v\n", d.ModelName)
			fmt.Fprintf(w, "ModelDescription\t%v\n", d.ModelDescription)
			fmt.Fprintf(w, "SerialNumber\t%v\n", d.SerialNumber)
			fmt.Fprintf(w, "UDN\t%v\n", d.UDN)
			fmt.Fprintf(w, "Location\t%v\n", d.Location.String())
			fmt.Fprintln(w, "")
		}
		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(discoverCmd)

	discoverCmd.PersistentFlags().IntP("timeout", "t", 3, "Discovery timeout in seconds")
}
