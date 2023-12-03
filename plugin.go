package main

import (
	"context"
	"errors"
	"os"
	"strconv"

	"codeberg.org/woodpecker-plugins/go-plugin"
	"github.com/joho/godotenv"
	"github.com/Masterminds/semver/v3"
)

type Plugin struct {
	*plugin.Plugin
}

// type alias which can hold environment vars and its data.
type EnvDataMap map[string]string

func (p *Plugin) execute(ctx context.Context) error {
	tag := p.Metadata.Curr.Tag
	if tag == "" {
		return nil
	}

	env, err := godotenv.Read()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	v, err := semver.NewVersion(tag)
	if err != nil {
		env["CI_COMMIT_TAG_IS_SEMVER"] = "false"
	} else {
		env["CI_COMMIT_TAG_IS_SEMVER"] = "true"
		env["CI_COMMIT_TAG_SEMVER"] = v.String()
		env["CI_COMMIT_TAG_SEMVER_MAJOR"] = strconv.FormatUint(v.Major(), 10)
		env["CI_COMMIT_TAG_SEMVER_MINOR"] = strconv.FormatUint(v.Minor(), 10)
		env["CI_COMMIT_TAG_SEMVER_PATCH"] = strconv.FormatUint(v.Patch(), 10)
		if v.Prerelease() != "" {
			env["CI_COMMIT_TAG_SEMVER_PRERELEASE"] = v.Prerelease()
			env["CI_COMMIT_TAG_IS_PRERELEASE"] = "true"
		}
	}

	return godotenv.Write(env, ".env")
}
