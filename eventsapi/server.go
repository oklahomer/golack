package eventsapi

import (
	"github.com/oklahomer/golack/event"
	"log"
	"net/http"
)

// EventReceiver defines an interface to subscribe to incoming events.
type EventReceiver interface {
	Receive(wrapper *EventWrapper)
}

// WithRequestValidator returns a function to set given rv on SetupHandler.
func WithRequestValidator(rv RequestValidator) func(*option) {
	return func(o *option) {
		o.RequestValidator = rv
	}
}

type option struct {
	RequestValidator RequestValidator
}

// SetupHandler construct http.HandlerFunc to serve Events API endpoint and receive incoming events.
func SetupHandler(receiver EventReceiver, opts ...func(*option)) http.HandlerFunc {
	opt := &option{}
	for _, o := range opts {
		o(opt)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		// Read the incoming request
		req, err := NewSlackRequest(request)
		if err != nil {
			switch err.(type) {
			case *BadRequestError, *event.MalformedPayloadError:
				writer.WriteHeader(http.StatusBadRequest)
				return

			default:
				writer.WriteHeader(http.StatusInternalServerError)
				return

			}
		}

		// Validate the request
		if opt.RequestValidator != nil && !opt.RequestValidator.Validate(req) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Decode payload
		ev, err := DecodePayload(req)
		if err != nil {
			switch err.(type) {
			case *event.MalformedPayloadError:
				writer.WriteHeader(http.StatusBadRequest)
				return

			default:
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Return HTTP response and dispatch task
		switch typed := ev.(type) {
		case *URLVerification:
			writer.Header().Set("Content-Type", "text/plain")
			writer.Write([]byte(typed.Challenge))
			return

		case *EventWrapper:
			receiver.Receive(typed)
			writer.WriteHeader(http.StatusOK)
			return

		default:
			writer.WriteHeader(http.StatusOK)
			log.Printf("Successfully decoded the payload but do not know how to handle %T", typed)
			return

		}
	}
}
