package server

import (
	"database/sql"
	"log"
	"net/http"
	"omisescb/endpoint/scb"
	omise_client "omisescb/external/omise"
	"omisescb/repository"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const defaultPort = "8080"

func init() {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("[ERROR] [NewConfigReader] init failed with error : ", err.Error())
		panic(err.Error())
	}
}

func Run() {
	handler, err := setUpServer()
	if err != nil {
		log.Printf("Error setUpServer ---> %+v", err.Error())
	}

	srv := &http.Server{Addr: ":" + defaultPort, Handler: handler}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Printf("Error Server shut down. Waiting for connections to drain. --> %+v", err)
			} else {
				log.Printf("Error failed to start server on port: %+v", srv.Addr)
			}
		}
	}()

	// Wait for an interrupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	<-sigint
}
func setUpServer() (http.Handler, error) {
	db, err := SetUpDB()
	if err != nil {
		return nil, err
	}

	omise, err := omise_client.SetupOmise()
	if err != nil {
		return nil, err
	}

	//new repository
	orderRepo := repository.NewOrderTransaction(db)

	//new service
	SCBService := scb.NewSCBPayment(omise, orderRepo)

	router := mux.NewRouter()
	router.Handle("/omise-scb", scb.MakeCreateSCBPaymentHandler(SCBService)).Methods(http.MethodPost)
	router.Handle("/omise-scb", scb.MakeGetSCBPaymentHandler(SCBService)).Methods(http.MethodGet)
	router.Handle("/omise-scb/{source_id}/complete", scb.MakeGetCallBackSCBPaymentHandler(SCBService)).Methods(http.MethodGet)
	return router, nil
}

func SetUpDB() (*sql.DB, error) {

	database, err := sql.Open("sqlite3", "db/order.db")
	if err != nil {
		return database, err
	}
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS orders (id INTEGER PRIMARY KEY, order_id INTEGER, charge_id VARCHAR,amount INTEGER,currency VARCHAR,payment_status VARCHAR,soure_id VARCHAR,soure_type VARCHAR,paid_at TIMESTAMP,create_at TIMESTAMP)")
	if err != nil {
		return database, err
	}
	if _, err := statement.Exec(); err != nil {
		return database, err
	}
	return database, err
}
