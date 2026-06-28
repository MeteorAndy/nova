package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"nova/internal/buildinfo"
)

const (
	githubAPIBase         = "https://api.github.com/repos"
	updateInstallTimeout  = 30 * time.Minute
	updateDownloadTimeout = 20 * time.Minute
)

type Service struct {
	repository     string
	currentVersion string
	httpClient     *http.Client
	executablePath string
	githubAPIBase  string
}

func NewService() *Service {
	exe, _ := os.Executable()
	return &Service{
		repository:     buildinfo.Repository,
		currentVersion: buildinfo.Version,
		httpClient:     &http.Client{Timeout: 60 * time.Second},
		executablePath: exe,
		githubAPIBase:  githubAPIBase,
	}
}

func (s *Service) Check(ctx context.Context) (CheckResult, error) {
	// Termux/Android: the in-app self-updater cannot safely run here — there is
	// no linux-arm64 release asset, and replacing the running binary risks
	// ETXTBSY on Android's overlay fs. Short-circuit so neither the auto-check
	// nor a manual "check now" offers an install; users rebuild in Termux.
	if termuxEnv() {
		return CheckResult{
			CurrentVersion: s.currentVersion,
			Platform:       platformKey(runtime.GOOS, runtime.GOARCH),
			Message:        "Termux/Android 暂不支持应用内更新，请在 Termux 内重新构建以升级",
		}, nil
	}
	release, err := s.latestRelease(ctx)
	if err != nil {
		return CheckResult{}, err
	}
	platform := platformKey(runtime.GOOS, runtime.GOARCH)
	asset := selectAsset(release.Assets, platform)
	current := s.currentVersion
	latest := normalizeVersion(release.TagName)
	updateAvailable := !isDevVersion(current) && latest != "" && compareVersions(current, latest) < 0
	result := CheckResult{
		CurrentVersion:  current,
		LatestVersion:   latest,
		UpdateAvailable: updateAvailable,
		CanInstall:      updateAvailable && asset != nil,
		Platform:        platform,
		ReleaseURL:      release.HTMLURL,
		PublishedAt:     release.PublishedAt,
		ReleaseNotes:    release.Body,
		Message:         "当前已是最新版本",
	}
	if asset != nil {
		result.Asset = &Asset{Name: asset.Name, Size: asset.Size, DownloadURL: asset.DownloadURL, BrowserDownloadURL: asset.BrowserDownloadURL}
	}
	switch {
	case isDevVersion(current):
		result.Message = "开发版本不支持应用内安装更新，请使用 Release 包运行后再检查"
	case latest == "":
		result.Message = "GitHub Release 未提供版本号"
	case !updateAvailable:
		result.Message = "当前已是最新版本"
	case asset == nil:
		result.Message = fmt.Sprintf("找到新版本，但没有匹配当前平台的安装包: %s", platform)
	default:
		result.Message = "发现可用更新"
	}
	return result, nil
}

// termuxEnv reports whether Nova is running under Termux on Android. Termux
// sets PREFIX to a path containing "com.termux" and HOME under
// /data/data/com.termux; both are absent on desktop Linux.
func termuxEnv() bool {
	return strings.Contains(os.Getenv("PREFIX"), "com.termux") || strings.HasPrefix(os.Getenv("HOME"), "/data/data/com.termux")
}

func (s *Service) latestRelease(ctx context.Context) (githubRelease, error) {
	url := s.githubLatestReleaseURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return githubRelease{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "nova-update-checker")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return githubRelease{}, fmt.Errorf("检查 GitHub Release 失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return githubRelease{}, fmt.Errorf("检查 GitHub Release 失败: HTTP %d %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return githubRelease{}, fmt.Errorf("解析 GitHub Release 响应失败: %w", err)
	}
	return release, nil
}

type githubRelease struct {
	TagName     string        `json:"tag_name"`
	HTMLURL     string        `json:"html_url"`
	Body        string        `json:"body"`
	PublishedAt time.Time     `json:"published_at"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	DownloadURL        string `json:"url"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func platformKey(goos, goarch string) string {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}
	return goos + "-" + arch
}

func selectAsset(assets []githubAsset, platform string) *githubAsset {
	for i := range assets {
		name := strings.ToLower(assets[i].Name)
		if strings.Contains(name, strings.ToLower(platform)) && strings.HasPrefix(name, "nova-") &&
			(strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip")) {
			return &assets[i]
		}
	}
	return nil
}

func selectChecksumAsset(assets []githubAsset) *githubAsset {
	for i := range assets {
		if strings.EqualFold(assets[i].Name, "checksums.txt") {
			return &assets[i]
		}
	}
	return nil
}

func (s *Service) githubLatestReleaseURL() string {
	base := strings.TrimRight(s.githubAPIBase, "/")
	if base == "" {
		base = githubAPIBase
	}
	return base + "/" + strings.Trim(s.repository, "/") + "/releases/latest"
}

func (s *Service) downloadHTTPClient() *http.Client {
	if s.httpClient == nil {
		return &http.Client{}
	}
	client := *s.httpClient
	client.Timeout = 0
	return &client
}

func updateAssetDownloadURL(asset *Asset) string {
	if asset == nil {
		return ""
	}
	if strings.TrimSpace(asset.BrowserDownloadURL) != "" {
		return asset.BrowserDownloadURL
	}
	return asset.DownloadURL
}

func githubAssetDownloadURL(asset githubAsset) string {
	if strings.TrimSpace(asset.BrowserDownloadURL) != "" {
		return asset.BrowserDownloadURL
	}
	return asset.DownloadURL
}
