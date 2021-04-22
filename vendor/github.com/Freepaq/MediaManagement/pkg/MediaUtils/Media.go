package MediaUtils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

/*
Get a list of of media files including sub-folders search
*/
func GetListOfFile(folder string, eligibleFiles string) []string {
	var files []string
	err := filepath.Walk(folder, fileStruct(&files, folder, eligibleFiles))
	if err != nil {
		panic(err)
	}

	return files
}

/*
Delete physically the file given as parameters (fullpath)
*/
func Delete(file FileStruct) {
	err := os.Remove(file.FullName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File [" + file.FullName + "] moved")
	}
}

/*
Rename physically the file given as parameters (fullpath)
Rename is based on the Creation date set on FileStruct return by func GetMeta
*/
func Rename(file *FileStruct) {
	var n = (*file).CreationDate.Format(TimestampFormat)
	(*file).NewName = n + (*file).Extension
	(*file).NewFullName = (*file).DestinationDir + (*file).NewName
	fmt.Println("New name : [" + (*file).NewName + "]")
}

/*
Copy physically the file from origin folder to destfolder.
In Dest folder, a sub folder corresponding to the year of media is created if not exists
In Dest year folder, a sub folder corresponding to the month of media is created if not exists
*/
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
		//(*ori).Proccessed = true
	} else {
		oriFile, err := ioutil.ReadFile((*ori).FullName)
		if nil != err {
			fmt.Println("Error reading  origin file :" + (*ori).FullName + " ")
		}
		if _, err := os.Stat(destFull); os.IsNotExist(err) {
			if err, result = writeFile(destFull, input); err != nil {
				fmt.Println(err)
				(*ori).Proccessed = false

			} else {
				(*ori).Proccessed = true
			}

		} else {
			destFile, err := ioutil.ReadFile(destFull)
			index := 1
			if nil != err {
				fmt.Println("Error reading destination file :" + destFull + " ")
				(*ori).Proccessed = false
			} else {
				if len(oriFile) != len(destFile) && (*ori).NewFullName == destFull {
					if err, result = writeFileLooping(destFull, input, index, ori); err != nil {
						fmt.Println(err)
						(*ori).Proccessed = false
					} else {
						(*ori).Proccessed = true
						fmt.Println("New File created as destination file exists with other size")
						fmt.Println((*ori).NewFullName)
					}
				} else {
					fmt.Println("Destination already exists with same size and name, no copy")
					(*ori).Proccessed = true
				}
			}
		}
	}
	input = nil
	return result
}

func writeFileLooping(destFull string, input []byte, index int, f *FileStruct) (error, bool) {
	s := strings.Split(destFull, ".")
	s[0] = s[0] + "(" + strconv.Itoa(index) + ")"
	f.NewFullName = s[0] + "." + s[1]
	if err, b := writeFile(f.NewFullName, input); !b {
		if index > 15 {
			return err, false
		}
		index++
		if err, b = writeFileLooping(destFull, input, index, f); !b {
			return err, b
		}
	}
	return nil, true
}

func writeFile(destFull string, input []byte) (error, bool) {
	err := ioutil.WriteFile(destFull, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destFull)
		fmt.Println(err)
		return err, false
	}
	fmt.Println("Destination file : [" + destFull + "] copied")
	return nil, true
}

/*
function to get the meta date of the media
Return a FileStruct
*/

func GetMeta(fname string) (FileStruct, error) {
	fileStr := FileStruct{}

	fileStr.FullName = fname
	fileStr.NewFullName = fname

	if IsVideoEligible(filepath.Ext(fname)) {
		fileStr.TypeOfMedia = VIDEO
		if err := ReadVideoMeta(fname, &fileStr); err != nil {
			return fileStr, err
		}
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
