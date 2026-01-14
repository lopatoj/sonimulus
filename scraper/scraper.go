package scraper

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"lopa.to/sonimulus/env"
	"lopa.to/sonimulus/repository"
)

const (
	imageRegexp = `background\-image:\surl\("([^>]+)"\);`
)

const (
	nameSelector          = "#content > div > div.l-user-hero.sc-px-2x > div > div.profileHeader__info > div > div.profileHeaderInfo__content.sc-media-content > h2.profileHeaderInfo__userName"
	imageSelector         = "#content > div > div.l-user-hero.sc-px-2x > div > div.profileHeader__info > div > div.profileHeaderInfo__avatar.sc-media-image.sc-mr-4x > div > span.sc-artwork"
	verifiedSelector      = "#content > div > div.l-user-hero.sc-px-2x > div > div.profileHeader__info > div > div.profileHeaderInfo__content.sc-media-content > h2 > div > span.verifiedBadge"
	trackCountSelector    = "#content > div > div.l-fluid-fixed > div.l-sidebar-right.l-user-sidebar-right > div > article.infoStats > table > tbody > tr > td:nth-child(3) > a > div"
	followerCountSelector = "#content > div > div.l-fluid-fixed > div.l-sidebar-right.l-user-sidebar-right > div > article.infoStats > table > tbody > tr > td:nth-child(1) > a > div"
	artistPlanSelector    = "#content > div > div.l-user-hero.sc-px-2x > div > div.profileHeader__info > div > div.profileHeaderInfo__content.sc-media-content > h3.profileHeaderInfo__additional > a.creatorBadge"
	handleLinkSelector    = "#content > div > div > div.l-main.g-main-scroll-area > div > div > ul > li > div > div.userBadgeListItem__title.sc-mt-2x.sc-mb-0\\.25x > a"
)

type HandleDepth struct {
	Handle string
	Depth  int
}

type Result struct {
	Handles     []string
	ParentDepth int
}

type Scraper struct {
	maxDepth int
	env      env.Env
	browser  *rod.Browser
}

func NewScraper(maxDepth int, e env.Env) *Scraper {
	u := launcher.New().
		NoSandbox(true).
		MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	return &Scraper{
		maxDepth: maxDepth,
		env:      e,
		browser:  browser,
	}
}

func (s *Scraper) ScrapePeopleConcurrent(
	numWorkers int,
	rootHandle string,
	onPerson func(handle, name, imageUrl string, verified bool, plan repository.Plan, trackCount int64) (id int64),
	onFollows func(followerId int64, followeeHandles []string),
) {
	visited := sync.Map{}
	wg := sync.WaitGroup{}
	queue := make(chan HandleDepth, 10000)
	results := make(chan Result, 1000)

	for i := range numWorkers {
		go func(workerId int) {
			for handleDepth := range queue {
				slog.Info("working new scraping job", "id", workerId, "handle", handleDepth.Handle, "depth", handleDepth.Depth)

				getFollows := handleDepth.Depth < s.maxDepth
				err := s.ScrapePerson(
					handleDepth.Handle,
					getFollows,
					onPerson,
					func(followerId int64, followeeHandles []string) {
						onFollows(followerId, followeeHandles)
						wg.Add(len(followeeHandles))
						results <- Result{
							Handles:     followeeHandles,
							ParentDepth: handleDepth.Depth,
						}
					},
				)
				if err != nil {
					slog.Error("error scraping person", "error", err)
				}
				wg.Done()
			}
		}(i)
	}

	go func() {
		for res := range results {
			for _, followeeHandle := range res.Handles {
				if _, seen := visited.LoadOrStore(followeeHandle, true); !seen {
					queue <- HandleDepth{
						Handle: followeeHandle,
						Depth:  res.ParentDepth + 1,
					}
				}
			}
		}
	}()

	visited.Store(rootHandle, true)
	wg.Add(1)
	queue <- HandleDepth{
		Handle: rootHandle,
		Depth:  0,
	}

	wg.Wait()
	close(queue)
	close(results)
}

func (s *Scraper) ScrapePerson(
	handle string,
	getFollows bool,
	onPerson func(handle, name, imageUrl string, verified bool, plan repository.Plan, trackCount int64) (id int64),
	onFollows func(followerId int64, followeeHandles []string),
) error {
	var (
		name       string
		imageUrl   string
		verified   bool = false
		plan       repository.Plan
		trackCount int64
	)

	userUrl := s.env.Soundcloud.URL + handle
	slog.Info("scraping user", "handle", handle)

	userPage := s.browser.MustPage(userUrl)
	defer userPage.Close()

	// Scrape user info
	name = strings.Trim(userPage.MustElement(nameSelector).MustText(), " ")

	style := userPage.MustElement(imageSelector).MustAttribute("style")
	if style == nil {
		slog.Error("user has no image", "handle", handle)
		return fmt.Errorf("user has no image")
	}
	matches := regexp.MustCompile(imageRegexp).FindStringSubmatch(*style)
	if len(matches) < 2 {
		slog.Error("user has no image", "handle", handle)
	} else {
		imageUrl = matches[1]
	}
	slog.Info("user image", "url", imageUrl)

	exists, _, err := userPage.Has(verifiedSelector)
	if err != nil {
		slog.Error("failed to find verified element", "handle", handle, "error", err)
		return fmt.Errorf("failed to find verified element")
	}
	if exists {
		verified = true
	}
	slog.Info("user verified", "verified", verified)

	exists, artistPlanElement, err := userPage.Has(artistPlanSelector)
	if err != nil {
		slog.Error("failed to find artist plan", "handle", handle, "error", err)
		return fmt.Errorf("failed to find artist plan")
	}
	if exists {
		title := artistPlanElement.MustAttribute("title")
		if title != nil {
			switch *title {
			case "Artist":
				plan = repository.PlanArtist
			case "Artist Pro":
				plan = repository.PlanArtistPro
			default:
				plan = repository.PlanNone
			}
		}
	}
	slog.Info("user plan", "plan", plan)

	count, err := strconv.Atoi(strings.ReplaceAll(userPage.MustElement(trackCountSelector).MustText(), ",", ""))
	if err != nil {
		slog.Error("failed to parse track count", "handle", handle, "error", err)
		return fmt.Errorf("failed to parse track count")
	}
	trackCount = int64(count)
	slog.Info("user track count", "count", trackCount)

	id := onPerson(
		handle,
		name,
		imageUrl,
		verified,
		plan,
		trackCount,
	)

	if id < 0 || !getFollows {
		return nil
	}

	followingPage := s.browser.MustPage(userUrl, "following")
	defer followingPage.Close()
	var lastHeight float64
	for {
		currentHeight := followingPage.MustEval(`() => document.documentElement.scrollHeight`).Num()

		if currentHeight == lastHeight {
			time.Sleep(2 * time.Second)
			if followingPage.MustEval(`() => document.documentElement.scrollHeight`).Num() == currentHeight {
				break
			}
		}

		slog.Info("scrolling")

		followingPage.MustEval(`() => window.scrollTo(0, document.documentElement.scrollHeight)`)

		followingPage.MustWaitIdle()

		lastHeight = currentHeight

		time.Sleep(500 * time.Millisecond)
	}

	usersFollowed := followingPage.MustElements(handleLinkSelector)
	follows := make([]string, 0, len(usersFollowed))
	for _, userFollowed := range usersFollowed {
		h, _ := strings.CutPrefix(*userFollowed.MustAttribute("href"), "/")
		follows = append(follows, h)
		slog.Info("user followed", "handle", h)
	}
	onFollows(id, follows)

	return nil
}
