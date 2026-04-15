package config

import (
	"os"
	"path/filepath"
)

// ScannedProject はスキャンで発見された未登録プロジェクト
type ScannedProject struct {
	Name string
	Path string
}

// ScanResult はスキャン結果
type ScanResult struct {
	New     []ScannedProject
	Skipped []ScannedProject
}

// ScanNewProjects はスキャンディレクトリから未登録プロジェクトを検出する
func ScanNewProjects(cfg *Config) (ScanResult, error) {
	scanDir := cfg.Settings.ScanDirectory
	if scanDir == "" {
		// デフォルト: ホームディレクトリ
		home, _ := os.UserHomeDir()
		scanDir = home
	}
	scanDir = ExpandPath(scanDir)

	entries, err := os.ReadDir(scanDir)
	if err != nil {
		return ScanResult{}, err
	}

	existing := map[string]bool{}
	for _, p := range cfg.Projects {
		existing[p.Name] = true
	}

	skippedSet := map[string]bool{}
	for _, p := range cfg.SkippedPaths {
		skippedSet[p] = true
	}

	var result ScanResult
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) == 0 || name[0] == '.' {
			continue
		}
		fullPath := filepath.Join(scanDir, name)
		if existing[name] {
			continue
		}
		sp := ScannedProject{Name: name, Path: fullPath}
		if skippedSet[fullPath] {
			result.Skipped = append(result.Skipped, sp)
		} else {
			result.New = append(result.New, sp)
		}
	}
	return result, nil
}
