package performance

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"path/filepath"
	"time"
)

type TestTS struct {
	LTime time.Time
	TTime time.Time
}

const (
	LTimeLayout = "15:04:05.999999"
	TTimeLayout = "15:04:05.999"
)

var (
	// Cmd is the quote command.
	Cmd = &cobra.Command{
		Use:     "performance",
		Short:   "performance testing",
		Long:    "performance testing",
		Aliases: []string{"pt"},
		Example: "qf performance marketdata.log",
		RunE:    execute,
	}
)

func execute(cmd *cobra.Command, args []string) error {
	argLen := len(args)
	if argLen != 1 {
		return fmt.Errorf("incorrect argument number")
	}

	testFileArg := args[0]
	// open the file
	dir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error:", err)
		return err
    }
	
	testFile, err := os.Open( filepath.Join( dir,"cmd/performance/", testFileArg))
	if err != nil {
		return fmt.Errorf("error opening %v, %v", testFile, err)
	}
	defer testFile.Close()

	// create a scanner to read the file
	scanner := bufio.NewScanner(testFile)

	//read the file line by line
	times := make([]TestTS,0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "35=W") {
			sub1 := strings.Split(line, " ")
			if len(sub1) > 2 {
				localTime := sub1[1]

				parsedLTime, err := time.Parse(LTimeLayout, localTime)
				if err != nil {
					fmt.Println("Error parsing time:", err)
					return err
				}

				//fmt.Printf("localTime:%s", localTime)

				sub2 := strings.Split(sub1[2], "\u0001")
				if len(sub2) > 5 {
					if strings.Contains(sub2[5], "52=") {
						sub3 := strings.Split(sub2[5], "-")
						if len(sub3) > 1 {
							talosTime := sub3[1]
							pasrsedTTime,err := time.Parse(TTimeLayout, talosTime)
							if err != nil {
								fmt.Println("Error parsing time:", err)
								return err
							}

							testTS := TestTS{
								LTime : parsedLTime,
								TTime : pasrsedTTime,
							}
							times = append(times, testTS)

							//fmt.Printf("times:%s\n", times)
						}
					}
				}
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	//for i :=1; i < len(times); i++ {

	//}

	
	for i :=1; i < len(times); i++ {
		t := times[i]
		t0 := times[i-1]

		tAtoT :=  t.LTime.Sub(t.TTime)
		interv := t.TTime.Sub(t0.TTime)
		fmt.Printf("i=%v, interv=%v, T->A=%v\n",i, interv,tAtoT)
		
	}

	return nil
}
