package store

import (
	"context"
	"database/sql"

	"github.com/c8112002/twitter_clone_go/utils"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/c8112002/twitter_clone_go/models"

	"github.com/c8112002/twitter_clone_go/entities"
)

type TweetStore struct {
	db  *sql.DB
	ctx context.Context
}

func NewTweetStore(db *sql.DB, ctx context.Context) *TweetStore {
	return &TweetStore{db: db, ctx: ctx}
}

func (ts *TweetStore) FetchTweets(maxID int, minID int, limit int) (*entities.Tweets, error) {
	tl, err := models.Tweets(
		qm.Where("id <= ?", maxID),
		qm.Where("id >= ?", minID),
		qm.Where("deleted_at is null"),
		qm.OrderBy(models.TweetColumns.ID+" desc"),
		qm.Limit(limit),
		qm.Load(qm.Rels(models.TweetRels.User, models.UserRels.Followers, models.FollowRels.Followee)),
		qm.Load(qm.Rels(models.TweetRels.User, models.UserRels.Followees, models.FollowRels.Follower)),
		qm.Load(qm.Rels(models.TweetRels.Likes, models.LikeRels.User, models.UserRels.Followers, models.FollowRels.Followee)),
		qm.Load(qm.Rels(models.TweetRels.Likes, models.LikeRels.User, models.UserRels.Followees, models.FollowRels.Follower)),
	).All(ts.ctx, ts.db)

	tweets := entities.Tweets{}

	for _, t := range tl {
		tweet := createTweet(*t)
		tweets = append(tweets, tweet)
	}

	return &tweets, err
}

func (ts *TweetStore) FetchFirstTweet() (*entities.Tweet, error) {
	t, err := models.Tweets(
		qm.Where("deleted_at is null"),
		qm.OrderBy(models.TweetColumns.ID),
		qm.Limit(1),
		qm.Load(qm.Rels(models.TweetRels.User, models.UserRels.Followers, models.FollowRels.Followee)),
		qm.Load(qm.Rels(models.TweetRels.User, models.UserRels.Followees, models.FollowRels.Follower)),
		qm.Load(qm.Rels(models.TweetRels.Likes, models.LikeRels.User, models.UserRels.Followers, models.FollowRels.Followee)),
		qm.Load(qm.Rels(models.TweetRels.Likes, models.LikeRels.User, models.UserRels.Followees, models.FollowRels.Follower)),
	).One(ts.ctx, ts.db)

	if err != nil {
		return nil, err
	}

	return createTweet(*t), nil
}

func (ts *TweetStore) FindTweet(id uint) (*entities.Tweet, error) {
	t, err := models.Tweets(
		qm.Where("id=?", id),
		qm.Limit(1),
		qm.Load(qm.Rels(models.TweetRels.User, models.UserRels.Followers, models.FollowRels.Followee)),
		qm.Load(qm.Rels(models.TweetRels.User, models.UserRels.Followees, models.FollowRels.Follower)),
		qm.Load(qm.Rels(models.TweetRels.Likes, models.LikeRels.User, models.UserRels.Followers, models.FollowRels.Followee)),
		qm.Load(qm.Rels(models.TweetRels.Likes, models.LikeRels.User, models.UserRels.Followees, models.FollowRels.Follower)),
	).One(ts.ctx, ts.db)

	if err != nil {
		return nil, err
	}

	return createTweet(*t), nil
}

func (ts *TweetStore) CreateTweet(tweet string, userID uint) (*entities.Tweet, error) {
	var m models.Tweet
	m.Tweet = tweet
	m.UserID = userID
	m.CreatedAt = utils.Now()
	m.UpdatedAt = utils.Now()
	err := m.Insert(ts.ctx, ts.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	t, err := ts.FindTweet(m.ID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (ts *TweetStore) Like(t *entities.Tweet, userID uint) (*entities.Tweet, error) {
	var like models.Like
	like.TweetID = t.ID
	like.UserID = userID
	like.CreatedAt = utils.Now()
	err := like.Insert(ts.ctx, ts.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	tweet, err := ts.FindTweet(t.ID)
	if err != nil {
		return nil, err
	}
	return tweet, nil
}

func (ts *TweetStore) Unlike(t *entities.Tweet, userID uint) (*entities.Tweet, error) {
	_, err := models.Likes(qm.Where(models.LikeColumns.UserID+"=?", userID), qm.Where(models.LikeColumns.TweetID+"=?", t.ID)).DeleteAll(ts.ctx, ts.db)
	if err != nil {
		return nil, err
	}

	tweet, err := ts.FindTweet(t.ID)
	if err != nil {
		return nil, err
	}
	return tweet, nil
}
