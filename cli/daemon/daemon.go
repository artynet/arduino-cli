//
// This file is part of arduino-cli.
//
// Copyright 2018 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to modify or
// otherwise use the software for commercial activities involving the Arduino
// software without disclosing the source code of your own applications. To purchase
// a commercial license, send an email to license@arduino.cc.
//

package daemon

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/arduino/arduino-cli/cli/globals"
	"github.com/arduino/arduino-cli/commands/daemon"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// NewCommand created a new `daemon` command
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "daemon",
		Short:   "Run as a daemon",
		Long:    "Running as a daemon the initialization of cores and libraries is done only once.",
		Example: "  " + os.Args[0] + " daemon",
		Args:    cobra.NoArgs,
		Run:     runDaemonCommand,
	}
}

func runDaemonCommand(cmd *cobra.Command, args []string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	userAgentValue := fmt.Sprintf("%s/%s daemon (%s; %s; %s) Commit:%s/Build:%s", globals.VersionInfo.Application,
		globals.VersionInfo.VersionString, runtime.GOARCH, runtime.GOOS, runtime.Version(), globals.VersionInfo.Commit, globals.VersionInfo.BuildDate)
	headers := http.Header{"User-Agent": []string{userAgentValue}}

	coreServer := daemon.ArduinoCoreServerImpl{
		DownloaderHeaders: headers,
		VersionString:     globals.VersionInfo.VersionString,
		Config:            globals.Config,
	}
	rpc.RegisterArduinoCoreServer(s, &coreServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("Done serving")
}
