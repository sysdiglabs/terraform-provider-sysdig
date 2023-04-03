package monitor

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigMonitorClient) GetNotificationChannelById(ctx context.Context, id int) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)

	if nc.Version == 0 {
		err = fmt.Errorf("NotificationChannel with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigMonitorClient) GetNotificationChannelUrl(id int) string {
	return fmt.Sprintf("%s/api/notificationChannels/%d", client.URL, id)
}
