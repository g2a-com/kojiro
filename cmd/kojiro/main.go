package main

import (
	"github.com/g2a-com/gojiro/pkg/git"
	"github.com/g2a-com/gojiro/pkg/jira"
	"github.com/g2a-com/gojiro/pkg/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	defaultRetries     = 3
	jiraVersionFlag    = "jira-version"
	tagFlag            = " tag"
	jiraEmailFlag      = "jira-email"
	jiraTokenFlag      = "jira-token"
	jiraProjectFlag    = "jira-project"
	jiraBaseUrlFlag    = "jira-base-url"
	jiraRetryTimesFlag = "jira-retry-times"
	dirFlag            = "dir"
	klioLoggerFlag     = "klio-logger"
	dryRunFlag         = "dry-run"
)

type cojiroContext struct {
	JiraVersion string
	JiraEmail   string
	JiraToken   string
	JiraProject string
	JiraBaseUrl string
	JiraRetries int
	Tag         string
	Dir         string
	DryRun      bool
	KlioLogger  bool
}

var ctx cojiroContext
var logger log.Logger

func main() {
	rootCmd := &cobra.Command{
		Use:   "kojiro",
		Short: "A simple release setter for Jira tasks since last tag",
		Long:  `Automatically create Jira release nad link all issues based commits. All automatically.`,
		RunE:  rootFunc,
	}
	// get current directory path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	pwd := filepath.Dir(ex)

	ctx = cojiroContext{}

	rootCmd.Flags().StringVarP(&ctx.JiraVersion, jiraVersionFlag, "v", "", "Version name for Jira")
	rootCmd.Flags().StringVarP(&ctx.Tag, tagFlag, "t", "", "Existing git tag")
	rootCmd.Flags().StringVarP(&ctx.JiraEmail, jiraEmailFlag, "e", "", "Jira email")
	rootCmd.Flags().StringVarP(&ctx.JiraToken, jiraTokenFlag, "k", "", "Jira token/key/password")
	rootCmd.Flags().StringVarP(&ctx.JiraProject, jiraProjectFlag, "p", "", "Jira project, it has to be ID, example: 10003")
	rootCmd.Flags().StringVarP(&ctx.JiraBaseUrl, jiraBaseUrlFlag, "u", "", "Jira service base url, example: https://example.atlassian.net")
	rootCmd.Flags().IntVarP(&ctx.JiraRetries, jiraRetryTimesFlag, "r", defaultRetries, "Jira retry times for HTTP requests if failed")
	rootCmd.Flags().StringVarP(&ctx.Dir, dirFlag, "d", pwd, "Absolute directory path to git repository")
	rootCmd.Flags().BoolVarP(&ctx.KlioLogger, klioLoggerFlag, "l", true, "Use Klio-compliant logger, hard to read without using Klio")
	rootCmd.Flags().BoolVarP(&ctx.DryRun, dryRunFlag, "", false, "Enable dry run mode")

	_ = rootCmd.MarkFlagRequired(tagFlag)
	_ = rootCmd.MarkFlagRequired(jiraEmailFlag)
	_ = rootCmd.MarkFlagRequired(jiraTokenFlag)
	_ = rootCmd.MarkFlagRequired(jiraProjectFlag)
	_ = rootCmd.MarkFlagRequired(jiraBaseUrlFlag)

	rootCmd.Example = "kojiro -e jira@example.com -k pa$$wor0 -p 10003 -t v1.1.0 -u https://example.atlassian.net"

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootFunc(c *cobra.Command, _ []string) error {
	if ctx.KlioLogger {

	} else {
		zapLogger := zap.NewExample().Sugar()

		defer func() { _ = zapLogger.Sync() }()
		logger = zapLogger
	}

	if ctx.JiraVersion == "" {
		ctx.JiraVersion = ctx.Tag
	}

	logger.Debugf(
		"starting with parameters: %+v",
		map[string]interface{}{
			"jiraEmail":      ctx.JiraEmail,
			"jiraToken":      ctx.JiraToken,
			"jiraProject":    ctx.JiraProject,
			"jiraBaseURL":    ctx.JiraBaseUrl,
			"jiraRetryTimes": ctx.JiraRetries,
			"gitDir":         ctx.Dir,
			"tag":            ctx.Tag,
			"version":        ctx.JiraVersion,
			"dryRun":         ctx.DryRun,
		},
	)
	logger.Infof("[JIRA-VERSIONER] git directory: %s", ctx.Dir)

	g := git.New(ctx.Dir, logger)

	tasks, err := g.GetTasks(ctx.Tag)
	if err != nil {
		logger.Errorf("[GIT] error while getting tasks since latest commit %+v", err)
		return err
	}

	jiraConfig := jira.Config{
		Username:       ctx.JiraEmail,
		Token:          ctx.JiraToken,
		ProjectID:      ctx.JiraProject,
		BaseURL:        ctx.JiraBaseUrl,
		Log:            logger,
		DryRun:         ctx.DryRun,
		HTTPMaxRetries: ctx.JiraRetries,
	}
	j, err := jira.New(&jiraConfig)
	if err != nil {
		logger.Errorf("[VERSION] error while connecting to jira server %+v", err)
		return err
	}

	_, err = j.CreateVersion(ctx.JiraVersion)
	if err != nil {
		logger.Errorf("[VERSION] error while creating version %+v", err)
		return err
	}

	j.LinkTasksToVersion(tasks)

	logger.Infof("[JIRA-VERSIONER] done")
	return nil
}
