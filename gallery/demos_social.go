package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/social"
	"github.com/egoist/mygo/ui"
)

var socialEntries = []Entry{
	{ID: 1191, Name: "Post", Category: "social", Summary: "Author, text, media, actions.", Demo: demoPost},
	{ID: 1192, Name: "Feed", Category: "social", Summary: "Posts with new pill; cue Notify.", Demo: demoSocialFeed},
	{ID: 1193, Name: "PostComposer", Category: "social", Summary: "Text with count and media; cue Send.", Demo: demoPostComposer},
	{ID: 1194, Name: "PostActions", Category: "social", Summary: "Like, comment, share, save; cue Pop.", Demo: demoPostActions},
	{ID: 1195, Name: "CommentThread", Category: "social", Summary: "Nested comments with reply; cue Send.", Demo: demoCommentThread},
	{ID: 1196, Name: "ProfileHeader", Category: "social", Summary: "Cover, avatar, name, bio, follow; cue Toggle.", Demo: demoProfileHeader},
	{ID: 1197, Name: "ProfileStats", Category: "social", Summary: "Posts, followers, following; roll.", Demo: demoProfileStats},
	{ID: 1198, Name: "WhoToFollow", Category: "social", Summary: "Suggestions with follow; cue Toggle.", Demo: demoWhoToFollow},
	{ID: 1199, Name: "StoryRing", Category: "social", Summary: "Avatars with rings for unseen stories; cue Open.", Demo: demoStoryRing},
	{ID: 1200, Name: "Hashtag", Category: "social", Summary: "Tag chip with count; cue Tap.", Demo: demoHashtag},
	{ID: 1201, Name: "Trending", Category: "social", Summary: "Trending topics ranked.", Demo: demoTrending},
	{ID: 1202, Name: "PollPost", Category: "social", Summary: "Poll with votes and time left; cue Select.", Demo: demoPollPost},
	{ID: 1203, Name: "MediaPost", Category: "social", Summary: "Image grid of up to four.", Demo: demoMediaPost},
	{ID: 1204, Name: "Mentions", Category: "social", Summary: "\"@\" suggestions while typing; cue Select.", Demo: demoMentions},
	{ID: 1205, Name: "UserCard", Category: "social", Summary: "Hover card of a user; cue Toggle.", Demo: demoUserCard},
	{ID: 1206, Name: "Followers", Category: "social", Summary: "List with follow back; cue Toggle.", Demo: demoFollowers},
	{ID: 1207, Name: "ProfileBadges", Category: "social", Summary: "Badges earned.", Demo: demoProfileBadges},
	{ID: 1208, Name: "ReactionBar", Category: "social", Summary: "Several reactions with counts; cue Pop.", Demo: demoReactionBar},
}

func demoSocialUser() social.User {
	return social.User{ID: "ana", Name: "Ana", Handle: "@ana", Bio: "Designer", Avatar: photos()[0], Cover: photos()[1]}
}

func demoSocialPost() social.PostData {
	return social.PostData{ID: "p1", Text: "Shipping today", Author: demoSocialUser(), Media: photos()[2], Likes: 3, Comments: 1, Shares: 2}
}

// studio is the ceramics studio cast the 1191–1208 demos share.
func studio() []social.User {
	return []social.User{
		{ID: "maya", Name: "Maya Chen", Handle: "@mayachen", Bio: "Product designer at Kiln & Co. Lamps, mugs and slow mornings.", Cover: photos()[0]},
		{ID: "daniel", Name: "Daniel Okafor", Handle: "@dokafor", Bio: "Still-life photographer in Portland."},
		{ID: "priya", Name: "Priya Nair", Handle: "@priyanair", Bio: "Ceramicist. Glazes, kilns, test tiles."},
		{ID: "lucas", Name: "Lucas Martin", Handle: "@lucasm", Bio: "Industrial designer working in brass and oak."},
		{ID: "sofia", Name: "Sofia Rossi", Handle: "@sofiarossi", Bio: "Art director. Color, light and print.", Cover: photos()[1]},
		{ID: "jonas", Name: "Jonas Berg", Handle: "@jonasberg", Bio: "Furniture maker from Malmö."},
	}
}

func studioPosts() []social.PostData {
	p, u := productShots(), studio()
	return []social.PostData{
		{ID: "p1", Author: u[0], Time: "2h", Media: p[0], Likes: 214, Comments: 18, Shares: 9,
			Text: "The brass desk lamp is finally out of the studio. Six months of prototypes, one very patient photographer."},
		{ID: "p2", Author: u[2], Time: "5h", Media: p[1], Likes: 96, Comments: 7, Shares: 3,
			Text: "New glaze test: speckled oat over terracotta. Two more firings before it goes on the mugs."},
		{ID: "p3", Author: u[3], Time: "1d", Likes: 41, Comments: 5, Shares: 1,
			Text: "Reminder: the oak tray drop opens Thursday at 10:00. Only 40 pieces this run."},
	}
}

func demoPost(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	state := ui.Local(root, "post", func() social.PostState { return social.PostState{} })
	root.Children(func() {
		if social.Post(c, studioPosts()[0], state, social.PostOptions{}).Changed() {
			log(fmt.Sprintf("liked %v saved %v action %s", state.Liked, state.Saved, state.Action))
		}
	})
}

func demoSocialFeed(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	fresh := boolState(root, "fresh", true)
	root.Children(func() {
		if social.Feed(c, studioPosts(), fresh, social.FeedOptions{}).Changed() {
			log("fresh " + fmt.Sprint(*fresh))
		}
	})
}

func demoPostComposer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	draft := textState(root, "draft", "Kiln day. Twelve mugs in, fingers crossed for the oat glaze.")
	added := boolState(root, "added", false)
	root.Children(func() {
		form, send := social.PostComposer(c, draft, productShots()[1], added, social.PostComposerOptions{})
		if send.Changed() {
			log("sent")
		} else if form.Changed() {
			log(fmt.Sprintf("runes %d media %v", len([]rune(*draft)), *added))
		}
	})
}

func demoPostActions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	state := ui.Local(root, "actions", func() social.PostState { return social.PostState{Saved: true} })
	root.Children(func() {
		if social.PostActions(c, state, social.PostActionsOptions{}).Changed() {
			log(fmt.Sprintf("liked %v saved %v action %s", state.Liked, state.Saved, state.Action))
		}
	})
}

func demoCommentThread(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	draft := textState(root, "draft", "Thursday works. I'll bring the test tiles.")
	u := studio()
	comments := []social.Comment{
		{ID: "c1", Author: u[1], Time: "1h", Text: "That shade reads so well on camera. Can we reshoot it in morning light?",
			Replies: []social.Comment{{ID: "c2", Author: u[0], Time: "48m", Text: "Yes, Thursday before the drop?"}}},
		{ID: "c3", Author: u[4], Time: "20m", Text: "Love the arm. The cable could be cloth-wrapped in a darker tone."},
	}
	root.Children(func() {
		form, send := social.CommentThread(c, comments, draft, social.CommentThreadOptions{})
		if send.Changed() {
			log("replied")
		} else if form.Changed() {
			log("runes " + fmt.Sprint(len([]rune(*draft))))
		}
	})
}

func demoProfileHeader(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	following := boolState(root, "following", false)
	root.Children(func() {
		if social.ProfileHeader(c, studio()[0], following, social.ProfileHeaderOptions{}).Changed() {
			log("following " + fmt.Sprint(*following))
		}
	})
}

func demoProfileStats(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() { social.ProfileStats(c, 128, 12480, 342, social.ProfileStatsOptions{}) })
}

// followed counts the true entries in a follow map.
func followed(m map[string]bool) int {
	n := 0
	for _, on := range m {
		if on {
			n++
		}
	}
	return n
}

func demoWhoToFollow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	following := ui.Local(root, "following", func() map[string]bool { return map[string]bool{"lucas": true} })
	root.Children(func() {
		if social.WhoToFollow(c, studio()[2:5], following, social.WhoToFollowOptions{}).Changed() {
			log("following " + fmt.Sprint(followed(*following)))
		}
	})
}

func demoStoryRing(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	opened := textState(root, "opened", "")
	u := studio()
	stories := []social.Story{{ID: "s1", User: u[0]}, {ID: "s2", User: u[1]}, {ID: "s3", User: u[2]}, {ID: "s4", User: u[3], Seen: true},
		{ID: "s5", User: u[4], Seen: true}, {ID: "s6", User: u[5], Seen: true}}
	root.Children(func() {
		if social.StoryRing(c, stories, opened).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoHashtag(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		if social.Hashtag(c, "#ceramics", 12480, opened, social.HashtagOptions{}).Changed() {
			log("opened " + fmt.Sprint(*opened))
		}
	})
}

func demoTrending(c *ui.Context, log func(string)) {
	topics := []social.Topic{
		{Tag: "#ceramics", Category: "Craft", Posts: 12480},
		{Tag: "#studiolife", Category: "Design", Posts: 8231},
		{Tag: "#glazetest", Category: "Craft", Posts: 3914},
		{Tag: "#brassandoak", Category: "Interiors", Posts: 1206},
		{Tag: "#productshoot", Category: "Photography", Posts: 987},
	}
	ui.Column(c).Width(380).Children(func() { social.Trending(c, topics, social.TrendingOptions{}) })
}

func demoPollPost(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	choices := ui.Local(root, "poll", func() []social.PollChoice {
		return []social.PollChoice{{Text: "Speckled oat", Votes: 412}, {Text: "Deep sea blue", Votes: 268}, {Text: "Matte charcoal", Votes: 155}}
	})
	voted := intState(root, "voted", -1)
	ends := ui.Local(root, "ends", func() time.Time { return c.Now().Add(18 * time.Hour) })
	root.Children(func() {
		if social.PollPost(c, "Which glaze should the spring mugs get?", choices, *ends, voted, social.PollPostOptions{}).Changed() {
			log("voted " + fmt.Sprint(*voted))
		}
	})
}

func demoMediaPost(c *ui.Context, log func(string)) {
	p := productShots()
	ui.Column(c).Width(440).Children(func() { social.MediaPost(c, []ui.ImageSource{p[1], p[5], p[3], p[4]}) })
}

func demoMentions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	text := textState(root, "text", "Loved the shoot, thanks @d")
	selected := textState(root, "selected", "")
	root.Children(func() {
		if social.Mentions(c, text, studio(), selected, social.MentionsOptions{}).Changed() {
			log(fmt.Sprintf("selected %s runes %d", *selected, len([]rune(*text))))
		}
	})
}

func demoUserCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := boolState(root, "open", true)
	root.Children(func() {
		if social.UserCard(c, studio()[4], open).Changed() {
			log("open " + fmt.Sprint(*open))
		}
	})
}

func demoFollowers(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	following := ui.Local(root, "following", func() map[string]bool { return map[string]bool{"daniel": true} })
	root.Children(func() {
		if social.Followers(c, studio()[1:5], following, social.FollowersOptions{}).Changed() {
			log("following " + fmt.Sprint(followed(*following)))
		}
	})
}

func demoProfileBadges(c *ui.Context, log func(string)) {
	badges := []social.Badge{
		{ID: "maker", Name: "Top maker", Icon: "award", Detail: "March 2026"},
		{ID: "streak", Name: "100-day streak", Icon: "flame", Detail: "Daily posts"},
		{ID: "early", Name: "Early member", Icon: "sparkles", Detail: "Since 2021"},
		{ID: "collab", Name: "Collaborator", Icon: "users", Detail: "12 shared drops"},
		{ID: "shoot", Name: "Featured", Icon: "star", Detail: "Editor's pick"},
		{ID: "mentor", Name: "Mentor", Icon: "medal", Detail: "8 apprentices"},
	}
	ui.Column(c).Width(420).Children(func() { social.ProfileBadges(c, badges) })
}

func demoReactionBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	reactions := ui.Local(root, "reactions", func() []social.Reaction {
		return []social.Reaction{
			{ID: "love", Name: "Love", Icon: "heart", Count: 24},
			{ID: "fire", Name: "Fire", Icon: "flame", Count: 11, Active: true},
			{ID: "clap", Name: "Applause", Icon: "party-popper", Count: 7},
			{ID: "idea", Name: "Insightful", Icon: "sparkles", Count: 3},
			{ID: "laugh", Name: "Funny", Icon: "laugh", Count: 2},
		}
	})
	root.Children(func() {
		if social.ReactionBar(c, reactions).Changed() {
			var parts []string
			for _, r := range *reactions {
				parts = append(parts, fmt.Sprintf("%s %d %v", r.ID, r.Count, r.Active))
			}
			log(strings.Join(parts, ", "))
		}
	})
}
