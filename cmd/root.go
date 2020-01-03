/*
Copyright © 2020 Bacworks sarl <contact@eyxance.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var readme string = `
=======================================================
Find below OpenVPN remote configuration
    
Remote site domain name= %s
Tunnel over port (local/remote)= %s
Protocol= %s
Cipher= %s
Bind to local address or port= %s
Interval for data channel renegociation (in seconds)= %s
Key Direction= %s

========================================================`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ovpn-splitter",
	Short: "Split client.ovpn to separate files and readme",
	Long: `OPVN Splitter is a CLI for OpenVPN applications.
This application create readme.txt, ca.crt, client.crt, client.key,
tls-auth files from a client.ovpn file.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Default opvn path
		path := "client.ovpn"

		// If user give a file path as arg, use this path instead
		if len(os.Args) == 2 {
			path = os.Args[1]
		}

		if _, err := os.Stat(path); err != nil {
			log.Fatalf("File '%s' does not exists.", path)
		}

		// Read file content
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Reading file: %s", err)
		}

		// Extract dir and filename from path
		var dir string = filepath.Dir(path)
		var fileName string
		_, fileName = filepath.Split(path)
		log.Printf("Reading '%s': OK", fileName)
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

		// Regexp Definition
		reCa := regexp.MustCompile(`<ca>\n([\(\)\#\n\-\+\/\= a-zA-Z0-9]+)</ca>`)
		reCert := regexp.MustCompile(`<cert>\n([\(\)\#\n\-\+\/\= a-zA-Z0-9]+)</cert>`)
		reKey := regexp.MustCompile(`<key>\n([\(\)\#\n\-\+\/\= a-zA-Z0-9]+)</key>`)
		reTlsAuth := regexp.MustCompile(`<tls-auth>\n([\(\)\#\n\-\+\/\= a-zA-Z0-9]+)</tls-auth>`)
		reRemote := regexp.MustCompile(`OVPN_ACCESS_SERVER_WSHOST=([\. a-zA-Z0-9]+):([0-9]+)`)
		reCipher := regexp.MustCompile(`cipher ([\-a-zA-Z0-9]+)`)
		reNegTime := regexp.MustCompile(`reneg-sec ([0-9]+)`)
		reKeyDir := regexp.MustCompile(`key-direction ([0-1])`)

		// Create crt and key files
		err = ioutil.WriteFile(dir+"/"+fileName+".ca.crt", reCa.Find(content), 0777)
		if err != nil {
			log.Printf("Writing '%s.ca.crt': %s", fileName, err)
		} else {
			log.Printf("Writing '%s.ca.crt': OK", fileName)
		}
		err = ioutil.WriteFile(dir+"/"+fileName+".crt", reCert.Find(content), 0777)
		if err != nil {
			log.Printf("Writing '%s.crt': %s", fileName, err)
		} else {
			log.Printf("Writing '%s.crt': OK", fileName)
		}
		err = ioutil.WriteFile(dir+"/"+fileName+".key", reKey.Find(content), 0777)
		if err != nil {
			log.Printf("Writing '%s.key': %s", fileName, err)
		} else {
			log.Printf("Writing '%s.key': OK", fileName)
		}
		err = ioutil.WriteFile(dir+"/"+fileName+".tls-auth", reTlsAuth.Find(content), 0777)
		if err != nil {
			log.Printf("Writing '%s.tls-auth': %s", fileName, err)
		} else {
			log.Printf("Writing '%s.tls-auth': OK", fileName)
		}

		// Create a Readme.txt
		domainTmp := reRemote.FindAllStringSubmatch(string(content), -1)
		var domainName string = domainTmp[0][1]
		cipherTmp := reCipher.FindStringSubmatch(string(content))
		var cipher string = string(cipherTmp[1])
		var bind string = "Yes"
		nobind, _ := regexp.Match(`nobind`, content)
		if nobind {
			bind = "No"
		}
		negTimeTmp := reNegTime.FindStringSubmatch(string(content))
		var negTime string = string(negTimeTmp[1])
		keyDirTmp := reKeyDir.FindStringSubmatch(string(content))
		var keyDir string = string(keyDirTmp[1])

		// Print to screen readme.txt content
		log.Printf(readme, domainName, "1194/1194", "UDP", cipher, bind, negTime, keyDir)

		// Save readme.txt to file
		err = ioutil.WriteFile(dir+"/"+fileName+".readme.txt", []byte(readme), 0777)
		if err != nil {
			log.Fatalf("Writing '%s.readme.txt': %s", fileName, err)
		}
		log.Printf("OpenVPN configuration is available in '%s.readme.txt'.", fileName)
	},
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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
