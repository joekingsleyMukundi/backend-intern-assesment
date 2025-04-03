package db

import (
	"context"
	"testing"
	"time"

	"github.com/joekingsleyMukundi/backend-intern-assesment/common/util"
	"github.com/stretchr/testify/require"
)

func CreatePayments(t *testing.T) Payment {
	user := CreateRandomUser(t)
	arg := CreatePaymentParams{
		Owner:  user.Username,
		Amount: util.RandomAmount(),
		Status: "Pending Payment",
	}
	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.Owner, payment.Owner)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.Status, payment.Status)

	require.NotEmpty(t, payment.CreatedAt)
	return payment
}

func TestCreatePaymet(t *testing.T) {
	CreatePayments(t)
}

func TestGetPayment(t *testing.T) {
	paymentCreated := CreatePayments(t)
	paymentGot, err := testQueries.GetPayment(context.Background(), paymentCreated.ID)
	require.NoError(t, err)
	require.NotEmpty(t, paymentGot)

	require.Equal(t, paymentCreated.Amount, paymentGot.Amount)
	require.Equal(t, paymentCreated.Status, paymentGot.Status)
	require.Equal(t, paymentCreated.Owner, paymentGot.Owner)

	require.WithinDuration(t, paymentCreated.CreatedAt, paymentGot.CreatedAt, time.Second)
}
func TestUpdatePayment(t *testing.T) {
	paymentCreated := CreatePayments(t)
	arg := UpdatePaymentParams{
		ID:     paymentCreated.ID,
		Status: "Completed",
	}
	updatedPayment, err := testQueries.UpdatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedPayment)

	require.Equal(t, paymentCreated.Amount, updatedPayment.Amount)
	require.Equal(t, arg.Status, updatedPayment.Status)
	require.Equal(t, paymentCreated.Owner, updatedPayment.Owner)

	require.WithinDuration(t, paymentCreated.CreatedAt, updatedPayment.CreatedAt, time.Second)
}
