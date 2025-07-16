package job

import (
	"front-office/common/constant"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapLoanRecordCheckerRow(t *testing.T) {
	t.Run("should map all fields correctly", func(t *testing.T) {
		message := "Succeed"
		result := mapLoanRecordCheckerRow(&logTransProductCatalog{
			Input:   &logTransInput{},
			Data:    &logTransData{},
			Message: &message,
		})

		expected := []string{constant.DummyName, constant.DummyNIK, constant.DummyPhoneNumber, "-", "", "Succeed"}
		assert.Equal(t, expected, result)
	})
}
