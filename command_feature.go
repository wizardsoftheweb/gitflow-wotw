package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	CommandFeature = cli.Command{
		Name:   "feature",
		Flags:  []cli.Flag{},
		Action: CommandFeatureAction,
		Subcommands: []cli.Commands{
			CommandFeatureList,
			CommandFeatureStart,
			CommandFeatureFinish,
			CommandFeaturePublish,
			CommandFeatureTrack,
			CommandFeatureDiff,
			CommandFeatureRebase,
			CommandFeatureCheckout,
			CommandFeaturePull,
		},
		Before: BeforeFeature,
	}
)

func BeforeFeature(context *cli.Context) error {
	Repo.Prefix = GitConfig.GetWithDefault(FEATURE_PREFIX_KEY, DefaultPrefixFeature.Value)
	Repo.HumanPrefix = strings.TrimSuffix(Repo.Prefix, "/")
}

func CommandFeatureAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureAction")
	return nil
}
func CommandFeatureListAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureListAction")

	return nil
}
func CommandFeatureStartAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureStartAction")
	return nil
}
func CommandFeatureFinishAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureFinishAction")
	return nil
}
func CommandFeaturePublishAction(context *cli.Context) error {
	logrus.Debug("CommandFeaturePublishAction")
	return nil
}
func CommandFeatureTrackAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureTrackAction")
	return nil
}
func CommandFeatureDiffAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureDiffAction")
	return nil
}
func CommandFeatureRebaseAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureRebaseAction")
	return nil
}
func CommandFeatureCheckoutAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureCheckoutAction")
	return nil
}
func CommandFeaturePullAction(context *cli.Context) error {
	logrus.Debug("CommandFeaturePullAction")
	return nil
}
