package gallery

import (
	"fmt"
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

func demoPost(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	state := ui.Local(root, "post", func() social.PostState { return social.PostState{} })
	root.Children(func() {
		if social.Post(c, demoSocialPost(), state, social.PostOptions{}).Changed() {
			log(fmt.Sprintf("liked %v saved %v", state.Liked, state.Saved))
		}
	})
}

func demoSocialFeed(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	fresh := boolState(root, "fresh", true)
	root.Children(func() {
		if social.Feed(c, []social.PostData{demoSocialPost()}, fresh, social.FeedOptions{}).Changed() {
			log("fresh " + fmt.Sprint(*fresh))
		}
	})
}

func demoPostComposer(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	draft := textState(root, "draft", "Hello")
	added := boolState(root, "added", false)
	root.Children(func() {
		form, send := social.PostComposer(c, draft, photos()[0], added, social.PostComposerOptions{})
		if send.Changed() {
			log("sent")
		} else if form.Changed() {
			log(fmt.Sprintf("runes %d media %v", len([]rune(*draft)), *added))
		}
	})
}

func demoPostActions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	state := ui.Local(root, "actions", func() social.PostState { return social.PostState{} })
	root.Children(func() {
		if social.PostActions(c, state, social.PostActionsOptions{}).Changed() {
			log(fmt.Sprintf("liked %v action %s", state.Liked, state.Action))
		}
	})
}

func demoCommentThread(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	draft := textState(root, "draft", "Thanks")
	comments := []social.Comment{{ID: "c1", Text: "First", Author: demoSocialUser(), Replies: []social.Comment{{ID: "c2", Text: "Nested", Author: demoSocialUser()}}}}
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
	root := ui.Column(c)
	following := boolState(root, "following", false)
	root.Children(func() {
		if social.ProfileHeader(c, demoSocialUser(), following, social.ProfileHeaderOptions{}).Changed() {
			log("following " + fmt.Sprint(*following))
		}
	})
}

func demoProfileStats(c *ui.Context, log func(string)) {
	social.ProfileStats(c, 12, 1200, 40, social.ProfileStatsOptions{})
}

func demoWhoToFollow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	following := ui.Local(root, "following", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		if social.WhoToFollow(c, []social.User{demoSocialUser()}, following, social.WhoToFollowOptions{}).Changed() {
			log("following " + fmt.Sprint((*following)["ana"]))
		}
	})
}

func demoStoryRing(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := textState(root, "opened", "")
	stories := []social.Story{{ID: "s1", User: demoSocialUser()}}
	root.Children(func() {
		if social.StoryRing(c, stories, opened).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoHashtag(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		if social.Hashtag(c, "#mygo", 42, opened, social.HashtagOptions{}).Changed() {
			log("opened " + fmt.Sprint(*opened))
		}
	})
}

func demoTrending(c *ui.Context, log func(string)) {
	social.Trending(c, []social.Topic{{Tag: "#mygo", Posts: 42}, {Tag: "#ui", Posts: 20}}, social.TrendingOptions{})
}

func demoPollPost(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	choices := ui.Local(root, "poll", func() []social.PollChoice {
		return []social.PollChoice{{Text: "Yes", Votes: 2}, {Text: "No", Votes: 1}}
	})
	voted := intState(root, "voted", -1)
	ends := ui.Local(root, "ends", func() time.Time { return c.Now().Add(time.Hour) })
	root.Children(func() {
		if social.PollPost(c, "Ship?", choices, *ends, voted, social.PollPostOptions{}).Changed() {
			log("voted " + fmt.Sprint(*voted))
		}
	})
}

func demoMediaPost(c *ui.Context, log func(string)) {
	p := photos()
	images := make([]ui.ImageSource, len(p))
	for i, image := range p {
		images[i] = image
	}
	social.MediaPost(c, images)
}

func demoMentions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	text := textState(root, "text", "Hi @a")
	selected := textState(root, "selected", "")
	root.Children(func() {
		if social.Mentions(c, text, []social.User{demoSocialUser()}, selected, social.MentionsOptions{}).Changed() {
			log(fmt.Sprintf("selected %s runes %d", *selected, len([]rune(*text))))
		}
	})
}

func demoUserCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := boolState(root, "open", false)
	root.Children(func() {
		if social.UserCard(c, demoSocialUser(), open).Changed() {
			log("open " + fmt.Sprint(*open))
		}
	})
}

func demoFollowers(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	following := ui.Local(root, "following", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		if social.Followers(c, []social.User{demoSocialUser()}, following, social.FollowersOptions{}).Changed() {
			log("following " + fmt.Sprint((*following)["ana"]))
		}
	})
}

func demoProfileBadges(c *ui.Context, log func(string)) {
	social.ProfileBadges(c, []social.Badge{{ID: "b1", Name: "Creator", Icon: "award"}})
}

func demoReactionBar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	reactions := ui.Local(root, "reactions", func() []social.Reaction {
		return []social.Reaction{{ID: "heart", Name: "Heart", Icon: "heart", Count: 2}}
	})
	root.Children(func() {
		if social.ReactionBar(c, reactions).Changed() {
			log("count " + fmt.Sprint((*reactions)[0].Count))
		}
	})
}
