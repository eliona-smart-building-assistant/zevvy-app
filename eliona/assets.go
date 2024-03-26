//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package eliona

import (
	"fmt"
	"template/apiserver"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func CreateAssets(config apiserver.Configuration, root asset.Root) error {
	for _, projectId := range *config.ProjectIDs {
		assetsCreated, err := asset.CreateAssets(root, projectId)
		if err != nil {
			return err
		}
		if assetsCreated != 0 {
			if err := notifyUser(*config.UserId, projectId, assetsCreated); err != nil {
				return fmt.Errorf("notifying user about CAC: %v", err)
			}
		}
	}
	return nil
}

func notifyUser(userId string, projectId string, assetsCreated int) error {
	receipt, _, err := client.NewClient().CommunicationAPI.
		PostNotification(client.AuthenticationContext()).
		Notification(
			api.Notification{
				User:      userId,
				ProjectId: *api.NewNullableString(&projectId),
				Message: *api.NewNullableTranslation(&api.Translation{
					De: api.PtrString(fmt.Sprintf("Template App hat %d neue Assets angelegt. Diese sind nun im Asset-Management verfügbar.", assetsCreated)),
					En: api.PtrString(fmt.Sprintf("Template app added %v new assets. They are now available in Asset Management.", assetsCreated)),
				}),
			}).
		Execute()
	log.Debug("eliona", "posted notification about CAC: %v", receipt)
	if err != nil {
		return fmt.Errorf("posting CAC notification: %v", err)
	}
	return nil
}
