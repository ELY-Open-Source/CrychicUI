package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/social"
	"github.com/egoist/mygo/ui"
)

var socialEntries2 = []Entry{
	{ID: 1209, Name: "Repost", Category: "social", Summary: "Repost and quote menu; cue Pop.", Demo: demoSocialRepost},
	{ID: 1210, Name: "Bookmark", Category: "social", Summary: "Saved toggle; cue Toggle.", Demo: demoSocialBookmark},
	{ID: 1211, Name: "Community", Category: "social", Summary: "Community card with join; cue Toggle.", Demo: demoSocialCommunity},
	{ID: 1212, Name: "GroupCard", Category: "social", Summary: "Group with members and privacy.", Demo: demoSocialGroupCard},
	{ID: 1213, Name: "Invite", Category: "social", Summary: "Invite link and contacts; cue Pop.", Demo: demoSocialInvite},
	{ID: 1214, Name: "Ranking", Category: "social", Summary: "Weekly leaderboard.", Demo: demoSocialRanking},
	{ID: 1215, Name: "Review", Category: "social", Summary: "Star review with photos.", Demo: demoSocialReview},
	{ID: 1216, Name: "Testimonial", Category: "social", Summary: "Quote, name, role.", Demo: demoSocialTestimonial},
	{ID: 1217, Name: "Question", Category: "social", Summary: "Question with votes and answers; cue Tick.", Demo: demoSocialQuestion},
	{ID: 1218, Name: "VoteArrows", Category: "social", Summary: "Up, score, down; cue Tick.", Demo: demoSocialVoteArrows},
	{ID: 1219, Name: "NotificationsFeed", Category: "social", Summary: "Likes and follows grouped.", Demo: demoSocialNotificationsFeed},
	{ID: 1220, Name: "ProfileComplete", Category: "social", Summary: "Steps to finish a profile.", Demo: demoSocialProfileComplete},
	{ID: 1221, Name: "Verified", Category: "social", Summary: "Check badge with tooltip.", Demo: demoSocialVerified},
	{ID: 1222, Name: "SocialLinks", Category: "social", Summary: "Icon links row; cue Tap.", Demo: demoSocialLinks},
	{ID: 1223, Name: "Live", Category: "social", Summary: "Live stream card with viewers.", Demo: demoSocialLive},
	{ID: 1224, Name: "Gift", Category: "social", Summary: "Send a gift with counts; cue Pop.", Demo: demoSocialGift},
	{ID: 1225, Name: "ShareCount", Category: "social", Summary: "Share count that rolls.", Demo: demoSocialShareCount},
}

func demoSocialRepost(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	action := textState(root, "action", "")
	root.Children(func() {
		if social.Repost(c, action, social.RepostOptions{}).Changed() {
			log("action " + *action)
		}
	})
}

func demoSocialBookmark(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	saved := boolState(root, "saved", false)
	root.Children(func() {
		if social.Bookmark(c, saved, social.BookmarkOptions{}).Changed() {
			log("saved " + fmt.Sprint(*saved))
		}
	})
}

func demoSocialCommunity(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	joined := boolState(root, "joined", false)
	community := social.CommunityData{ID: "ui", Name: "UI Builders", Description: "Build native UI", Members: 42}
	root.Children(func() {
		if social.Community(c, community, joined, social.CommunityOptions{}).Changed() {
			log("joined " + fmt.Sprint(*joined))
		}
	})
}

func demoSocialGroupCard(c *ui.Context, log func(string)) {
	social.GroupCard(c, social.GroupData{ID: "ui", Name: "UI Builders", Privacy: "Public", Members: 42}, social.GroupCardOptions{})
}

func demoSocialInvite(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	invited := ui.Local(root, "invited", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		if social.Invite(c, "https://example.test/invite", []social.User{demoSocialUser()}, invited, social.InviteOptions{}).Changed() {
			log("invited " + fmt.Sprint((*invited)["ana"]))
		}
	})
}

func demoSocialRanking(c *ui.Context, log func(string)) {
	social.Ranking(c, []social.RankedUser{{User: demoSocialUser(), Score: 100}}, social.RankingOptions{})
}

func demoSocialReview(c *ui.Context, log func(string)) {
	p := photos()
	review := social.ReviewData{ID: "r1", Author: demoSocialUser(), Rating: 5, Text: "Great work", Photos: []ui.ImageSource{p[0], p[1]}}
	social.Review(c, review)
}

func demoSocialTestimonial(c *ui.Context, log func(string)) {
	social.Testimonial(c, "It feels native.", "Ana", "Designer")
}

func demoSocialQuestion(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	question := ui.Local(root, "question", func() social.QuestionData {
		return social.QuestionData{ID: "q1", Title: "Native UI?", Body: "How do I start?", Answers: 2}
	})
	vote := intState(root, "vote", 0)
	root.Children(func() {
		if social.Question(c, question, vote, social.QuestionOptions{}).Changed() {
			log("votes " + fmt.Sprint(question.Votes))
		}
	})
}

func demoSocialVoteArrows(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	score, vote := intState(root, "score", 0), intState(root, "vote", 0)
	root.Children(func() {
		if social.VoteArrows(c, score, vote, social.VoteArrowsOptions{}).Changed() {
			log("score " + fmt.Sprint(*score))
		}
	})
}

func demoSocialNotificationsFeed(c *ui.Context, log func(string)) {
	notices := []social.SocialNotification{{ID: "n1", Kind: "like", Actors: []string{"Ana"}, Time: "now"}, {ID: "n2", Kind: "follow", Actors: []string{"Ben"}, Time: "1m"}}
	social.NotificationsFeed(c, notices, social.NotificationsFeedOptions{})
}

func demoSocialProfileComplete(c *ui.Context, log func(string)) {
	steps := []social.ProfileStep{{ID: "photo", Title: "Add photo", Done: true}, {ID: "bio", Title: "Add bio"}}
	social.ProfileComplete(c, steps, social.ProfileCompleteOptions{})
}

func demoSocialVerified(c *ui.Context, log func(string)) {
	social.Verified(c, "Verified", social.VerifiedOptions{})
}

func demoSocialLinks(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := textState(root, "opened", "")
	links := []social.SocialLink{{ID: "gh", Name: "GitHub", Icon: "link", URL: "https://github.com"}}
	root.Children(func() {
		if social.SocialLinks(c, links, opened).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoSocialLive(c *ui.Context, log func(string)) {
	live := social.LiveData{ID: "l1", Title: "UI review", Streamer: "Ana", Viewers: 42, Cover: photos()[1]}
	social.Live(c, live, social.LiveOptions{})
}

func demoSocialGift(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	gifts := ui.Local(root, "gifts", func() []social.GiftData {
		return []social.GiftData{{ID: "flower", Name: "Flower", Icon: "flower"}}
	})
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if social.Gift(c, gifts, selected, social.GiftOptions{}).Changed() {
			log("count " + fmt.Sprint((*gifts)[0].Count))
		}
	})
}

func demoSocialShareCount(c *ui.Context, log func(string)) {
	social.ShareCount(c, 42, social.ShareCountOptions{})
}
