package main

import (
	"goadmin/config"
	"goadmin/routes"
	"os"

	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "start",
	Short: "Starts the EAP server",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		port, _ := cmd.Flags().GetInt("port")
		startAdmin(dbUrl, port)
	},
}

func startAdmin(dbUrl string, port int) {
	// setup postgres db connection
	config.ConnectDb(dbUrl)

	// initialize http router
	router := httprouter.New()

	// get the current working directory
	workingDir, _ := os.Getwd()

	// setup the media directory for file uploads
	config.DefineMediaPath(workingDir + "/media")

	// set the server port for viewing media content
	config.SetPort(port)

	// serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// serve media files
	router.ServeFiles("/media/*filepath", http.Dir("media"))

	// configure route definitions
	router.GET("/tables/:tableName/new-object", routes.AddObject)
	router.POST("/tables/:tableName/new-object", routes.AddObject)
	router.GET("/tables/:tableName/old-object/:id", routes.EditObject)
	router.PATCH("/tables/:tableName/old-object/:id", routes.EditObject)
	router.DELETE("/tables/:tableName/old-object", routes.DeleteObject)
	router.GET("/tables/:tableName", routes.ListTableObjects)
	router.GET("/tables", routes.ListTables)

	// display the access URL on the terminal
	fmt.Printf("Server is running on: http://localhost:%d/tables\n", port)

	// run the http server
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	command.Flags().StringP("db", "d", "postgres://postgres:pgadmin@localhost:5432/goadmin", "PostgreSQL database URL")
	command.Flags().IntP("port", "p", 7000, "Port number to run the admin panel server")
}

func main() {
	var rootCmd = &cobra.Command{Use: "EasyAdminPanel"}
	rootCmd.AddCommand(command)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
