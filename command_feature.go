package main

import (
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
}

func CommandFeatureAction(context *cli.Context) error {
	logrus.Debug("CommandFeatureAction")
	return nil
}
func CommandFeatureList(context *cli.Context) error {
	logrus.Debug("CommandFeatureList")
	return nil
}
func CommandFeatureStart(context *cli.Context) error {
	logrus.Debug("CommandFeatureStart")
	return nil
}
func CommandFeatureFinish(context *cli.Context) error {
	logrus.Debug("CommandFeatureFinish")
	return nil
}
func CommandFeaturePublish(context *cli.Context) error {
	logrus.Debug("CommandFeaturePublish")
	return nil
}
func CommandFeatureTrack(context *cli.Context) error {
	logrus.Debug("CommandFeatureTrack")
	return nil
}
func CommandFeatureDiff(context *cli.Context) error {
	logrus.Debug("CommandFeatureDiff")
	return nil
}
func CommandFeatureRebase(context *cli.Context) error {
	logrus.Debug("CommandFeatureRebase")
	return nil
}
func CommandFeatureCheckout(context *cli.Context) error {
	logrus.Debug("CommandFeatureCheckout")
	return nil
}
func CommandFeaturePull(context *cli.Context) error {
	logrus.Debug("CommandFeaturePull")
	return nil
}
