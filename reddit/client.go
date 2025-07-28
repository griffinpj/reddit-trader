package reddit

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

// UserProfile represents Reddit user information
type UserProfile struct {
    ID               string  `json:"id"`
    Name             string  `json:"name"`
    IconImg          string  `json:"icon_img"`
    SnoovatarImg     string  `json:"snoovatar_img"`
    Created          float64 `json:"created"`
    CreatedUTC       float64 `json:"created_utc"`
    LinkKarma        int     `json:"link_karma"`
    CommentKarma     int     `json:"comment_karma"`
    TotalKarma       int     `json:"total_karma"`
    IsFriend         bool    `json:"is_friend"`
    IsEmployee       bool    `json:"is_employee"`
    IsMod            bool    `json:"is_mod"`
    IsGold           bool    `json:"is_gold"`
    Verified         bool    `json:"verified"`
    HasVerifiedEmail bool    `json:"has_verified_email"`
    Description      string  `json:"subreddit.public_description"`
}

// Post represents a Reddit post
type Post struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"` // Full ID (t3_xxxxx)
    Title       string  `json:"title"`
    Author      string  `json:"author"`
    Subreddit   string  `json:"subreddit"`
    Score       int     `json:"score"`
    NumComments int     `json:"num_comments"`
    Created     float64 `json:"created"`
    CreatedUTC  float64 `json:"created_utc"`
    URL         string  `json:"url"`
    Permalink   string  `json:"permalink"`
    Selftext    string  `json:"selftext"`
    IsSelf      bool    `json:"is_self"`
    Over18      bool    `json:"over_18"`
    Spoiler     bool    `json:"spoiler"`
    Stickied    bool    `json:"stickied"`
}

// ListingResponse represents Reddit's listing response structure
type ListingResponse struct {
    Kind string `json:"kind"`
    Data struct {
        After    string `json:"after"`
        Before   string `json:"before"`
        Children []struct {
            Kind string          `json:"kind"`
            Data json.RawMessage `json:"data"`
        } `json:"children"`
    } `json:"data"`
}

// GetUserProfile fetches the authenticated user's profile information
func (c *Client) GetUserProfile(accessToken string) (*UserProfile, error) {
    req, err := http.NewRequest("GET", apiURL+"/api/v1/me", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Authorization", "Bearer "+accessToken)
    req.Header.Set("User-Agent", userAgent)

    // Execute request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Check status
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    // Parse response
    var profile UserProfile
    if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    return &profile, nil
}

// PostHistoryOptions contains options for filtering post history
type PostHistoryOptions struct {
    Limit      int      // Number of posts to retrieve (max 100)
    After      string   // Pagination - get posts after this ID
    Before     string   // Pagination - get posts before this ID
    Sort       string   // Sort order: "new", "hot", "top", "controversial"
    TimeFilter string   // Time filter for "top" and "controversial": "hour", "day", "week", "month", "year", "all"
    Subreddits []string // Filter by specific subreddits (applied client-side)
}

// GetUserPosts fetches the authenticated user's post history
func (c *Client) GetUserPosts(accessToken, username string, options PostHistoryOptions) ([]*Post, string, error) {
    // Build URL with query parameters
    params := url.Values{}
    if options.Limit > 0 && options.Limit <= 100 {
        params.Set("limit", fmt.Sprintf("%d", options.Limit))
    } else {
        params.Set("limit", "25") // Default
    }

    if options.After != "" {
        params.Set("after", options.After)
    }
    if options.Before != "" {
        params.Set("before", options.Before)
    }
    if options.Sort != "" {
        params.Set("sort", options.Sort)
    }
    if options.TimeFilter != "" && (options.Sort == "top" || options.Sort == "controversial") {
        params.Set("t", options.TimeFilter)
    }

    endpoint := fmt.Sprintf("%s/user/%s/submitted?%s", apiURL, username, params.Encode())

    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, "", fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Authorization", "Bearer "+accessToken)
    req.Header.Set("User-Agent", userAgent)

    // Execute request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, "", fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Check status
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    // Parse listing response
    var listing ListingResponse
    if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
        return nil, "", fmt.Errorf("failed to parse response: %w", err)
    }

    // Extract posts
    posts := make([]*Post, 0)
    for _, child := range listing.Data.Children {
        if child.Kind == "t3" { // t3 = link/post
            var post Post
            if err := json.Unmarshal(child.Data, &post); err != nil {
                continue // Skip malformed posts
            }

            // Apply subreddit filter if specified
            if len(options.Subreddits) > 0 {
                found := false
                for _, sub := range options.Subreddits {
                    if strings.EqualFold(post.Subreddit, sub) {
                        found = true
                        break
                    }
                }
                if !found {
                    continue
                }
            }

            posts = append(posts, &post)
        }
    }

    return posts, listing.Data.After, nil
}

// GetUserPostsBySubreddits is a convenience method that fetches all user posts from specific subreddits
func (c *Client) GetUserPostsBySubreddits(accessToken, username string, subreddits []string, limit int) ([]*Post, error) {
    allPosts := make([]*Post, 0)
    after := ""

    // Reddit API returns max 100 posts per request
    for len(allPosts) < limit {
        remaining := limit - len(allPosts)
        requestLimit := 100
        if remaining < 100 {
            requestLimit = remaining
        }

        posts, nextAfter, err := c.GetUserPosts(accessToken, username, PostHistoryOptions{
            Limit:      requestLimit,
            After:      after,
            Sort:       "new",
            Subreddits: subreddits,
        })

        if err != nil {
            return allPosts, err
        }

        allPosts = append(allPosts, posts...)

        // No more posts available
        if nextAfter == "" || len(posts) == 0 {
            break
        }

        after = nextAfter
    }

    return allPosts, nil
}

