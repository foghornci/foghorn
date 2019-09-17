package bitbucketserver

type webhook struct {
	EventKey            string      `json:"eventKey"`
	Date                string      `json:"date"`
	Actor               user        `json:"actor"`
	PullRequest         pullRequest `json:"pullRequest"`
	PreviousTitle       string      `json:"previousTitle"`
	PreviousDescription string      `json:"previousDescription"`
	PreviousTarget      target      `json:"previousTarget"`
	RemovedReviewers    []reviewer  `json:"removedReviewers"`
	AddedReviewers      []reviewer  `json:"addedReviewers"`
	PreviousStatus      string      `json:"previousStatus"`
	Comment             comment     `json:"comment"`
	CommentParentID     int64       `json:"commentParentId"`
	PreviousComment     string      `json:"previousComment"`
	Repository          repository  `json:"repository"`
	Changes             []change    `json:"changes"`
	Old                 repository  `json:"old"`
	New                 repository  `json:"new"`
	Commit              string      `json:"commit"`
}

type pullRequest struct {
	ID           int64      `json:"id"`
	Version      string     `json:"version"`
	Title        string     `json:"title"`
	State        string     `json:"state"`
	Open         bool       `json:"open"`
	Closed       bool       `json:"closed"`
	CreatedDate  string     `json:"createdDate"`
	UpdatedDate  string     `json:"updatedDate"`
	FromRef      ref        `json:"fromRef"`
	ToRef        ref        `json:"toRef"`
	Locked       bool       `json:"locked"`
	Author       author     `json:"author"`
	Reviewers    []reviewer `json:"reviewers"`
	Participants []user     `json:"participants"`
	Links        links      `json:"links"`
}

type ref struct {
	ID           string     `json:"id"`
	DisplayID    string     `json:"displayId"`
	LatestCommit string     `json:"latestCommit"`
	Repository   repository `json:"repository"`
	Public       bool       `json:"public"`
	Type         string     `json:"type"`
}

type author struct {
	User     user   `json:"user"`
	Role     string `json:"role"`
	Approved bool   `json:"approved"`
	Status   string `json:"status"`
}

type user struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int64  `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type reviewer struct {
	User     user   `json:"user"`
	Role     string `json:"role"`
	Approved string `json:"approved"`
	Status   string `json:"status"`
}

type participant struct {
	User               user   `json:"user"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Approved           bool   `json:"approved"`
	Status             string `json:"status"`
}

type repository struct {
	Slug          string  `json:"string"`
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	SCMID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       project `json:"project"`
	Public        bool    `json:"public"`
}

type project struct {
	Key    string `json:"key"`
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
	Type   string `json:"type"`
}

type links struct {
	Self linkSet `json:"self"`
}

type linkSet struct {
	Href string `json:"href"`
}

type target struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
}

type comment struct {
	Properties  properties `json:"properties"`
	ID          int64      `json:"id"`
	Version     string     `json:"version"`
	Text        string     `json:"text"`
	Author      author     `json:"author"`
	CreatedDate string     `json:"createDate"`
	UpdatedDate string     `json:"updatedDate"`
}

type properties struct {
	RepositoryID int64 `json:"repositoryId"`
}

type change struct {
	Ref      ref    `json:"ref"`
	RefID    string `json:"refId"`
	FromHash string `json:"fromHash"`
	ToHash   string `json:"toHash"`
	Type     string `json:"type"`
}
