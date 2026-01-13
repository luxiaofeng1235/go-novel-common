/*
 * @Descripttion: 本地文件服务（上传等）
 * @Author: red
 * @Date: 2026-01-12 12:20:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 12:20:00
 */
package file_service

import (
	"fmt"
	"go-novel/utils"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UploadResult struct {
	LocalPath   string `json:"localPath"`
	PublicPath  string `json:"publicPath"`
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
}

// LocalUpload 本地上传：把文件保存到 upload.baseDir 下，并返回可被 source 静态服务访问的 URL
// - subDir：相对目录（可为空），会被安全清洗，禁止 .. 等路径穿越
func LocalUpload(c *gin.Context, subDir string) (*UploadResult, error) {
	maxSizeMB := viper.GetInt("upload.maxSizeMB")
	if maxSizeMB <= 0 {
		maxSizeMB = 50
	}
	maxBytes := int64(maxSizeMB) * 1024 * 1024
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败")
	}
	if maxBytes > 0 && fileHeader.Size > maxBytes {
		return nil, fmt.Errorf("文件过大")
	}
	if fileHeader == nil {
		return nil, fmt.Errorf("上传文件不能为空")
	}
	if err := validateUploadFile(fileHeader); err != nil {
		return nil, err
	}
	baseDir := viper.GetString("upload.baseDir")
	if strings.TrimSpace(baseDir) == "" {
		baseDir = "./public/upload"
	}

	publicPrefix := viper.GetString("upload.publicPathPrefix")
	if strings.TrimSpace(publicPrefix) == "" {
		publicPrefix = "/public/upload"
	}

	cleanSubDir, err := sanitizeSubDir(subDir)
	if err != nil {
		return nil, err
	}

	dstDir := baseDir
	publicPath := publicPrefix
	if cleanSubDir != "" {
		dstDir = filepath.Join(baseDir, filepath.FromSlash(cleanSubDir))
		publicPath = path.Join(publicPrefix, cleanSubDir)
	}

	dstFilename := utils.NewRandomFilename(fileHeader.Filename)
	localPath, err := utils.SaveMultipartFile(fileHeader, dstDir, dstFilename)
	if err != nil {
		return nil, err
	}

	publicPath = path.Join(publicPath, dstFilename)
	fullURL := buildPublicURL(publicPath)

	return &UploadResult{
		LocalPath:  filepath.ToSlash(localPath),
		PublicPath: publicPath,
		URL:        fullURL,
		Filename:   dstFilename,
		Size:       fileHeader.Size,
		ContentType: func() string {
			if fileHeader.Header == nil {
				return ""
			}
			return fileHeader.Header.Get("Content-Type")
		}(),
	}, nil
}

func sanitizeSubDir(subDir string) (string, error) {
	subDir = strings.TrimSpace(subDir)
	if subDir == "" {
		return "", nil
	}
	subDir = strings.ReplaceAll(subDir, "\\", "/")

	// 只允许常见安全字符，避免奇怪路径/编码问题
	re := regexp.MustCompile(`^[a-zA-Z0-9/_-]{1,100}$`)
	if !re.MatchString(subDir) {
		return "", fmt.Errorf("dir 不合法")
	}
	if strings.Contains(subDir, "..") {
		return "", fmt.Errorf("dir 不合法")
	}

	// 进一步 clean，确保不会生成绝对路径
	clean := path.Clean("/" + subDir)
	clean = strings.TrimPrefix(clean, "/")
	if clean == "." || clean == "" {
		return "", nil
	}
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("dir 不合法")
	}
	return clean, nil
}

func buildPublicURL(publicPath string) string {
	if strings.TrimSpace(publicPath) == "" {
		return ""
	}
	base := strings.TrimRight(utils.GetSourceBaseUrl(), "/")
	if base == "" {
		// 兜底：至少保证能拼出一个 URL，避免上传接口因配置缺失直接不可用
		base = strings.TrimRight(utils.GetApiUrl(), "/")
	}
	if base == "" {
		return publicPath
	}
	return base + publicPath
}

func validateUploadFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return fmt.Errorf("上传文件不能为空")
	}

	// 后缀白名单（推荐配置）
	allowedExts := viper.GetStringSlice("upload.allowedExts")
	if len(allowedExts) > 0 {
		ext := strings.ToLower(path.Ext(fileHeader.Filename))
		ok := false
		for _, allow := range allowedExts {
			allow = strings.ToLower(strings.TrimSpace(allow))
			if allow == "" {
				continue
			}
			if !strings.HasPrefix(allow, ".") {
				allow = "." + allow
			}
			if ext == allow {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("不支持的文件类型")
		}
	}

	// Content-Type 前缀白名单（可选）
	allowedMimePrefixes := viper.GetStringSlice("upload.allowedMimePrefixes")
	if len(allowedMimePrefixes) > 0 && fileHeader.Header != nil {
		contentType := strings.ToLower(strings.TrimSpace(fileHeader.Header.Get("Content-Type")))
		if contentType != "" {
			ok := false
			for _, p := range allowedMimePrefixes {
				p = strings.ToLower(strings.TrimSpace(p))
				if p != "" && strings.HasPrefix(contentType, p) {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("不支持的文件类型")
			}
		}
	}

	return nil
}
