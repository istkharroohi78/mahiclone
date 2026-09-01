package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"ANJALI/config" // Assuming tumhara config package yahan hai
	"ANJALI/logging"
)

// InstallReq executes a shell command and returns stdout, stderr, exit code, and pid
func InstallReq(command string) (string, string, int, int) {
	// Equivalent to shlex.split(cmd)
	args := strings.Fields(command)
	if len(args) == 0 {
		return "", "No command provided", -1, 0
	}

	cmd := exec.Command(args[0], args[1:]...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	returnCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		} else {
			returnCode = -1
		}
	}

	pid := 0
	if cmd.Process != nil {
		pid = cmd.Process.Pid
	}

	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), returnCode, pid
}

// Git handles upstream repository syncing and dependency installation
func Git() {
	logging.InfoLogger.Println("Railway Deployment Detected: Skipping Git Upstream Check.")
	return // Early return exact match to the Python script

	// --- REACHABLE ONLY IF RETURN IS REMOVED ---

	repoLink := config.UPSTREAM_REPO
	var upstreamRepo string

	if config.GIT_TOKEN != "" {
		// Parsing the URL assuming format: https://github.com/username/repo
		parts := strings.Split(repoLink, "com/")
		if len(parts) > 1 {
			gitUsername := strings.Split(parts[1], "/")[0]
			tempRepo := strings.Split(repoLink, "https://")[1]
			upstreamRepo = fmt.Sprintf("https://%s:%s@%s", gitUsername, config.GIT_TOKEN, tempRepo)
		} else {
			upstreamRepo = config.UPSTREAM_REPO
		}
	} else {
		upstreamRepo = config.UPSTREAM_REPO
	}

	// In Go, using os/exec for git commands is often cleaner than a heavy git library
	// for simple deployment scripts.
	checkGit := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := checkGit.Run(); err == nil {
		logging.InfoLogger.Println("Git Client Found [VPS DEPLOYER]")
	} else {
		// Initialize repo if invalid
		exec.Command("git", "init").Run()

		// Setup remote
		remoteCheck := exec.Command("git", "remote", "get-url", "origin")
		if err := remoteCheck.Run(); err != nil {
			exec.Command("git", "remote", "add", "origin", upstreamRepo).Run()
		}

		exec.Command("git", "fetch", "origin").Run()
		exec.Command("git", "checkout", "-b", config.UPSTREAM_BRANCH, "origin/"+config.UPSTREAM_BRANCH).Run()

		// Attempt to set standard remote
		exec.Command("git", "remote", "set-url", "origin", config.UPSTREAM_REPO).Run()

		// Pull changes
		exec.Command("git", "fetch", "origin", config.UPSTREAM_BRANCH).Run()
		pullCmd := exec.Command("git", "pull", "origin", config.UPSTREAM_BRANCH)
		if err := pullCmd.Run(); err != nil {
			exec.Command("git", "reset", "--hard", "FETCH_HEAD").Run()
		}

		InstallReq("pip3 install --no-cache-dir -r requirements.txt")
		logging.InfoLogger.Println("Fetching updates from upstream repository...")
	}
}
