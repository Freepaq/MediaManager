package MediaUtils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func fileStruct(files *[]string, origin string, eligibleFiles string) filepath.WalkFunc {
	fmt.Println("Start scanning folder :" + origin)
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() && origin != path {
			//		filepath.Walk(path, fileStruct(files, path))
		} else {
			if VIDEO == eligibleFiles {
				if IsVideoEligible(filepath.Ext(path)) {
					*files = append(*files, path)
				}
			}
			if PHOTO == eligibleFiles {
				if IsPhotoEligible(filepath.Ext(path)) {
					*files = append(*files, path)
				}
			}
			if ALL == eligibleFiles {
				if IsMediEligible(filepath.Ext(path)) {
					*files = append(*files, path)
				}
			}
		}
		return nil
	}
}
func GetListOfFile(folder string, eligibleFiles string) []string {
	var files []string
	err := filepath.Walk(folder, fileStruct(&files, folder, eligibleFiles))
	if err != nil {
		panic(err)
	}

	return files
}

func Delete(file FileStruct) {
	err := os.Remove(file.FullName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File [" + file.FullName + "] moved")
	}
}

func Rename(file *FileStruct) {
	var n = (*file).CreationDate.Format(TimestampFormat)
	(*file).NewName = n + (*file).Extension
	(*file).NewFullName = (*file).DestinationDir + (*file).NewName
	fmt.Println("New name : [" + (*file).NewName + "]")
}

func Copy(ori *FileStruct, destFoler string, force bool) bool {
	year, month, _ := (*ori).CreationDate.Date()
	dest := filepath.Join(destFoler, strconv.Itoa(year))
	destMonth := filepath.Join(destFoler, strconv.Itoa(year), Months[int(month)])
	destFull := filepath.Join(destMonth, (*ori).NewName)
	(*ori).DestinationDir = destMonth
	(*ori).NewFullName = filepath.Join((*ori).DestinationDir, (*ori).NewName)
	if _, err := os.Stat(destFoler); os.IsNotExist(err) {
		os.Mkdir(destFoler, 0755)
	}
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		os.Mkdir(dest, 0755)
	}
	if _, err := os.Stat(destMonth); os.IsNotExist(err) {
		os.Mkdir(destMonth, 0755)
	}
	input, err := ioutil.ReadFile((*ori).FullName)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var result = true
	(*ori).Proccessed = true
	if !force {
		if _, err := os.Stat(destFull); os.IsNotExist(err) {
			if err, result = writeFile(destFull, input); err != nil {
				fmt.Println(err)
				(*ori).Proccessed = false
			}
		} else {
			fmt.Println("Destination file : [" + destFull + "] exists, not overrided")
			(*ori).Proccessed = false
		}
		(*ori).Proccessed = true
	} else {
		oriFile, err := ioutil.ReadFile((*ori).FullName)
		if nil != err {
			fmt.Println("Error reading  origin file :" + (*ori).FullName + " ")
		}
		destFile, err := ioutil.ReadFile(destFull)
		if nil != err {
			fmt.Println("Error reading destination file :" + destFull + " ")
		}
		if len(oriFile) == len(destFile) {
			if err, result = writeFile(destFull, input); err != nil {
				fmt.Println(err)
				(*ori).Proccessed = false
			}
			(*ori).Proccessed = true
		}
		//TODO if force == true and dest file same origin file
	}
	input = nil
	return result
}

func writeFile(destFull string, input []byte) (error, bool) {
	err := ioutil.WriteFile(destFull, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destFull)
		fmt.Println(err)
		return nil, false
	}
	fmt.Println("Destination file : [" + destFull + "] copied")
	return err, true
}

func GetMeta(fname string) (FileStruct, error) {
	fileStr := FileStruct{}

	fileStr.FullName = fname
	fileStr.NewFullName = fname

	if IsVideoEligible(filepath.Ext(fname)) {
		_ = ReadVideoMeta(fname, &fileStr)
		fileStr.TypeOfMedia = VIDEO
	}
	if IsPhotoEligible(filepath.Ext(fname)) {
		ReadPhotoMeta(fname, &fileStr)
		fileStr.TypeOfMedia = PHOTO
	}

	fileStr.OriginDir, fileStr.Name = filepath.Split(fname)
	fileStr.Extension = filepath.Ext(fname)
	fileStr.NewName = fileStr.Name
	return fileStr, nil
}