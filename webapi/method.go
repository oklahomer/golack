package webapi

import "golang.org/x/xerrors"

var (
	// JSONAcceptableMethods lists all Web API methods that support JSON serialized payload.
	// See https://api.slack.com/web#methods_supporting_json
	JSONAcceptableMethods = []string{
		"admin.apps.approve",
		"admin.apps.restrict",
		"admin.conversations.setTeams",
		"admin.inviteRequests.approve",
		"admin.inviteRequests.approved.list",
		"admin.inviteRequests.denied.list",
		"admin.inviteRequests.deny",
		"admin.inviteRequests.list",
		"admin.teams.create",
		"admin.teams.list",
		"admin.teams.settings.info",
		"admin.teams.settings.setDescription",
		"admin.teams.settings.setDiscoverability",
		"admin.teams.settings.setName",
		"admin.usergroups.addChannels",
		"admin.usergroups.addTeams",
		"admin.usergroups.listChannels",
		"admin.usergroups.removeChannels",
		"admin.users.assign",
		"admin.users.invite",
		"admin.users.list",
		"admin.users.remove",
		"admin.users.session.reset",
		"admin.users.setAdmin",
		"admin.users.setExpiration",
		"admin.users.setOwner",
		"admin.users.setRegular",
		"api.test",
		"auth.test",
		"calls.add",
		"calls.end",
		"calls.info",
		"calls.participants.add",
		"calls.participants.remove",
		"calls.update",
		"channels.archive",
		"channels.create",
		"channels.invite",
		"channels.join",
		"channels.kick",
		"channels.leave",
		"channels.mark",
		"channels.rename",
		"channels.setPurpose",
		"channels.setTopic",
		"channels.unarchive",
		"chat.delete",
		"chat.deleteScheduledMessage",
		"chat.meMessage",
		"chat.postEphemeral",
		"chat.postMessage",
		"chat.scheduleMessage",
		"chat.scheduledMessages.list",
		"chat.unfurl",
		"chat.update",
		"conversations.archive",
		"conversations.close",
		"conversations.create",
		"conversations.invite",
		"conversations.join",
		"conversations.kick",
		"conversations.leave",
		"conversations.mark",
		"conversations.open",
		"conversations.rename",
		"conversations.setPurpose",
		"conversations.setTopic",
		"conversations.unarchive",
		"dialog.open",
		"dnd.endDnd",
		"dnd.endSnooze",
		"files.comments.delete",
		"files.delete",
		"files.revokePublicURL",
		"files.sharedPublicURL",
		"groups.archive",
		"groups.create",
		"groups.invite",
		"groups.kick",
		"groups.leave",
		"groups.mark",
		"groups.open",
		"groups.rename",
		"groups.setPurpose",
		"groups.setTopic",
		"groups.unarchive",
		"im.close",
		"im.mark",
		"im.open",
		"mpim.close",
		"mpim.mark",
		"mpim.open",
		"pins.add",
		"pins.remove",
		"reactions.add",
		"reactions.remove",
		"reminders.add",
		"reminders.complete",
		"reminders.delete",
		"stars.add",
		"stars.remove",
		"usergroups.create",
		"usergroups.disable",
		"usergroups.enable",
		"usergroups.update",
		"usergroups.users.update",
		"users.profile.set",
		"users.setActive",
		"users.setPresence",
		"views.open",
		"views.publish",
		"views.push",
		"views.update",
	}
	jsonAcceptableMethodMap = func() map[string]struct{} {
		m := map[string]struct{}{}
		for _, method := range JSONAcceptableMethods {
			m[method] = struct{}{}
		}
		return m
	}()
	ErrJSONPayloadNotSupported = xerrors.New("JSON payload is not supported")
)

func IsJSONPayloadSupportedMethod(slackMethod string) bool {
	_, ok := jsonAcceptableMethodMap[slackMethod]
	return ok
}