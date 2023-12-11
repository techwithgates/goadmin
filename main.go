package main

import (
	"embed"
	"fmt"

	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"

	"github.com/techwithgates/goadmin/config"
	"github.com/techwithgates/goadmin/routes"
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

//go:embed template/*
//go:embed static/*
//go:embed media/*
var embedder embed.FS

func startAdmin(dbUrl string, port int) {
	// setup postgres db connection
	config.ConnectDb(dbUrl)

	// initialize http router
	router := httprouter.New()

	workingDir, _ := os.Getwd()

	// setup the media directory for file uploads
	config.DefineMediaPath(workingDir + "/media")

	// set the server port for viewing media content
	config.SetPort(port)

	// set the embedder to apply in routes
	routes.SetEmbedder(&embedder)

	// set the static root directory
	staticRoot, err := fs.Sub(embedder, "static")

	// set the media root directory
	mediaRoot, err := fs.Sub(embedder, "media")

	// serve embedded static files
	router.ServeFiles("/static/*filepath", http.FS(staticRoot))

	// serve embedded media files
	router.ServeFiles("/media/*filepath", http.FS(mediaRoot))

	// configure route definitions
	router.GET("/tables/:tableName/new-object", routes.AddObject)
	router.POST("/tables/:tableName/new-object", routes.AddObject)
	router.GET("/tables/:tableName/old-object/:id", routes.EditObject)
	router.PATCH("/tables/:tableName/old-object/:id", routes.EditObject)
	router.DELETE("/tables/:tableName/old-object", routes.DeleteObject)
	router.GET("/tables/:tableName", routes.ListTableObjects)
	router.GET("/tables", routes.ListTables)

	fmt.Printf("EAP server is running on: http://localhost:%d/tables\n", port)

	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatal(err)
	}
}

func init() {
	command.Flags().StringP("db", "d", "postgres://postgres:pgadmin@localhost:5432/goadmin", "PostgreSQL database URL")
	command.Flags().IntP("port", "p", 7000, "Port number to run EAP server")
}

func main() {
	var rootCmd = &cobra.Command{Use: "EasyAdminPanel"}
	rootCmd.AddCommand(command)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
