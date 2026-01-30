package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var skippedDirs = []string{".git", "node_modules", "vendor"}

func isTempFile(name string) bool {
	base := filepath.Base(name)
	return strings.HasPrefix(base, ".") ||
		strings.HasSuffix(name, "~") ||
		strings.HasSuffix(name, ".swp")
}

func isSkippedDir(path string) bool {
	for _, d := range skippedDirs {
		if strings.Contains(path, d) {
			return true
		}
	}
	return false
}

func resolveWatchDirs(extraDirs string, basePath string) (watchList, extrasList []string) {
	watchList = []string{basePath}

	if extraDirs == "" {
		return watchList, nil
	}

	for _, e := range strings.Split(extraDirs, ",") {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		absExtra, err := filepath.Abs(e)
		if err != nil {
			log.Printf("Warning: Skipping invalid watch dir %s: %v", e, err)
			continue
		}
		extrasList = append(extrasList, absExtra)
		watchList = append(watchList, absExtra)
	}
	return watchList, extrasList
}

func isSourceFile(absPath string, sourceDirs []string) bool {
	for _, src := range sourceDirs {
		if strings.HasPrefix(absPath, src) {
			return true
		}
	}
	return false
}

func isOutputFile(absPath string, servedDir string) bool {
	return strings.HasPrefix(absPath, servedDir)
}

func watchFiles(watchDirs []string, sourceDirs []string, cmdStr string, servedDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = watcher.Close() }()

	for _, dir := range watchDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			log.Printf("Warning: couldn't resolve path %s: %v", dir, err)
			continue
		}

		err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if isSkippedDir(path) {
					return filepath.SkipDir
				}
				if err := watcher.Add(path); err != nil {
					log.Printf("Warning: couldn't watch %s: %v", path, err)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Warning: error walking path for %s: %v", absDir, err)
		} else {
			log.Printf("Watching for changes in: %s", absDir)
		}
	}

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if isTempFile(event.Name) {
				continue
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Printf("File changed: %s", event.Name)

				absEventPath, err := filepath.Abs(event.Name)
				if err != nil {
					log.Printf("Error resolving path %s: %v", event.Name, err)
					absEventPath = event.Name
				}

				isSource := isSourceFile(absEventPath, sourceDirs)
				isOutput := isOutputFile(absEventPath, servedDir)

				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(200*time.Millisecond, func() {
					if isSource && cmdStr != "" {
						log.Printf("Executing command: %s", cmdStr)
						cmd := exec.Command("sh", "-c", cmdStr)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						if err := cmd.Run(); err != nil {
							log.Printf("Command execution failed: %v", err)
						} else {
							log.Println("Command execution completed")
						}
						return
					}

					if isOutput {
						log.Println("Output changed, notifying clients...")
						notifyClients()
					}
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
