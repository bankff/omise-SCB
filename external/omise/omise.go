package omise_client

import (
	"fmt"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/spf13/viper"
)

type omiseEx struct {
	client *omise.Client
}

type OmiseEx interface {
	CreateSoure(amount int64, currency string, peymentType string) (*omise.Source, error)
	CreateCharge(amount int64, currency string, source string, ID int64) (*omise.Charge, error)
	GetCharge(ID string) (*omise.Charge, error)
}

func SetupOmise() (OmiseEx, error) {
	publicKey := viper.GetString("Omise.public-key")
	privateKey := viper.GetString("Omise.private-key")

	client, err := omise.NewClient(publicKey, privateKey)
	if err != nil {
		return omiseEx{}, err
	}

	return omiseEx{
		client: client,
	}, nil
}

func (o omiseEx) CreateSoure(amount int64, currency string, paymentType string) (*omise.Source, error) {

	source, createSource := &omise.Source{}, &operations.CreateSource{
		Amount:   amount,
		Currency: currency,
		Type:     paymentType,
	}

	if err := o.client.Do(source, createSource); err != nil {
		return source, err
	}

	return source, nil
}

func (o omiseEx) CreateCharge(amount int64, currency string, source string, ID int64) (*omise.Charge, error) {

	charge, create := &omise.Charge{}, &operations.CreateCharge{
		Amount:    amount,
		Currency:  currency,
		Source:    source,
		ReturnURI: fmt.Sprintf("http://localhost:8080/omise-scb/%d/complete", ID),
	}

	if err := o.client.Do(charge, create); err != nil {
		return charge, err
	}

	return charge, nil
}

func (o omiseEx) GetCharge(ID string) (*omise.Charge, error) {
	charge, retrieve := &omise.Charge{}, &operations.RetrieveCharge{ID}
	if err := o.client.Do(charge, retrieve); err != nil {
		return charge, err
	}

	return charge, nil

}
