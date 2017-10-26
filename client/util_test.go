package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/api/types/swarm"

	"github.com/franela/goblin"
)

func Test_utils(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Tests for client/util.go", func() {
		responseFile, err := os.Open("../fixtures/tests/service.json")
		if err != nil {
			t.Fatalf("Unable to read test data: %s", err.Error())
		}
		defer func() { _ = responseFile.Close() }()
		responseBody, _ := ioutil.ReadAll(responseFile)
		var service []swarm.Service
		err = json.Unmarshal(responseBody, &service)
		g.Describe("The CheckName function", func() {
			if err != nil {
				t.Fatal("Unable to unmarshal data", err.Error())
			}
			g.It("should return true if name isn't in use", func() {

				err = CheckName("foo", service)
				g.Assert(err == nil).IsTrue()
			})
			g.It("should return false if name is in use", func() {

				err = CheckName("angry_goldstine", service)
				g.Assert(err == nil).IsFalse()
			})
		})
		g.Describe("The CheckPort function", func() {
			if err != nil {
				t.Fatal("Unable to unmarshal data", err.Error())
			}
			g.It("should return true if port isn't in use", func() {

				err = CheckPort(50763, service)
				g.Assert(err == nil).IsTrue()
			})
			g.It("should return false if port is in use", func() {

				err = CheckPort(50789, service)
				g.Assert(err == nil).IsFalse()
			})
		})
	})
}
