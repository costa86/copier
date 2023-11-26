package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	host     string
	port     int
	username string
	password string
	srcFile  string
	destFile string
	version  string = "1.0.0"
)

func handleFailure(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "copier",
	Short: fmt.Sprintf("Upload a file via SFTP to a remote server\nDeveloped by PCA team\nVersion %s", version),
	Run: func(cmd *cobra.Command, args []string) {

		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)

		handleFailure(err)

		defer client.Close()

		sftpClient, err := sftp.NewClient(client)

		handleFailure(err)

		defer sftpClient.Close()

		localFile, err := os.Open(srcFile)

		handleFailure(err)

		defer localFile.Close()

		remoteFile, err := sftpClient.Create(destFile)

		handleFailure(err)

		defer remoteFile.Close()

		_, err = io.Copy(remoteFile, localFile)

		handleFailure(err)

		fmt.Printf("File '%s' sent to '%s'", srcFile, destFile)
		os.Exit(0)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "t", "", "SSH server host")
	rootCmd.Flags().IntVarP(&port, "port", "p", 22, "SSH server port")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "SSH username")
	rootCmd.Flags().StringVarP(&password, "password", "w", "", "SSH password")
	rootCmd.Flags().StringVarP(&srcFile, "src", "s", "", "Source file path")
	rootCmd.Flags().StringVarP(&destFile, "dest", "d", "", "Destination file path on the server")

	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("src")
	rootCmd.MarkFlagRequired("dest")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		handleFailure(err)
	}
}
