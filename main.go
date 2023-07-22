package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Action string

const (
	Separator            = ","
	ActionEncrypt Action = "encrypt"
	ActionDecrypt Action = "decrypt"
)

var (
	buildVersion string
	buildTime    string

	action          string
	sensitiveFields string
	curJson         string
	lastJson        string
	filePath        string
)

func main() {
	flag.StringVar(&sensitiveFields, "fd", "", "sensitive fields")
	flag.StringVar(&action, "action", "", "encrypt or decrypt")
	flag.StringVar(&curJson, "new", "", "new config json")
	flag.StringVar(&lastJson, "old", "", "old config json")
	flag.StringVar(&filePath, "out", "", "out file name")
	flag.Parse()

	if action != string(ActionDecrypt) && action != string(ActionEncrypt) {
		fmt.Printf("BuildVersion: %s, BuildTime: %s\n", buildVersion, buildTime)
		flag.Usage()
		os.Exit(1)
	}

	task := NewJson(Action(action))

	err := task.ProcessMap(task.MergeMap(curJson, lastJson), filePath)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

type Json interface {
	MergeMap(cur string, last ...string) map[string]interface{}
	ProcessMap(data map[string]interface{}, filePath ...string) error
	Err() error
}

func NewJson(action Action) Json {
	if action == ActionDecrypt {
		return &DecryptMap{}
	}

	return &EncryptMap{}
}

type DecryptMap struct {
	err error
}

// MergeMap just return cur
func (d *DecryptMap) MergeMap(cur string, _ ...string) map[string]interface{} {
	var curMap map[string]interface{}

	err := json.Unmarshal([]byte(cur), &curMap)
	if err != nil {
		d.err = fmt.Errorf("decrypt failed marshall, err: %w", err)
		return nil
	}

	return curMap
}

func (d *DecryptMap) Err() error {
	return d.err
}

// ProcessMap process value if key in sensitiveFields, print to stdout
func (d *DecryptMap) ProcessMap(data map[string]interface{}, filePath ...string) error {
	if d.err != nil {
		return d.err
	}

	processedMap, err := processMap(ActionDecrypt)(data, sensitiveFields)
	if err != nil {
		return fmt.Errorf("process failed, err: %w", err)
	}

	return outputMap(processedMap)
}

type EncryptMap struct {
	err error
}

// MergeMap merge last some key-value to cur only when cur miss this key
func (e *EncryptMap) MergeMap(cur string, last ...string) map[string]interface{} {
	if len(last) == 0 {
		e.err = errors.New("merge miss src")
		return nil
	}

	var (
		curMap  map[string]interface{}
		lastMap map[string]interface{}
	)

	err := json.Unmarshal([]byte(cur), &curMap)
	if err != nil {
		e.err = fmt.Errorf("failed unmarshall val: %s, err: %w", string(cur), err)
		return nil
	}

	err = json.Unmarshal([]byte(last[0]), &lastMap)
	if err != nil {
		e.err = fmt.Errorf("failed unmarshall val: %s, err: %w", string(cur), err)
		return nil
	}

	for lackkey, val := range lastMap {
		if _, ok := curMap[lackkey]; !ok { // last +> cur
			curMap[lackkey] = val
		}
	}

	return curMap
}

// ProcessMap process value if key in sensitiveFields, save result to file
func (e *EncryptMap) ProcessMap(data map[string]interface{}, filePath ...string) error {
	if e.err != nil {
		return e.err
	}

	if len(filePath) == 0 || filePath[0] == "" {
		return fmt.Errorf("process failed, err: miss fielPath")
	}

	processedMap, err := processMap(ActionEncrypt)(data, sensitiveFields)
	if err != nil {
		return fmt.Errorf("process failed, err: %w", err)
	}

	return writeMap(filePath[0], processedMap)
}

func (d *EncryptMap) Err() error {
	return d.err
}

// processMap decrypt or encrypt value if key in fields
func processMap(action Action) func(configMap map[string]interface{}, fields string) (map[string]interface{}, error) {
	sensitiveFields = strings.TrimSpace(sensitiveFields)

	return func(configMap map[string]interface{}, fields string) (map[string]interface{}, error) {
		if sensitiveFields == "" {
			return configMap, nil
		}

		for _, field := range strings.Split(sensitiveFields, Separator) {
			val, ok := configMap[field]
			if !ok {
				continue
			}

			valStr, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("failed, err: action: %s feild: %s invalid not str", action, field)
			}

			s, err := execCmd("echo", valStr+"-"+string(action))
			if err != nil || strings.Contains(s, "error Exception") {
				return nil, fmt.Errorf("failed, action: %s feild: %s valStr: %s, %s err: %w", action, field, valStr, s, err)
			}

			configMap[field] = s
		}

		return configMap, nil
	}
}

// execCmd exec shell command
func execCmd(name string, arg ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)
	output, err := cmd.Output()
	if err != nil {
		return string(output), err
	}

	return strings.Replace(string(output), "\n", "", -1), nil
}

// writeMap 将Map写文件
func writeMap(filePath string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(b)

	return nil
}

// outputMap print to stdout
func outputMap(content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	fmt.Printf("%s", b)

	return nil
}
