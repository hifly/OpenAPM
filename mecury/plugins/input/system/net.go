package system

import (
	"fmt"
	"net"

	"github.com/aiyun/openapm/mecury/agent"
)

type NetIOStats struct {
	ps PS

	skipChecks bool
	Interfaces []string
}

func (_ *NetIOStats) Description() string {
	return "Read metrics about network interface usage"
}

var netSampleConfig = `
  ## By default, telegraf gathers stats from any up interface (excluding loopback)
  ## Setting interfaces will tell it to gather these explicit interfaces,
  ## regardless of status.
  ##
  # interfaces = ["eth0"]
`

func (_ *NetIOStats) SampleConfig() string {
	return netSampleConfig
}

func (s *NetIOStats) Gather(acc agent.Accumulator) error {
	netio, err := s.ps.NetIO()
	if err != nil {
		return fmt.Errorf("error getting net io info: %s", err)
	}

	for _, io := range netio {
		if len(s.Interfaces) != 0 {
			var found bool

			for _, name := range s.Interfaces {
				if name == io.Name {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		} else if !s.skipChecks {
			iface, err := net.InterfaceByName(io.Name)
			if err != nil {
				continue
			}

			if iface.Flags&net.FlagLoopback == net.FlagLoopback {
				continue
			}

			if iface.Flags&net.FlagUp == 0 {
				continue
			}
		}

		tags := map[string]string{
			"interface": io.Name,
		}

		fields := map[string]interface{}{
			"bytes_sent":   io.BytesSent,
			"bytes_recv":   io.BytesRecv,
			"packets_sent": io.PacketsSent,
			"packets_recv": io.PacketsRecv,
			"err_in":       io.Errin,
			"err_out":      io.Errout,
			"drop_in":      io.Dropin,
			"drop_out":     io.Dropout,
		}
		acc.AddFields("net", fields, tags)
	}

	// Get system wide stats for different network protocols
	// (ignore these stats if the call fails)
	// netprotos, _ := s.ps.NetProto()
	// fields := make(map[string]interface{})
	// for _, proto := range netprotos {
	// 	for stat, value := range proto.Stats {
	// 		name := fmt.Sprintf("%s_%s", strings.ToLower(proto.Protocol),
	// 			strings.ToLower(stat))
	// 		fields[name] = value
	// 	}
	// }
	// tags := map[string]string{
	// 	"interface": "all",
	// }
	// acc.AddFields("net", fields, tags)

	return nil
}

func init() {
	agent.AddInput("net", &NetIOStats{ps: &systemPS{}})
}
