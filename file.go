package main

import (
	"bufio"
	"os"
)

func WriteStringToFile(filepath, output string) error {
	/*	Open file in Write Mode	*/
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	/*	Truncate the size of file to 0 */
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	/*	Go to file initial position	*/
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	/* Initialize writer and write to file	*/
	writer := bufio.NewWriter(file)
	writer.WriteString(output)

	/*	Flush data do writer	*/
	return writer.Flush()
}
