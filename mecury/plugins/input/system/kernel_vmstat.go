// +build linux

package system

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/aiyun/openapm/mecury/agent"
)

type KernelVmstat struct {
	statFile string
}

func (k *KernelVmstat) Description() string {
	return "Get kernel statistics from /proc/vmstat"
}

func (k *KernelVmstat) SampleConfig() string {
	return ""
}

func (k *KernelVmstat) Gather(acc agent.Accumulator) error {
	data, err := k.getProcVmstat()
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})

	dataFields := bytes.Fields(data)
	for i, field := range dataFields {

		// dataFields is an array of {"stat1_name", "stat1_value", "stat2_name",
		// "stat2_value", ...}
		// We only want the even number index as that contain the stat name.
		if i%2 == 0 {
			// Convert the stat value into an integer.
			m, err := strconv.Atoi(string(dataFields[i+1]))
			if err != nil {
				return err
			}

			fields[string(field)] = int64(m)
		}
	}

	acc.AddFields("kernel_vmstat", fields, map[string]string{})
	return nil
}

func (k *KernelVmstat) getProcVmstat() ([]byte, error) {
	if _, err := os.Stat(k.statFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("kernel_vmstat: %s does not exist!", k.statFile)
	} else if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(k.statFile)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func init() {
	agent.AddInput("kernel_vmstat",
		&KernelVmstat{
			statFile: "/proc/vmstat",
		},
	)
}
