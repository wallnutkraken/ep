package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wallnutkraken/ep/poddata"

	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dbPath string
	// data is defined here in cmd so that any command will have it readily initialized for
	// their own tasks. It's specifically loaded on application start instead of at a later point
	// to allow for the "dbpath" argument to be used.
	data poddata.Data
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ep",
	Short: "CLI-based podcast streamer.",
	Long:  `Streams podcasts from RSS feeds`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&dbPath, "dbpath", "", "Path to the ep database file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if dbPath == "" {
		// No path provided, use the default path
		dbPath = getDefaultDBPath()
	}
	// Start the database
	db, err := poddata.New(dbPath)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	data = db
}

func getDefaultDBPath() string {
	cfg := configdir.New("", "")
	configFolders := cfg.QueryFolders(configdir.Global)
	if len(configFolders) == 0 {
		logrus.Fatal("Could not find any configuration folders... what are you running this on?")
	}
	configDir := configFolders[0].Path
	// Add the filename onto the directory
	return filepath.Join(configDir, "epdata.sqlite")
}
