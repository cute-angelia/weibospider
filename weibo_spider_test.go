package weibospider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetUserInfoUsesWeiboAjaxCookieAndParsesUser(t *testing.T) {
	var gotCookie string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ajax/profile/info" {
			t.Fatalf("path = %q, want /ajax/profile/info", r.URL.Path)
		}
		if r.URL.Query().Get("uid") != "2993720115" {
			t.Fatalf("uid query = %q", r.URL.Query().Get("uid"))
		}
		gotCookie = r.Header.Get("Cookie")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ok": true,
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"id":                2993720115,
					"screen_name":       "tester",
					"verified":          true,
					"verified_type":     0,
					"verified_reason":   "reason",
					"description":       "desc",
					"gender":            "f",
					"followers_count":   123,
					"friends_count":     45,
					"profile_image_url": "https://example.test/avatar.jpg",
				},
			},
		})
	}))
	defer server.Close()

	wb := NewWeiboSpider(WithCookie("SUB=abc"), withBaseURL(server.URL), WithDelay(0))
	user, err := wb.GetUserInfo(2993720115)
	if err != nil {
		t.Fatalf("GetUserInfo() error = %v", err)
	}

	if gotCookie != "SUB=abc" {
		t.Fatalf("Cookie header = %q, want SUB=abc", gotCookie)
	}
	if user.ID != 2993720115 || user.Name != "tester" || user.FollowersCount != 123 || user.FollowCount != 45 {
		t.Fatalf("parsed user = %#v", user)
	}
}

func TestGetUserPostsUsesSearchProfileParsesTextRawAndLongText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ajax/statuses/searchProfile":
			if r.URL.Query().Get("uid") != "2993720115" || r.URL.Query().Get("page") != "2" {
				t.Fatalf("query = %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ok": true,
				"data": map[string]interface{}{
					"list": []map[string]interface{}{
						{
							"id":              501,
							"mid":             501,
							"mblogid":         "Nabc123",
							"created_at":      "Thu Aug 13 10:00:00 +0800 2026",
							"text_raw":        "short text",
							"reposts_count":   1,
							"comments_count":  2,
							"attitudes_count": 3,
							"pic_num":         0,
							"page_info": map[string]interface{}{
								"type": 11,
							},
							"user": map[string]interface{}{
								"id":          2993720115,
								"screen_name": "tester",
							},
						},
						{
							"mid":             "502",
							"mblogid":         "Nlong456",
							"text_raw":        "truncated",
							"isLongText":      true,
							"reposts_count":   4,
							"comments_count":  5,
							"attitudes_count": 6,
							"pic_num":         1,
							"user": map[string]interface{}{
								"id":          2993720115,
								"screen_name": "tester",
							},
						},
					},
				},
			})
		case "/ajax/statuses/longtext":
			if r.URL.Query().Get("id") != "Nlong456" {
				t.Fatalf("longtext id = %q", r.URL.Query().Get("id"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ok":   true,
				"data": map[string]interface{}{"longTextContent": "full long text"},
			})
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	wb := NewWeiboSpider(WithCookie("SUB=abc"), withBaseURL(server.URL), WithDelay(0), WithLongText(true))
	posts, err := wb.GetUserPosts(2993720115, 2)
	if err != nil {
		t.Fatalf("GetUserPosts() error = %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("len(posts) = %d, want 2", len(posts))
	}
	if posts[0].ID != "501" || posts[0].MID != "501" || posts[0].Text != "short text" {
		t.Fatalf("first post = %#v", posts[0])
	}
	if posts[0].URL != "https://weibo.com/2993720115/Nabc123" {
		t.Fatalf("first URL = %q", posts[0].URL)
	}
	if posts[1].Text != "full long text" {
		t.Fatalf("long text = %q", posts[1].Text)
	}
}

func TestMissingCookieReturnsHelpfulError(t *testing.T) {
	t.Setenv("WEIBO_COOKIE", "")

	wb := NewWeiboSpider(WithDelay(0))
	_, err := wb.GetUserPosts(2993720115, 1)
	if err == nil || !strings.Contains(err.Error(), "cookie") {
		t.Fatalf("error = %v, want cookie error", err)
	}
}

func TestPackageGetUserInfoUsesDefaultSpider(t *testing.T) {
	t.Setenv("WEIBO_COOKIE", "")

	_, err := GetUserInfo(2993720115)
	if err == nil || !strings.Contains(err.Error(), "cookie") {
		t.Fatalf("error = %v, want cookie error", err)
	}
}

func TestRandomSleepAllowsEqualBounds(t *testing.T) {
	start := time.Now()
	RandomSleep(0, 0)
	if time.Since(start) > time.Second {
		t.Fatal("RandomSleep(0, 0) took too long")
	}
}
