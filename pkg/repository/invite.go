package repository

import (
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/utils"
	"math"
	"sync"
)

type invite struct {
	mu *sync.Mutex
}

var Invite = &invite{mu: &sync.Mutex{}}

func (i *invite) CreateCode() string {
	rand := utils.NewRandom()

	var (
		key    = rand.String(10)
		trials = 0.0
	)

	i.mu.Lock()
	if IsExists[models.Invite]([]string{"invites.key"}, []any{key}) || IsExists[models.Group]([]string{"custom_invite"}, []any{key}) {
		trials += 0.3
		times := int(math.Trunc(trials))
		if times == 1 {
			times = 2
		}

		key = rand.String(10 * times)
	}
	i.mu.Unlock()

	return key
}

func (i *invite) FindInviteByCode(id string) *models.Invite {
	return FindModelBy[models.Invite]([]string{"key"}, []any{id})

}

func (i *invite) FindGroupInviteByCode(id string) *models.Group {
	return FindModelBy[models.Group]([]string{"custom_invite"}, []any{id})
}
