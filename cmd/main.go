package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shivanshkc/rosenbridge/pkg/bridge"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Listening on 8080")
	if err := http.ListenAndServe(":8080", &handler{}); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

type handler struct{}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	wsBridge := bridge.NewWebsocketBridge()
	if err := wsBridge.Connect(request, writer); err != nil {
		errHTTP := errutils.ToHTTPError(err)
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return
	}

	for {
		ctx, cancelFunc := context.WithTimeout(request.Context(), time.Minute)

		reqMessage := fmt.Sprintf(`{"id": "%s", "data": "something useful"}`, uuid.NewString())
		idFunc := func(msg []byte) any {
			decoded := map[string]any{}
			if err := json.Unmarshal(msg, &decoded); err != nil {
				return nil
			}
			return decoded["id"]
		}

		response, err := wsBridge.SendSync(ctx, []byte(reqMessage), idFunc)
		if err != nil {
			fmt.Println("error in SendSync call:", err)
			cancelFunc()
			continue
		}

		fmt.Println("Client response:", string(response))
		cancelFunc()
	}
}
