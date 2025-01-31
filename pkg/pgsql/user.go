package pgsql

import (
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) SelectUserByID(id int) (*types.User, error) {
	sess := c.GetSession()

	res := &types.User{}

	q := `select * from users where id = ?`
	if err := sess.SelectBySql(q, id).LoadOne(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) CountUsersByTelegramID(tid int64) (int, error) {
	sess := c.GetSession()

	var count int

	q := `select count(*) from users where telegram_id = ?`
	if _, err := sess.SelectBySql(q, tid).Load(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Client) SelectUsersByTelegramID(tid int64) ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)

	q := `select * from users where telegram_id = ? order by id desc`
	if _, err := sess.SelectBySql(q, tid).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectBotUsersByTelegramID(bid int, tid int64) ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)

	q := `select * from users where bot_id = ? and telegram_id = ? order by id desc limit 1`
	if _, err := sess.SelectBySql(q, bid, tid).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectUsersWithoutDeeplinkID(limit int) ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)

	q := `select * from users where deeplink_id = 0 order by id desc limit ?`
	if _, err := sess.SelectBySql(q, limit).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectAllUsersWithEmptyCreationEvent() ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)

	q := `select * from users where event_created = 'message' order by id desc`
	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectUsersByForwardSenderName(fsn string) ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)

	q := `select * from users where forward_sender_name = ?`
	if _, err := sess.SelectBySql(q, fsn).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectUsersByUsername(username string) ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)
	q := `select * from users where username = ?`
	if _, err := sess.SelectBySql(q, username).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectRandomReadyUsersByDepotChannelHash(botID int, hash string, limit int) ([]*types.User, error) {
	sess := c.GetSession()

	res := make([]*types.User, 0)
	q := `select * from users where mailing_state = ? and bot_id = ? and depot_channel_hash = ? order by random() limit ?`
	if _, err := sess.SelectBySql(q, types.UserMailingStateReady, botID, hash, limit).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) UpdateUsersMailingState(state string, ids []int) error {
	sess := c.GetSession()

	q := `update users set mailing_state = ?, mailing_state_updated_at = now() where id = any (?)`
	if _, err := sess.UpdateBySql(q, state, pq.Array(ids)).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SetBotUsersMailingStatesReady(botID int, hash string) error {
	sess := c.GetSession()

	states := []string{types.UserMailingStateInProgress, types.UserMailingStateFinished}

	q := `update users
	set mailing_state            = ?,
		mailing_state_updated_at = now()
	where bot_id = ?
	  and depot_channel_hash = ?
	  and mailing_state = any (?)`
	if _, err := sess.UpdateBySql(q,
		types.UserMailingStateReady,
		botID,
		hash,
		pq.Array(states),
	).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotUsersSetDefaultChannel(botID int, hash string, tgchid int64, tgchurl string) error {
	sess := c.GetSession()

	q := `update users
	set telegram_channel_id  = ?,
		telegram_channel_url = ?,
		depot_channel_hash   = ?
	where bot_id = ?
	  and depot_channel_hash = ''
	  and telegram_channel_id = 0
	  and telegram_channel_url = ''`
	if _, err := sess.UpdateBySql(q, tgchid, tgchurl, hash, botID).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotUsersTelegramChannel(botID int, hash string, tgchid int64, tgchurl string) error {
	sess := c.GetSession()

	q := `update users
	set telegram_channel_id  = ?,
		telegram_channel_url = ?
	where bot_id = ?
	  and depot_channel_hash = ?`
	if _, err := sess.UpdateBySql(q, tgchid, tgchurl, botID, hash).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotUserTelegramChannel(userID int, hash string, tgchid int64, tgchurl string) error {
	sess := c.GetSession()

	q := `update users
	set telegram_channel_id  = ?,
		telegram_channel_url = ?,
		depot_channel_hash   = ?
	where id = ?`
	if _, err := sess.UpdateBySql(q, tgchid, tgchurl, hash, userID).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotUserMailingState(botToken string, telegramID int64, state string) error {
	sess := c.GetSession()

	if state == types.UserMailingStateBlocked {
		qtpl := `update users
set mailing_state            = case when mailing_failed_attempts + 1 >= %d then '%s' else '%s' end,
    mailing_failed_attempts  = mailing_failed_attempts + 1,
    mailing_state_updated_at = now()
where bot_id = (select id from bots where bot_token = ?)
  and telegram_id = ?;`

		q := fmt.Sprintf(qtpl, types.UserMailingMaxFailedAttempts, types.UserMailingStateBlocked, types.UserMailingStateFinished)
		if _, err := sess.UpdateBySql(q, botToken, telegramID).Exec(); err != nil {
			return err
		}

	} else {
		q := `update users
set mailing_state            = ?,
    mailing_failed_attempts  = 0,
    mailing_state_updated_at = now()
where bot_id = (select id from bots where bot_token = ?)
  and telegram_id = ?;`
		if _, err := sess.UpdateBySql(q, state, botToken, telegramID).Exec(); err != nil {
			return err
		}

	}

	return nil
}

func (c *Client) CreateUser(user *types.User) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("users").
		Columns(
			"bot_id",
			"deeplink_id",
			"telegram_id",
			"depot_channel_hash",
			"telegram_channel_id",
			"telegram_channel_url",
			"first_name",
			"last_name",
			"username",
			"seen",
			"forward_sender_name",
			"is_bot",
			"is_premium",
			"language_code",
			"event_created",
		).
		Record(user).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUserMessagedAt(id int) error {
	sess := c.GetSession()

	q := `update users set messaged_at = now() where id = ?`
	if _, err := sess.UpdateBySql(q, id).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUserOnMessage(old, new *types.User) error {
	sess := c.GetSession()

	updated := false

	stmt := sess.Update("users").
		Set("messaged_at", "now()")

	if old.Firstname != new.Firstname && new.Firstname != "" {
		stmt.Set("first_name", new.Firstname)
		updated = true
	}

	if old.Lastname != new.Lastname && new.Lastname != "" {
		stmt.Set("last_name", new.Lastname)
		updated = true
	}

	if old.Username != new.Username && new.Username != "" {
		stmt.Set("username", new.Username)
		updated = true
	}

	if old.IsPremium != new.IsPremium {
		stmt.Set("is_premium", new.IsPremium)
		updated = true
	}

	if old.LanguageCode != new.LanguageCode && new.LanguageCode != "" {
		stmt.Set("language_code", new.LanguageCode)
		updated = true
	}

	if old.ForwardSenderName != new.ForwardSenderName && new.ForwardSenderName != "" {
		stmt.Set("forward_sender_name", new.ForwardSenderName)
		updated = true
	}

	if old.MailingState == types.UserMailingStateBlocked {
		stmt.Set("mailing_state", types.UserMailingStateReady)
		stmt.Set("mailing_failed_attempts", 0)
		stmt.Set("mailing_state_updated_at", "now()")
		updated = true
	}

	if old.DepotChannelHash != new.DepotChannelHash && new.DepotChannelHash != "" {
		stmt.Set("depot_channel_hash", new.DepotChannelHash)
		stmt.Set("telegram_channel_id", new.TelegramChannelID)
		stmt.Set("telegram_channel_url", new.TelegramChannelURL)
		updated = true
	}

	if updated {
		stmt.Set("updated_at", "now()")
	}

	stmt = stmt.Where("id = ?", old.ID)
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUserDeeplink(uid, did int) error {
	sess := c.GetSession()

	q := `update users set deeplink_id = ? where id = ?`
	if _, err := sess.UpdateBySql(q, did, uid).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUserDeeplinkAndCreatedAt(uid, did int, createdAt time.Time) error {
	sess := c.GetSession()

	q := `update users set deeplink_id = ?, created_at = ?, event_created = 'register' where id = ?`
	if _, err := sess.UpdateBySql(q, did, createdAt, uid).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUserDepositState(id int) error {
	sess := c.GetSession()

	q := `update users set deposited = true, deposited_at = now() where id = ? and deposited = false`
	if _, err := sess.UpdateBySql(q, id).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUserHeadersInfo(id int, ip, userAgent, countryCode, osName, deviceType string) error {
	sess := c.GetSession()

	stmt := sess.Update("users")

	if ip != "" {
		stmt.Set("ip", ip)
	}

	if userAgent != "" {
		stmt.Set("user_agent", userAgent)
	}

	if countryCode != "" {
		stmt.Set("country_code", countryCode)
	}

	if osName != "" {
		stmt.Set("os_name", osName)
	}

	if deviceType != "" {
		stmt.Set("device_type", deviceType)
	}

	stmt.Set("updated_at", "now()").Where("id = ?", id)
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}
