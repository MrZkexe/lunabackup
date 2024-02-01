package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Configs struct {
	DB           bool     `json:"db"`
	Folders      []string `json:"folders"`
	UserMariaDB  string   `json:"user_mariadb"`
	PassMariaDB  string   `json:"pass_mariadb"`
}

var (
	configFilePath = "lunaconf.json"
	baseFolder     = "/etc/LunaBackup"
	backupPath     = "/backup"
)

func startBackup(backupPath string, configData []byte, dateTime string) {
	var config Configs
	var fileList []string

	json.Unmarshal(configData, &config)

	for _, folder := range config.Folders {
		fmt.Println("Compiling folder →", folder)
		fileList = append(fileList, getAllFilesInFolder(folder)...)
	}

	if config.DB {
		fmt.Println("Compiling → database")
		dbBackupCommand := exec.Command("mariadb-dump", fmt.Sprintf("-u%s", config.UserMariaDB), fmt.Sprintf("-p%s", config.PassMariaDB), "--all-databases")
		dbOutput, err := dbBackupCommand.CombinedOutput()
		if err != nil {
			fmt.Println("Error running mariadb-dump:", err)
			return
		}
		err = ioutil.WriteFile(fmt.Sprintf("/db-%s.sql", dateTime), dbOutput, 0644)
		if err != nil {
			fmt.Println("Error writing database backup file:", err)
			return
		}
		fileList = append(fileList, fmt.Sprintf("/db-%s.sql", dateTime))
	}

	fmt.Println("Compilation started")

	tarGzFile := fmt.Sprintf("%s/bkp-%s.tar.gz", backupPath, dateTime)
	err := createTarGz(tarGzFile, fileList)
	if err != nil {
		fmt.Println("Error creating tar.gz file:", err)
		return
	}

	fmt.Printf("File compiled → %s\n", tarGzFile)
	err = os.Remove(fmt.Sprintf("/db-%s.sql", dateTime))
	if err != nil {
		fmt.Println("Error removing temporary database backup file:", err)
	}
}

func getAllFilesInFolder(folderPath string) []string {
	var fileList []string
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info.Mode().IsRegular() {
				fileList = append(fileList, path)
			}
		}
		return nil
	})
	return fileList
}

func createTarGz(tarGzFile string, fileList []string) error {
	file, err := os.Create(tarGzFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range fileList {
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", filePath, err)
			continue
		}

		relFilePath, err := filepath.Rel("/", filePath)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %s\n", filePath, err)
			continue
		}

		tarHeader := &tar.Header{
			Name: relFilePath,
			Mode: 0644,
			Size: int64(len(fileData)),
		}

		err = tarWriter.WriteHeader(tarHeader)
		if err != nil {
			fmt.Printf("Error writing tar header for %s: %s\n", filePath, err)
			continue
		}

		_, err = tarWriter.Write(fileData)
		if err != nil {
			fmt.Printf("Error writing file data for %s: %s\n", filePath, err)
		}
	}

	return nil
}

func createBackupFolder(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		fmt.Println("Creating folder", folder)
		os.MkdirAll(folder, os.ModePerm)
	}
	fmt.Println("Folder", folder, "is OK")
}

func verifyJSON(configFilePath string) (bool, error) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		os.MkdirAll(baseFolder, os.ModePerm)
		fmt.Println("File does not exist\nCreating file")

		config := &Configs{
			DB:           false,
			Folders:      []string{"/etc", "/home", "/var/log", "/var/www", "/usr/local/bin", "/usr/local/sbin", "/var/spool/cron"},
			UserMariaDB:  "root",
			PassMariaDB:  "root",
		}
		jsonData, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			return false, fmt.Errorf("Error marshalling default config: %s", err)
		}
		err = ioutil.WriteFile(configFilePath, jsonData, 0640)
		if err != nil {
			return false, fmt.Errorf("Error writing default config file: %s", err)
		}
		fmt.Println("File created.")
	} else {
		fmt.Println("File exists")
	}
	return true, nil
}

func main() {
	fmt.Println("Checking:", filepath.Join(baseFolder, configFilePath))
	resp, err := verifyJSON(filepath.Join(baseFolder, configFilePath))
	switch resp {
	case true:
		read, err := ioutil.ReadFile(filepath.Join(baseFolder, configFilePath))
		if err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}
		fmt.Println(string(read))
		fmt.Println("Validating backup folder")
		createBackupFolder(backupPath)
		fmt.Println("Starting backup")
		startBackup(backupPath, read, time.Now().Format("02-01-2006"))
	default:
		fmt.Println(err)
		return
	}
}
