package main

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"time"
)

var list []int

func main() {

	start := time.Now()

	params := g.GoSNMP{
		Target:    "172.16.8.2",
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(1) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	err1 := params.Connect()
	if err1 != nil {
		fmt.Println("Unable to connect")
	}

	var oidList []string
	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.2.")
	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.3.")
	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.7.")
	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.8.")

	var data = make(map[int]interface{})

	for k := 0; k < 1; k++ {

		err := params.Walk(".1.3.6.1.2.1.2.2.1.1", walkFunc)
		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < len(list); i++ {

			var newList []string
			for j := 0; j < len(oidList); j++ {
				newList = append(newList, oidList[j]+strconv.Itoa(list[i]))
			}

			ans, _ := params.Get(newList)

			var value []string

			for _, result := range ans.Variables {

				switch result.Type {

				case g.Integer:
					value = append(value, fmt.Sprintf("%d", result.Value))

				case g.OctetString:
					value = append(value, string(result.Value.([]byte)))

				case g.IPAddress:
					value = append(value, result.Value.(string))

				default:
					value = append(value, result.Value.(string))

				}

			}

			newList = newList[:0]

			data[list[i]] = value

			value = nil

		}

	}

	end := time.Now()

	fmt.Println(end.Sub(start))

}

func walkFunc(pdu g.SnmpPDU) error {

	list = append(list, pdu.Value.(int))

	return nil
}
