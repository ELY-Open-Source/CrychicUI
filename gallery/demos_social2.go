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
	root := ui.Column(c).Width(360)
	action := textState(root, "action", "")
	root.Children(func() {
		if social.Repost(c, action, social.RepostOptions{}).Changed() {
			log("action " + *action)
		}
	})
}

func demoSocialBookmark(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	saved := boolState(root, "saved", false)
	root.Children(func() {
		if social.Bookmark(c, saved, social.BookmarkOptions{}).Changed() {
			log("saved " + fmt.Sprint(*saved))
		}
	})
}

func demoSocialCommunity(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	joined := boolState(root, "joined", false)
	community := social.CommunityData{ID: "potters", Name: "Studio Potters", Icon: "flame", Members: 18342, Cover: productShots()[1],
		Description: "Glaze recipes, kiln logs and weekly critiques from people who throw every day.",
		Faces:       []string{"Priya Nair", "Jonas Berg", "Sofia Rossi", "Daniel Okafor"}}
	root.Children(func() {
		if social.Community(c, community, joined, social.CommunityOptions{}).Changed() {
			log("joined " + fmt.Sprint(*joined))
		}
	})
}

func demoSocialGroupCard(c *ui.Context, log func(string)) {
	groups := []social.GroupData{
		{ID: "glaze", Name: "Glaze Chemistry Club", Privacy: "Public", Icon: "flask-conical", Members: 4210,
			Faces: []string{"Priya Nair", "Maya Chen", "Lucas Martin", "Jonas Berg"}},
		{ID: "kiln", Name: "Kiln & Co. team", Privacy: "Private", Private: true, Icon: "lock", Members: 12,
			Faces: []string{"Maya Chen", "Daniel Okafor", "Sofia Rossi"}},
		{ID: "wood", Name: "Oak and Brass Makers", Privacy: "Public", Icon: "hammer", Members: 987,
			Faces: []string{"Lucas Martin", "Jonas Berg"}},
	}
	ui.Column(c).Width(440).Gap(10).Children(func() {
		for _, group := range groups {
			c.Key(group.ID)
			social.GroupCard(c, group, social.GroupCardOptions{})
		}
	})
}

func demoSocialInvite(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	invited := ui.Local(root, "invited", func() map[string]bool { return map[string]bool{"jonas": true} })
	root.Children(func() {
		invite, copy := social.Invite(c, "https://kiln.co/join/spring-drop", studio()[1:6], invited, social.InviteOptions{})
		if invite.Changed() {
			log("invited " + fmt.Sprint(followed(*invited)))
		}
		if copy.Changed() {
			log("copied link")
		}
	})
}

func demoSocialRanking(c *ui.Context, log func(string)) {
	u := studio()
	rows := []social.RankedUser{{User: u[2], Score: 2480}, {User: u[0], Score: 2215}, {User: u[3], Score: 1830},
		{User: u[5], Score: 1204}, {User: u[1], Score: 960}}
	ui.Column(c).Width(420).Children(func() { social.Ranking(c, rows, social.RankingOptions{}) })
}

func demoSocialReview(c *ui.Context, log func(string)) {
	p := productShots()
	review := social.ReviewData{ID: "r1", Author: studio()[4], Rating: 4, Time: "3d", Title: "Gorgeous glaze, a little heavy",
		Text:   "The speckled oat finish looks even better in person and the handle is comfortable. It is heavier than I expected for daily coffee.",
		Photos: []ui.ImageSource{p[1], p[5], p[3], p[4], p[0]}}
	ui.Column(c).Width(440).Children(func() { social.Review(c, review) })
}

func demoSocialTestimonial(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		social.Testimonial(c, "Kiln & Co. made our studio feel like a real brand. The lamp sold out in four hours.",
			"Jonas Berg", "Founder, Berg Furniture")
	})
}

func demoSocialQuestion(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	question := ui.Local(root, "question", func() social.QuestionData {
		return social.QuestionData{ID: "q1", Votes: 41, Answers: 6, Tags: []string{"glaze", "firing", "stoneware"},
			Title: "Why does my speckled glaze crawl on the rim?",
			Body:  "Cone 6 oxidation, two coats dipped. The rim pulls back into beads after the glaze firing."}
	})
	vote := intState(root, "vote", 0)
	root.Children(func() {
		if social.Question(c, question, vote, social.QuestionOptions{}).Changed() {
			log("votes " + fmt.Sprint(question.Votes))
		}
	})
}

func demoSocialVoteArrows(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	score, vote := intState(root, "score", 1284), intState(root, "vote", 0)
	root.Children(func() {
		if social.VoteArrows(c, score, vote, social.VoteArrowsOptions{}).Changed() {
			log("score " + fmt.Sprint(*score))
		}
	})
}

func demoSocialNotificationsFeed(c *ui.Context, log func(string)) {
	notices := []social.SocialNotification{
		{ID: "n1", Kind: "like", Time: "2m", Actors: []string{"Priya Nair", "Daniel Okafor", "Sofia Rossi", "Jonas Berg"}},
		{ID: "n2", Kind: "like", Time: "1h", Actors: []string{"Lucas Martin"}},
		{ID: "n3", Kind: "follow", Time: "3h", Actors: []string{"Sofia Rossi", "Jonas Berg"}},
		{ID: "n4", Kind: "follow", Time: "1d", Actors: []string{"Daniel Okafor"}},
	}
	ui.Column(c).Width(440).Children(func() { social.NotificationsFeed(c, notices, social.NotificationsFeedOptions{}) })
}

func demoSocialProfileComplete(c *ui.Context, log func(string)) {
	steps := []social.ProfileStep{
		{ID: "photo", Title: "Add a profile photo", Done: true},
		{ID: "bio", Title: "Write a short bio", Done: true},
		{ID: "links", Title: "Link your shop and portfolio", Done: true},
		{ID: "post", Title: "Share your first post"},
		{ID: "follow", Title: "Follow five makers"},
	}
	ui.Column(c).Width(420).Children(func() { social.ProfileComplete(c, steps, social.ProfileCompleteOptions{}) })
}

func demoSocialVerified(c *ui.Context, log func(string)) {
	ui.Row(c).Width(320).Gap(6).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "Maya Chen").Bold()
		social.Verified(c, "Verified: Maya Chen", social.VerifiedOptions{})
	})
}

func demoSocialLinks(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	opened := textState(root, "opened", "")
	links := []social.SocialLink{
		{ID: "site", Name: "Website", Icon: "globe", URL: "https://kiln.co"},
		{ID: "photos", Name: "Photos", Icon: "camera", URL: "https://photos.kiln.co"},
		{ID: "videos", Name: "Videos", Icon: "play", URL: "https://video.kiln.co"},
		{ID: "journal", Name: "Journal", Icon: "rss", URL: "https://kiln.co/journal"},
		{ID: "mail", Name: "Email", Icon: "mail", URL: "mailto:studio@kiln.co"},
	}
	root.Children(func() {
		if social.SocialLinks(c, links, opened).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoSocialLive(c *ui.Context, log func(string)) {
	live := social.LiveData{ID: "l1", Title: "Throwing the spring mugs, live from the studio", Streamer: "Priya Nair", Viewers: 1204,
		Cover: productShots()[1]}
	ui.Column(c).Width(420).Children(func() { social.Live(c, live, social.LiveOptions{}) })
}

func demoSocialGift(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	gifts := ui.Local(root, "gifts", func() []social.GiftData {
		return []social.GiftData{
			{ID: "heart", Name: "Heart", Icon: "heart", Count: 128},
			{ID: "coffee", Name: "Coffee", Icon: "coffee", Count: 42},
			{ID: "flower", Name: "Flower", Icon: "flower", Count: 17},
			{ID: "gem", Name: "Gem", Icon: "gem", Count: 4},
			{ID: "rocket", Name: "Rocket", Icon: "rocket", Count: 2},
			{ID: "cake", Name: "Cake", Icon: "cake", Count: 9},
		}
	})
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if social.Gift(c, gifts, selected, social.GiftOptions{}).Changed() {
			g := (*gifts)[*selected]
			log(fmt.Sprintf("%s %d", g.ID, g.Count))
		}
	})
}

func demoSocialShareCount(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Gap(10).Children(func() {
		social.ShareCount(c, 3214, social.ShareCountOptions{})
		social.ShareCount(c, 1, social.ShareCountOptions{})
	})
}
