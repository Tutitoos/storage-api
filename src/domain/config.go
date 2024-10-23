package domain

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type IConfig struct {
	Port                      int
	ApiUrl                    string
	Token                     string
	CloudflareAccountId       string
	CloudflareAccessKeyId     string
	CloudflareSecretAccessKey string
	BucketName                string
	BucketRegion              string
	BucketUrl                 string
	ExcludeFolders            []string
	ExcludeFiles              []string
	WhitelistIps              []string
	BypassWhitelist           string
}

func Config() *IConfig {
	isDocker := runningInDocker()
	if !isDocker {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Invalid PORT value")
	}

	tryPort, err := strconv.Atoi(os.Getenv("TRY_PORT"))
	if err != nil {
		log.Fatalf("Invalid TRY_PORT value")
	}

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		log.Fatalf("Invalid API_URL value")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalf("Invalid TOKEN value")
	}

	cloudflareAccountId := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	if cloudflareAccountId == "" {
		log.Fatalf("Invalid CLOUDFLARE_ACCOUNT_ID value")
	}

	cloudflareAccessKeyId := os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")
	if cloudflareAccessKeyId == "" {
		log.Fatalf("Invalid CLOUDFLARE_ACCESS_KEY_ID value")
	}

	cloudflareSecretAccessKey := os.Getenv("CLOUDFLARE_ACCESS_KEY_SECRET")
	if cloudflareSecretAccessKey == "" {
		log.Fatalf("Invalid CLOUDFLARE_ACCESS_KEY_SECRET value")
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("Invalid BUCKET_NAME value")
	}

	bucketRegion := os.Getenv("BUCKET_REGION")
	if bucketRegion == "" {
		log.Fatalf("Invalid BUCKET_REGION value")
	}

	whitelistIps := os.Getenv("WHITELIST_IPS")
	if whitelistIps == "" {
		whitelistIps = "127.0.0.1,::1"
	}

	if port == 0 {
		port = tryPort
	}

	return &IConfig{
		Port:                      port,
		ApiUrl:                    strings.Replace(apiUrl, "{PORT}", strconv.Itoa(port), -1),
		Token:                     token,
		CloudflareAccountId:       cloudflareAccountId,
		CloudflareAccessKeyId:     cloudflareAccessKeyId,
		CloudflareSecretAccessKey: cloudflareSecretAccessKey,
		BucketName:                bucketName,
		BucketRegion:              bucketRegion,
		BucketUrl:                 "https://" + cloudflareAccountId + ".r2.cloudflarestorage.com",
		ExcludeFolders:            strings.Split(os.Getenv("EXCLUDE_FOLDER"), ","),
		ExcludeFiles:              strings.Split(os.Getenv("EXCLUDE_FILE"), ","),
		WhitelistIps:              strings.Split(whitelistIps, ","),
		BypassWhitelist:           os.Getenv("BYPASS_WHITELIST"),
	}
}

func runningInDocker() bool {
	_, err := os.Stat("/proc/1/cgroup")
	return err == nil
}

var CONFIG = Config()
