package links

import (
	"encoding/csv"
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"strconv"
)

func ConvertJSONToCSV(source, destination string, loggingUtil *zap.SugaredLogger, clientName string) error {
	// 2. Read the JSON file into the struct array
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	// remember to close the file at the end of the function
	defer func(sourceFile *os.File) {
		err = sourceFile.Close()
		if err != nil {
			loggingUtil.Errorw("error while closing the source file", zap.Error(err),
				"client", clientName)
		}
	}(sourceFile)

	var links = make(map[string]LinkStruct)
	if err = json.NewDecoder(sourceFile).Decode(&links); err != nil {
		loggingUtil.Errorw("error while decoding the json file", zap.Error(err),
			"client", clientName)
		return err
	}

	// 3. Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	if err != nil {
		loggingUtil.Errorw("error while creating the output file", zap.Error(err),
			"client", clientName)
		return err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			loggingUtil.Errorw("error while closing the output file", zap.Error(err),
				"client", clientName)
		}
	}(outputFile)

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"Link", "Broken Status"}
	if err = writer.Write(header); err != nil {
		return err
	}

	for _, r := range links {
		var csvRow []string
		csvRow = append(csvRow, r.Link, strconv.FormatBool(r.IsBroken))
		if err = writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
