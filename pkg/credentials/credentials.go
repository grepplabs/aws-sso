package credentials

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"time"
)

const TokenMinDuration = 15 * time.Second

func RetrieveRoleCredentials(profile string) (*RoleCredentials, error) {
	ssoProfile, err := GetSSOProfile(profile)
	if err != nil {
		return nil, err
	}
	accessToken, err := GetAccessToken(ssoProfile.StartURL, ssoProfile.Region)
	if err != nil {
		return nil, err
	}
	if accessToken == "" {
		return nil, errors.Errorf("please login with 'aws sso login --profile=%s'\n", profile)
	}
	roleCredentials, err := GetRoleCredentials(ssoProfile, accessToken)
	if err != nil {
		return nil, err
	}
	if roleCredentials == nil {
		return nil, errors.Errorf("please login with 'aws sso login --profile=%s'\n", profile)
	}
	return roleCredentials, nil
}

func RefreshProfileCredentials(profileName string, roleCredentials *RoleCredentials) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	configPath := path.Join(home, ".aws", "credentials")

	p, err := configparser.NewConfigParserFromFile(configPath)
	if err != nil {
		return err
	}
	if !p.HasSection(profileName) {
		err = p.AddSection(profileName)
		if err != nil {
			return err
		}
	}
	if err = p.Set(profileName, "aws_access_key_id", roleCredentials.AccessKeyId); err != nil {
		return err
	}
	if err = p.Set(profileName, "aws_secret_access_key", roleCredentials.SecretAccessKey); err != nil {
		return err
	}
	if err = p.Set(profileName, "aws_session_token", roleCredentials.SessionToken); err != nil {
		return err
	}
	if err = p.SaveWithDelimiter(configPath, "="); err != nil {
		return err
	}
	return nil
}

type RoleCredentials struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	Expiration      int64  `json:"expiration"`
}

func GetRoleCredentials(ssoProfile *SSOProfile, accessToken string) (*RoleCredentials, error) {
	cmd := exec.Command(
		"aws", "sso", "get-role-credentials",
		"--profile", ssoProfile.ProfileName,
		"--role-name", ssoProfile.RoleName,
		"--account-id", ssoProfile.AccountID,
		"--access-token", accessToken,
		"--output", "json")

	cmd.Dir = ""
	outputBuffer := new(bytes.Buffer)
	outputWriter := bufio.NewWriter(outputBuffer)
	cmd.Stdout = outputWriter

	errBuffer := new(bytes.Buffer)
	errWriter := bufio.NewWriter(errBuffer)
	cmd.Stderr = errWriter

	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, errBuffer.String())
	}
	output := outputBuffer.String()

	result := struct {
		RoleCredentials RoleCredentials `json:"roleCredentials"`
	}{}
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return nil, err
	}
	return &result.RoleCredentials, nil
}

func getAccessTokenFromFile(filename string, ssoStartUrl string, ssoRegion string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	result := struct {
		StartUrl    string `json:"startUrl"`
		Region      string `json:"Region"`
		ExpiresAt   string `json:"expiresAt"`
		AccessToken string `json:"accessToken"`
	}{}

	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}
	if result.StartUrl != ssoStartUrl || result.Region != ssoRegion {
		return "", nil
	}
	expiresAt, err := time.Parse("2006-01-02T15:04:05Z", strings.Replace(result.ExpiresAt, "UTC", "Z", 1))
	if err != nil {
		return "", err
	}
	if expiresAt.Before(time.Now().Add(-TokenMinDuration)) {
		return "", err
	}
	return result.AccessToken, nil
}

func GetAccessToken(ssoStartURL string, ssoRegion string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	cacheDir := path.Join(home, ".aws", "sso", "cache")
	fileInfos, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return "", err
	}

	var accessToken string
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		cacheFile := path.Join(cacheDir, fileInfo.Name())
		accessToken, err = getAccessTokenFromFile(cacheFile, ssoStartURL, ssoRegion)
		if err == nil && len(accessToken) != 0 {
			return accessToken, nil
		}
	}
	return "", nil
}

type SSOProfile struct {
	ProfileName string

	StartURL  string
	Region    string
	AccountID string
	RoleName  string
}

func GetSSOProfile(profileName string) (*SSOProfile, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	configPath := path.Join(home, ".aws", "config")

	p, err := configparser.NewConfigParserFromFile(configPath)
	if err != nil {
		return nil, err
	}
	section := fmt.Sprintf("profile %s", profileName)
	dict, err := p.Items(section)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find profile '%s' in '%s'", profileName, configPath)
	}
	ret := &SSOProfile{ProfileName: profileName}
	if ret.StartURL, err = getAttribute(dict, "sso_start_url"); err != nil {
		return nil, err
	}
	if ret.Region, err = getAttribute(dict, "sso_region"); err != nil {
		return nil, err
	}
	if ret.AccountID, err = getAttribute(dict, "sso_account_id"); err != nil {
		return nil, err
	}
	if ret.RoleName, err = getAttribute(dict, "sso_role_name"); err != nil {
		return nil, err
	}
	return ret, nil
}

func getAttribute(dict map[string]string, tag string) (string, error) {
	value := dict[tag]
	if value == "" {
		return "", errors.Errorf("'%s' not in '%s' profile", tag, dict)
	}
	return value, nil
}
