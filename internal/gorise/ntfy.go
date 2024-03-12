package gorise

type Ntfy struct {
	authorization string
	url           string
	topics        []string
}

// func (n Ntfy) Send(ctx context.Context, msg Message) error {
// 	var errs []error
//
// 	if len(msg.Attachments) == 0 {
// 		for _, topic := range n.topics {
//
// 		}
// 	} else {
// 		for _, topic := range n.topics {
// 			req, _ := http.NewRequestWithContext(ctx, "PUT", "https://ntfy.sh/flowers", file)
// 			req.Header.Set("Filename", "flower.jpg")
// 			http.DefaultClient.Do(req)
// 		}
// 	}
//
// 	return errors.Join(errs...)
// }
