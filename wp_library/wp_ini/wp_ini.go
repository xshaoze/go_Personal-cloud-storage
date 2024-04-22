package wp_ini

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	fileName       string
	fileLocation   string
	CfgSectionData map[string]map[string]string
	Config         map[string]string
}

var ConfigData Config

func init() {
	ConfigData, _ = newConfig("config.ini")
}

func newConfig(fileName string) (Config, error) {
	var cfg Config
	// 初始化目录
	path, _ := os.Getwd()
	cfg.fileName = fileName
	cfg.fileLocation = path + "\\" + cfg.fileName
	cfg.CfgSectionData = make(map[string]map[string]string)
	cfg.Config = make(map[string]string)
	// 检测文件是否存在
	file, err := os.Open(cfg.fileLocation)
	if err != nil {
		file, err = os.Create(cfg.fileLocation)
		if err != nil {
			return cfg, errors.New("文件创建失败：" + err.Error())
		}
		file.Close()
		return cfg, nil
	}
	defer file.Close()
	cfg.readCfg()

	return cfg, nil
}

func (cfg Config) readCfg() {
	file, err := os.Open(cfg.fileLocation)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var section string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		// 检测节(section)的开始
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(line[1 : len(line)-1])
			cfg.CfgSectionData[section] = make(map[string]string)
		} else if section == "" {
			section = "default"
			cfg.CfgSectionData[section] = make(map[string]string)
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				cfg.CfgSectionData[section][key] = value
			}
		} else if section != "" {
			// 解析选项和值
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				cfg.CfgSectionData[section][key] = value
			}
		}
	}
	for i := range cfg.CfgSectionData {
		for j := range cfg.CfgSectionData[i] {
			cfg.Config[j] = cfg.CfgSectionData[i][j]
		}
	}
}

func (cfg Config) SaveCFG() {
	file, err := os.OpenFile(cfg.fileLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
	}
	var lineStr string
	for key, value := range ConfigData.Config {
		lineStr = lineStr + key + "=" + value + "\n"
	}
	_, err = file.Write([]byte(lineStr))
	if err != nil {
		log.Println(err)
	}
	log.Println(lineStr)
}
