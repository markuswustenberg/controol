package osc_test

import (
	"controol/osc"
	"testing"
	"time"

	osc2 "github.com/hypebeast/go-osc/osc"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	t.Run("Send errors if there are less than four args", func(t *testing.T) {
		t.Parallel()
		err := osc.Send(nil)
		require.Error(t, err)

		err = osc.Send([]string{"1", "2", "3"})
		require.Error(t, err)
	})

	t.Run("Send errors if the port is not an int", func(t *testing.T) {
		t.Parallel()
		err := osc.Send([]string{"host", "pony", "/address", "1"})
		require.Error(t, err)
	})

	t.Run("Send sends messages with int/float/bool/string args", func(t *testing.T) {
		msgs := make(chan *osc2.Message, 1)

		server := &osc2.Server{Addr: "localhost:9000"}

		server.Handle("*", func(msg *osc2.Message) {
			msgs <- msg
		})

		go func() {
			server.ListenAndServe()
		}()

		err := osc.Send([]string{"localhost", "9000", "/address", "1", "1.0", "true", "pony"})
		require.NoError(t, err)

		select {
		case msg := <-msgs:
			require.Equal(t, "/address", msg.Address)
			args := msg.Arguments

			a1, ok := args[0].(int32)
			require.True(t, ok)
			require.Equal(t, int32(1), a1)

			a2, ok := args[1].(float32)
			require.True(t, ok)
			require.Equal(t, float32(1), a2)

			a3, ok := args[2].(bool)
			require.True(t, ok)
			require.Equal(t, true, a3)

			a4, ok := args[3].(string)
			require.True(t, ok)
			require.Equal(t, "pony", a4)
		case <-time.After(100 * time.Millisecond):
			require.FailNow(t, "did not receive message")
		}
	})
}

func TestReceive(t *testing.T) {
	t.Run("Receive errors if there are less than two args", func(t *testing.T) {
		t.Parallel()
		err := osc.Receive(nil)
		require.Error(t, err)

		err = osc.Receive([]string{"1"})
		require.Error(t, err)
	})

	t.Run("Receive errors on bad host and/or port", func(t *testing.T) {
		err := osc.Receive([]string{"flerp", "1234"})
		require.Error(t, err)

		err = osc.Receive([]string{"localhost", "pony"})
		require.Error(t, err)
	})

	t.Run("Receive receives messages", func(t *testing.T) {
		msgs := make(chan *osc2.Message, 1)

		go func() {
			osc.Receive([]string{"localhost", "9000"}, func(msg *osc2.Message) {
				msgs <- msg
			})
		}()

		client := osc2.NewClient("localhost", 9000)
		go func() {
			// Keep sending messages
			for {
				msg := osc2.NewMessage("/pony")
				err := client.Send(msg)
				require.NoError(t, err)
				time.Sleep(10 * time.Millisecond)
			}
		}()

		select {
		case msg := <-msgs:
			require.Equal(t, "/pony", msg.Address)
		case <-time.After(100 * time.Millisecond):
			require.FailNow(t, "did not receive message")
		}
	})
}
