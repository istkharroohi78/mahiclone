package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"ANJALI/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

var gitConfigCache = make(map[string]string)
var gitConfigTimestamp time.Time

func RegisterRepoUploadHandlers(b *tb.Bot) {
	cfg := config.LoadConfig()

	// /gitconfig username email token
	b.Handle("/gitconfig", func(m *tb.Message) {
		if int64(m.Sender.ID) != cfg.OwnerID {
			return
		}

		args := strings.Split(m.Text, " ")
		if len(args) < 4 {
			b.Send(m.Chat, "Usage: `/gitconfig [username] [email] [github_token]`", tb.ModeMarkdown)
			return
		}

		gitConfigCache["name"] = args[1]
		gitConfigCache["email"] = args[2]
		gitConfigCache["token"] = args[3]
		gitConfigTimestamp = time.Now()

		b.Delete(m) // Security: Delete message with token
		b.Send(m.Chat, "✅ **GitHub credentials configured for next 5 minutes.**", tb.ModeMarkdown)
	})

	// /gitupload repo_name private/public branch_name
	b.Handle("/gitupload", func(m *tb.Message) {
		if int64(m.Sender.ID) != cfg.OwnerID {
			return
		}

		if time.Since(gitConfigTimestamp) > 5*time.Minute || gitConfigCache["token"] == "" {
			b.Send(m.Chat, "⚠️ **Git config expired or not set. Run `/gitconfig` first.**", tb.ModeMarkdown)
			return
		}

		if m.ReplyTo == nil || m.ReplyTo.Document == nil || !strings.HasSuffix(m.ReplyTo.Document.FileName, ".zip") {
			b.Send(m.Chat, "⚠️ **Please reply to a `.zip` file of your project repo.**", tb.ModeMarkdown)
			return
		}

		args := strings.Split(m.Text, " ")
		repoName := "uploaded-repo"
		if len(args) > 1 {
			repoName = args[1]
		}

		statusMsg, _ := b.Send(m.Chat, fmt.Sprintf("⏳ **Processing & uploading repository `%s` to GitHub...**", repoName), tb.ModeMarkdown)

		// Download Document
		zipPath := fmt.Sprintf("/tmp/%s.zip", repoName)
		err := b.Download(&m.ReplyTo.Document.File, zipPath)
		if err != nil {
			b.Edit(statusMsg, fmt.Sprintf("❌ **Download Error:** %v", err))
			return
		}
		defer os.Remove(zipPath)

		// Shell execution to unzip and push
		extractDir := fmt.Sprintf("/tmp/%s_extract", repoName)
		os.MkdirAll(extractDir, 0755)
		defer os.RemoveAll(extractDir)

		unzipCmd := exec.Command("unzip", "-q", zipPath, "-d", extractDir)
		unzipCmd.Run()

		gitCmds := fmt.Sprintf(`
			cd %s
			git init
			git config user.name "%s"
			git config user.email "%s"
			git add .
			git commit -m "Uploaded via THE SHIV Music Engine"
			git branch -M main
			git remote add origin https://%s@github.com/%s/%s.git
			git push -u origin main --force
		`, extractDir, gitConfigCache["name"], gitConfigCache["email"], gitConfigCache["token"], gitConfigCache["name"], repoName)

		cmd := exec.Command("bash", "-c", gitCmds)
		var out, stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			b.Edit(statusMsg, fmt.Sprintf("❌ **Git Push Error:**\n`%s`", stderr.String()), tb.ModeMarkdown)
			return
		}

		repoURL := fmt.Sprintf("https://github.com/%s/%s", gitConfigCache["name"], repoName)
		b.Edit(statusMsg, fmt.Sprintf("✅ **Repository `%s` successfully pushed to GitHub!**\n\n🔗 URL: %s", repoName, repoURL), tb.ModeMarkdown)
	})
}
