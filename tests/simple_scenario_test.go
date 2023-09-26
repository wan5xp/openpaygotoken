package openpaygotoken_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/wan5xp/openpaygotoken/pkg/openpaygotoken"
	"github.com/wan5xp/openpaygotoken/pkg/simulators"
)

func TestSimpleScenario(t *testing.T) {
	fmt.Println("Device: We initiate the device simulator with our device")
	deviceSimulator, err := simulators.NewDeviceSimulator(startingCode, &key, 1, false, false, 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Server: We initiate the server simulator with our device")
	serverSimulator := simulators.NewSingleDeviceServerSimulator(startingCode, &key, 1, false, 1)

	fmt.Println("\nDevice: We try entering an invalid token into the device: 123456789")
	err = deviceSimulator.EnterToken("123456789")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be still inactive)")
	deviceSimulator.PrintStatus()

	fmt.Println("\nServer: We add 1 days of activation for the device")
	thisToken, err := serverSimulator.GenerateTokenFromDate(time.Now().Add(24*1*time.Hour), false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token:", thisToken)
	fmt.Println("Device: We enter the generated token into the device")
	err = deviceSimulator.EnterToken(thisToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 1 day)")
	deviceSimulator.PrintStatus()
	fmt.Println("Device: We check the device status (should be active)")
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < (24*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (24*time.Hour+1*time.Second) {
		t.Errorf("Expected expiration timestamp to be %s, got %s", serverSimulator.ExpirationDate, deviceSimulator.ExpirationTimestamp)
	}

	fmt.Println("Device: We enter the token a second time to make sure it doesnt add the days again")
	err = deviceSimulator.EnterToken(thisToken)
	if err != nil && !errors.Is(err, &simulators.ErrOldToken{}) {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 1 day)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if deviceSimulator.ExpirationTimestamp.Sub(serverSimulator.ExpirationDate) > 1*time.Second {
		t.Errorf("Expected expiration timestamp to be %s, got %s", serverSimulator.ExpirationDate, deviceSimulator.ExpirationTimestamp)
	}

	fmt.Println("\nServer: We set it to expire in 30 days")
	thisToken, err = serverSimulator.GenerateTokenFromDate(time.Now().Add(30*24*time.Hour), false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token:", thisToken)
	fmt.Println("Device: We enter the generated token into the device")
	err = deviceSimulator.EnterToken(thisToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 30 days)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < (30*24*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (30*24*time.Hour+1*time.Second) {
		t.Errorf("Expected expiration timestamp to be %s, got %s", serverSimulator.ExpirationDate, deviceSimulator.ExpirationTimestamp)
	}

	fmt.Println("\nServer: We set it to expire in 7 days (removing 23 days)")
	thisToken, err = serverSimulator.GenerateTokenFromDate(time.Now().Add(7*24*time.Hour), false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token:", thisToken)
	fmt.Println("Device: We enter the generated token into the device")
	err = deviceSimulator.EnterToken(thisToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 7 days)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < (7*24*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (7*24*time.Hour+1*time.Second) {
		t.Errorf("Expected expiration timestamp to be %s, got %s", serverSimulator.ExpirationDate, deviceSimulator.ExpirationTimestamp)
	}

	fmt.Println("\nServer: We generate a token for putting the device in PAYG-OFF mode")
	thisPaygOffCode, err := serverSimulator.GeneratePaygDisableToken()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token:", thisPaygOffCode)
	fmt.Println("Device: We enter the generated token into the device")
	err = deviceSimulator.EnterToken(thisPaygOffCode)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active forever)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if deviceSimulator.PaygEnabled {
		t.Errorf("Expected PAYG to be disabled")
	}

	fmt.Println("\nServer: We generate a token for putting the device in PAYG-ON mode with 0 days")
	thisToken, err = serverSimulator.GenerateTokenFromDate(time.Now(), false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token:", thisToken)
	fmt.Println("Device: We enter the generated token into the device")
	err = deviceSimulator.EnterToken(thisToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should not be active)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if !deviceSimulator.PaygEnabled {
		t.Errorf("Expected PAYG to be enabled")
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < -1*time.Second || time.Until(deviceSimulator.ExpirationTimestamp) > 1*time.Second {
		t.Errorf("Expected expiration in 0 seconds, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}

	fmt.Println("\nServer: We generate a a bunch of 1 day tokens but only enter the latest one")
	for xn := 0; xn < 5; xn++ {
		_, err = serverSimulator.GenerateTokenFromDate(time.Now().Add(24*time.Hour), false)
		if err != nil {
			t.Error(err)
		}
	}
	thisToken, err = serverSimulator.GenerateTokenFromDate(time.Now().Add(24*time.Hour), false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token:", thisToken)
	fmt.Println("Device: We enter the generated token into the device")
	err = deviceSimulator.EnterToken(thisToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 1 day and the count synchronised with the server)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < (24*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (24*time.Hour+1*time.Second) {
		t.Errorf("Expected time until expiration is 24 hours, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}

	fmt.Println("\nWe add generate 9 tokens each add-time of 1 day")
	tokens := make([]string, 0)
	for xn := 0; xn < 9; xn++ {
		thisToken, err = serverSimulator.GenerateTokenFromValue(1, openpaygotoken.AddTime)
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, thisToken)
	}
	fmt.Print("Tokens: ")
	for _, token := range tokens {
		fmt.Print(token, " ")
	}
	fmt.Println("\nDevice: We enter the 9th token into the device")
	err = deviceSimulator.EnterToken(tokens[8])
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with +1 day (2 days total))")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < (48*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (48*time.Hour+1*time.Second) {
		t.Errorf("Expected time until expiration is 48 hours, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}
	fmt.Println("Device: We enter the 1st token into the device")
	err = deviceSimulator.EnterToken(tokens[0])
	if err != nil && !errors.Is(err, &simulators.ErrOldToken{}) {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status , it should not have changed, because its more than 5 add times before")
	deviceSimulator.PrintStatus()
	if time.Until(deviceSimulator.ExpirationTimestamp) < (48*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (48*time.Hour+1*time.Second) {
		t.Errorf("Expected time until expiration is 48 hours, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}
	fmt.Println("Device: We enter the tokens 5, 4, 3 and 2 into the device")
	for xn := 5; xn > 1; xn-- {
		err = deviceSimulator.EnterToken(tokens[xn])
		if err != nil && !errors.Is(err, &simulators.ErrOldToken{}) {
			t.Error(err)
		}
	}
	fmt.Println("Device: We check the device status (should be active with +4 day (6 days total))")
	deviceSimulator.PrintStatus()
	if time.Until(deviceSimulator.ExpirationTimestamp) < (6*24*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (6*24*time.Hour+1*time.Second) {
		t.Errorf("Expected time until expiration is 6 days, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}

	fmt.Println("\nServer: We add generate 2 tokens, first add-time and then set-time")
	token1, err := serverSimulator.GenerateTokenFromValue(1, openpaygotoken.AddTime)
	if err != nil {
		t.Error(err)
	}
	token2, err := serverSimulator.GenerateTokenFromValue(0, openpaygotoken.SetTime)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Tokens:", token1, token2)
	fmt.Println("Device: We enter the 2nd token")
	err = deviceSimulator.EnterToken(token2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 0 days)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < -1*time.Second || time.Until(deviceSimulator.ExpirationTimestamp) > 1*time.Second {
		t.Errorf("Expected time until expiration is 0 seconds, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}
	fmt.Println("Device: We enter the 1st token into the device")
	err = deviceSimulator.EnterToken(token1)
	if err != nil && !errors.Is(err, &simulators.ErrOldToken{}) {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status, it should not have changed, because you cannot use an add-time token older than a set-time")
	deviceSimulator.PrintStatus()
	if time.Until(deviceSimulator.ExpirationTimestamp) < -1*time.Second || time.Until(deviceSimulator.ExpirationTimestamp) > 1*time.Second {
		t.Errorf("Expected time until expiration is 0 seconds, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}

	fmt.Println("\nServer: We add generate 2 tokens, first set-time and then add-time")
	token1, err = serverSimulator.GenerateTokenFromValue(1, openpaygotoken.SetTime)
	if err != nil {
		t.Error(err)
	}
	token2, err = serverSimulator.GenerateTokenFromValue(2, openpaygotoken.AddTime)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Tokens:", token1, token2)
	fmt.Println("Device: We enter the 2nd token")
	err = deviceSimulator.EnterToken(token2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status (should be active with 2 day)")
	deviceSimulator.PrintStatus()
	if deviceSimulator.Count != serverSimulator.Count {
		t.Errorf("Expected count to be %d, got %d", serverSimulator.Count, deviceSimulator.Count)
	}
	if time.Until(deviceSimulator.ExpirationTimestamp) < (48*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (48*time.Hour+1*time.Second) {
		t.Errorf("Expected time until expiration is 48 hours, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}
	fmt.Println("Device: We enter the 1st token into the device")
	err = deviceSimulator.EnterToken(token1)
	if err != nil && !errors.Is(err, &simulators.ErrOldToken{}) {
		t.Error(err)
	}
	fmt.Println("Device: We check the device status, it should not have changed, because you cannot use an older set-time token")
	deviceSimulator.PrintStatus()
	if time.Until(deviceSimulator.ExpirationTimestamp) < (48*time.Hour-1*time.Second) || time.Until(deviceSimulator.ExpirationTimestamp) > (48*time.Hour+1*time.Second) {
		t.Errorf("Expected time until expiration is 48 hours, got %s", time.Until(deviceSimulator.ExpirationTimestamp))
	}

}
