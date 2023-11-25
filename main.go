package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

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
	showHelp bool
	version  string = "1.0.0"
)

func escape(input string) string {
	return strconv.QuoteToASCII(input)
}

var rootCmd = &cobra.Command{
	Use:   "copier",
	Short: fmt.Sprintf("Upload a file to an SFTP server. Version %s", version),
	Run: func(cmd *cobra.Command, args []string) {

		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(escape(password)),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
		if err != nil {
			log.Fatal("Failed to dial:", err)
		}
		defer client.Close()

		sftpClient, err := sftp.NewClient(client)
		if err != nil {
			log.Fatal("Failed to create SFTP client:", err)
		}
		defer sftpClient.Close()

		localFile, err := os.Open(srcFile)
		if err != nil {
			log.Fatal("Failed to open local file:", err)
		}
		defer localFile.Close()

		remoteFile, err := sftpClient.Create(destFile)
		if err != nil {
			log.Fatal("Failed to create remote file:", err)
		}
		defer remoteFile.Close()

		_, err = io.Copy(remoteFile, localFile)
		if err != nil {
			log.Fatal("Failed to copy file contents:", err)
		}

		fmt.Printf("File '%s' sent to '%s'", srcFile, destFile)
		os.Exit(0)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "T", "", "SSH server host")
	rootCmd.Flags().IntVarP(&port, "port", "P", 22, "SSH server port")
	rootCmd.Flags().StringVarP(&username, "username", "U", "", "SSH username")
	rootCmd.Flags().StringVarP(&password, "password", "W", "", "SSH password")
	rootCmd.Flags().StringVarP(&srcFile, "src", "S", "", "Source file path")
	rootCmd.Flags().StringVarP(&destFile, "dest", "D", "", "Destination file path on the server")
	rootCmd.Flags().BoolVarP(&showHelp, "help", "H", false, "Show help message")

	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("src")
	rootCmd.MarkFlagRequired("dest")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
