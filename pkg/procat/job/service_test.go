package job

import (
	"front-office/common/constant"
	"front-office/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapLoanRecordCheckerRow(t *testing.T) {
	t.Run("should map all fields correctly", func(t *testing.T) {
		message := "Succeed"
		result := mapLoanRecordCheckerRow(&logTransProductCatalog{
			Input: &logTransInput{
				Name:        helper.StringPtr(constant.DummyName),
				NIK:         helper.StringPtr(constant.DummyNIK),
				PhoneNumber: helper.StringPtr(constant.DummyPhoneNumber),
			},
			Data: &logTransData{
				Remarks: helper.StringPtr("-"),
				Status:  helper.StringPtr(""),
			},
			Message: &message,
		})

		expected := []string{constant.DummyName, constant.DummyNIK, constant.DummyPhoneNumber, "-", "", "Succeed"}
		assert.Equal(t, expected, result)
	})
}
