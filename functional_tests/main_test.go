package functional_tests

import (
	"context"
	"log"
	"testing"
	"time"
)

// Run server only once for all functional tests
func TestMain(m *testing.M) {
	srv := setupServer()
	time.Sleep(time.Second * 1)
	defer func() {
		err := srv.Shutdown(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}()
	m.Run()
}
