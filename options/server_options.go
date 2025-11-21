package options

import (
	"strconv"

	"github.com/spf13/pflag"
)

type ServerOptions struct {
	Host           string   `json:"host"  mapstructure:"host"  description:"Server host address"`
	Port           int      `json:"port"`
	Mode           string   `json:"mode"`
	Middlewares    []string `json:"middlewares" mapstructure:"middlewares"`
	Node           int      `json:"node" mapstructure:"node" description:"Node ID for the server, used for clustering and routing purposes"`
	TrustedProxies []string `json:"trusted-proxies" mapstructure:"trusted-proxies" description:"List of trusted proxy IP addresses or CIDR blocks, comma separated. If this list is empty default trusted proxies will be used."`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Host:           "0.0.0.0",
		Port:           9000,
		Mode:           "release",
		Node:           1,
		TrustedProxies: []string{},
	}
}

func (o *ServerOptions) Validate() []error {
	var errs []error

	return errs
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "server.host", o.Host, "Server host address")
	fs.IntVar(&o.Port, "server.port", o.Port, "Server port")
	fs.StringVar(&o.Mode, "server.mode", o.Mode, "Server mode")
	fs.StringSliceVar(&o.Middlewares, "server.middlewares", o.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.")
	fs.IntVar(&o.Node, "server.node", o.Node, "Node ID for the server, used for clustering and routing purposes")
	fs.StringSliceVar(&o.TrustedProxies, "server.trusted-proxies", o.TrustedProxies, "List of trusted proxy IP addresses or CIDR blocks, comma separated. If this list is empty default trusted proxies will be used.")
}

func (o *ServerOptions) Address() string {
	return o.Host + ":" + strconv.Itoa(o.Port)
}
