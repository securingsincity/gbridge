package gbridge

import (
	"log"
	"net/http"
	"encoding/json"
	"os"
	"io"
)

func (b *Bridge) HandleSmartHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

	req := IntentMessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding request:", err)
	}
	for _, i := range req.Inputs {
		log.Printf("Intent: %s-> %s\n", i.Intent, string(i.Payload))
		switch i.Intent {
		case "action.devices.QUERY":
			requestBody := QueryRequest{}
			if err := json.Unmarshal(i.Payload, &requestBody); err != nil {
				log.Println(err)
				return
			}
			log.Printf("QUERY: %+v\n", requestBody)

			responseBody := QueryResponse{
				Devices: make(map[string]DeviceState),
			}

			for _, d := range requestBody.Devices {
				if ctx, ok := b.Devices[d.ID]; ok {
					if ctx.Query != nil {
						res := DeviceState{}
						ctx.Query(ctx.Device, &res)
						responseBody.Devices[d.ID] = res
					}
				}
			}

			if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(IntentMessageResponse{
				RequestId: req.RequestId,
				Payload:   responseBody,
			}); err != nil {
				log.Println(err)
			}
			return
		case "action.devices.SYNC":
			log.Println("SYNC")
			devices := []Device{}
			for _, ctx := range b.Devices {
				devices = append(devices, ctx.Device)
			}
			if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(IntentMessageResponse{
				RequestId: req.RequestId,
				Payload: SyncResponse{
					AgentUserId: b.AgentUserId,
					Devices:     devices,
				},
			}); err != nil {
				log.Println("result error:", err)
			}
			return
		case "action.devices.EXECUTE":
			requestBody := ExecRequest{}
			if err := json.Unmarshal(i.Payload, &requestBody); err != nil {
				log.Println(err)
				return
			}
			log.Printf("EXEC: %+v\n", requestBody)
			ids := []string{}
			for _, c := range requestBody.Commands {
				for _, d := range c.Devices {
					ids = append(ids, d.ID)
				}
			}
			responseBody := ExecResponse{}
			for _, c := range requestBody.Commands {
				for _, d := range c.Devices {
					if ctx, ok := b.Devices[d.ID]; ok {
						if ctx.Exec != nil {
							for _, e := range c.Execution {
								r := CommandResponse{
									Ids:    ids,
									Status: CommandStatusError,
								}
								ctx.Exec(b.Devices[d.ID].Device, e, &r)
								if r.Status == CommandStatusError && r.ErrorCode == "" {
									r.ErrorCode = DeviceErrorUnknownError
								}
								responseBody.Commands = append(responseBody.Commands, r)
							}
						}
					} else {
						responseBody.Commands = append(responseBody.Commands, CommandResponse{
							Ids:       ids,
							Status:    CommandStatusError,
							ErrorCode: DeviceErrorDeviceNotFound,
						})
					}
				}
			}

			if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(IntentMessageResponse{
				RequestId: req.RequestId,
				Payload:   responseBody,
			}); err != nil {
				log.Println(err)
			}
			return
		}
	}
}
