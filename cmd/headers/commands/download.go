package commands

import (
	"github.com/EVRICE/tgeth_alpha/cmd/headers/download"
	"github.com/spf13/cobra"
)

var (
	bufferSizeStr string // Size of buffer
	combined      bool   // Whether downloader also includes sentry
	timeout       int    // Timeout for delivery requests
	window        int    // Size of sliding window for downloading block bodies
)

func init() {
	downloadCmd.Flags().StringVar(&filesDir, "filesdir", "", "path to directory where files will be stored")
	downloadCmd.Flags().StringVar(&bufferSizeStr, "bufferSize", "512M", "size o the buffer")
	downloadCmd.Flags().StringVar(&sentryAddr, "sentryAddr", "localhost:9091", "sentry address <host>:<port>")
	downloadCmd.Flags().BoolVar(&combined, "combined", false, "run downloader and sentry in the same process")
	downloadCmd.Flags().IntVar(&timeout, "timeout", 30, "timeout for devp2p delivery requests, in seconds")
	downloadCmd.Flags().IntVar(&window, "window", 65536, "size of sliding window for downloading block bodies, block")

	// Options below are only used in the combined mode
	downloadCmd.Flags().StringVar(&natSetting, "nat", "any", "NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")
	downloadCmd.Flags().IntVar(&port, "port", 30303, "p2p port number")
	downloadCmd.Flags().StringArrayVar(&staticPeers, "staticpeers", []string{}, "static peer list [enode]")
	downloadCmd.Flags().BoolVar(&discovery, "discovery", true, "discovery mode")
	downloadCmd.Flags().StringVar(&netRestrict, "netrestrict", "", "CIDR range to accept peers from <CIDR>")

	withChaindata(downloadCmd)
	withLmdbFlags(downloadCmd)
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download headers backwards",
	RunE: func(cmd *cobra.Command, args []string) error {
		db := openDatabase(chaindata)
		defer db.Close()
		if combined {
			return download.Combined(natSetting, port, staticPeers, discovery, netRestrict, filesDir, bufferSizeStr, db, timeout, window)
		}
		return download.Download(filesDir, bufferSizeStr, sentryAddr, coreAddr, db, timeout, window)
	},
}
