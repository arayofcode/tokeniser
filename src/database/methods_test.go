package database

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/arayofcode/tokeniser/src/models"
	"github.com/google/uuid"

	"github.com/arayofcode/tokeniser/src/common"
)

var (
	ctx context.Context
	db  Database
)

func configInit(t *testing.T) {
	t.Log("Setting up the database for testing")
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	db = DatabaseInit(common.GetDbURL())
}

func TestInsert(t *testing.T) {
	type testData struct {
		input         models.CreditCardDetails
		expectedError bool
		want          models.CreditCardRow
	}

	configInit(t)
	tests := map[string]testData{
		"empty name": {
			input: models.CreditCardDetails{
				CardHolderName:      "",
				CardNumberEncrypted: []byte("a-very-valid-card"),
				ExpiryDateEncrypted: []byte("a-very-valid-expiry"),
			},
			expectedError: true,
			want:          models.CreditCardRow{},
		},
		"empty card number": {
			input: models.CreditCardDetails{
				CardHolderName:      "Tamato Rolli",
				CardNumberEncrypted: []byte(""),
				ExpiryDateEncrypted: []byte("a-very-valid-expiry"),
			},
			expectedError: true,
			want:          models.CreditCardRow{},
		},
		"empty expiry": {
			input: models.CreditCardDetails{
				CardHolderName:      "Tamato Rolli",
				CardNumberEncrypted: []byte("a-very-valid-card"),
				ExpiryDateEncrypted: []byte(""),
			},
			expectedError: true,
			want:          models.CreditCardRow{},
		},
		"no data": {
			input:         models.CreditCardDetails{},
			expectedError: true,
			want:          models.CreditCardRow{},
		},
		"valid data": {
			input: models.CreditCardDetails{
				CardHolderName:      "Tamato Rolli",
				CardNumberEncrypted: []byte("a-very-valid-card"),
				ExpiryDateEncrypted: []byte("a-very-valid-expiry"),
			},
			expectedError: false,
			want:          models.CreditCardRow{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := db.InsertCard(ctx, tc.input)

			// Unexpected Error
			if err != nil && !tc.expectedError {
				t.Errorf("database.InsertCard() - %s: expected no error, found: %s", name, err)
				return
			}

			// Wanted error, didn't get any
			if err == nil && tc.expectedError {
				t.Errorf("database.InsertCard() - %s: expected error, found none", name)
				return
			}

			// Invalid UUID
			if err = uuid.Validate(got.Token.String()); err != nil {
				t.Errorf("database.InsertCard() - %s: expected valid UUID, found error: %s", name, err)
				return
			}

			// Success, Performing cleanups if card was created
			t.Logf("database.InsertCard() - %s: Test successful, performing cleanups if needed", name)
			if got.Token != uuid.Nil && !db.DeleteCard(ctx, got.Token) {
				t.Errorf("database.InsertCard()+database.DeleteCard() - %s: successful cleanup, failed while deleting", name)
				return
			}
		})
	}
}

func TestGet(t *testing.T) {
	type testData struct {
		input         uuid.UUID
		expectedError bool
		want          models.CreditCardRow
		createRow     bool
		createRowData models.CreditCardDetails
	}
	configInit(t)
	tests := map[string]testData{
		"invalid token": {
			input:         [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedError: true,
		},
		"empty token": {
			expectedError: true,
		},
		"token not in table": {
			input: uuid.Nil,
			want:  models.CreditCardRow{},
		},
		"token present in table": {
			createRow: true,
			createRowData: models.CreditCardDetails{
				CardHolderName:      "Tamato Present in Table",
				CardNumberEncrypted: []byte("a-very-valid-card"),
				ExpiryDateEncrypted: []byte("a-very-valid-expiry"),
			},
			want: models.CreditCardRow{
				CardHolderName:      "Tamato Present in Table",
				CardNumberEncrypted: []byte("a-very-valid-card"),
				ExpiryDateEncrypted: []byte("a-very-valid-expiry"),
			},
		},
	}

	for name, tc := range tests {
		t.Logf("database.GetCard() - %s: Setting rows in DB if needed before running test", name)
		if tc.createRow {
			want, _ := db.InsertCard(ctx, tc.createRowData)
			tc.input = want.Token
			tc.want.Token = want.Token
			tc.want.CreatedAt = want.CreatedAt
			tc.want.UpdatedAt = want.UpdatedAt
			tc.want.RowID = want.RowID
			tests[name] = tc
		}

		got, err := db.GetCard(ctx, tc.input)
		// Expected no error, received error
		if err != nil && !tc.expectedError {
			t.Errorf("database.GetCard() - %s: expected no error, found: %s", name, err)
			return
		}

		// When received and expected are different
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("database.GetCard() - %s: expected %+v \nfound: %+v", name, tc.want, got)
			return
		}

		//	Add more corner cases to check later

		// Success, Performing cleanups if card was created
		t.Logf("database.GetCard() - %s: Test successful, performing cleanups if needed", name)
		if got.Token != uuid.Nil && !db.DeleteCard(ctx, got.Token) {
			t.Errorf("database.GetCard()+database.DeleteCard() - %s: successful cleanup, failed while deleting", name)
			return
		}
	}
}
