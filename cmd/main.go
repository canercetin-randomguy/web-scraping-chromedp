package main

import (
	"canercetin/pkg/backend"
	"canercetin/pkg/logger"
	"canercetin/pkg/sqlpkg"
	"go.uber.org/zap"
	"log"
)

func main() {
	dbLogFile := logger.CreateNewFile("./logs/db")
	dbLogger, err := logger.NewLoggerWithFile(dbLogFile)
	if err != nil {
		log.Println(err)
	}
	// get a fresh database connection
	dbConn := sqlpkg.SqlConn{}
	err = dbConn.GetSQLConn("")
	if err != nil {
		dbLogger.Error(err.Error())
	}
	go func() {
		err = dbConn.CreateClientSchema()
		if err != nil {
			dbLogger.Error("Error creating client schema.", zap.Error(err))
		}
		err = dbConn.CreateClientTable()
		if err != nil {
			dbLogger.Errorw("Error creating client table.", zap.Error(err))
		}
		err = dbConn.CreateClientFileTable()
		if err != nil {
			dbLogger.Errorw("Error creating client file table.", zap.Error(err))
		}
	}()
	err = backend.StartWebPageBackend(backend.Port)
	if err != nil {
		log.Println(err)
	}
}
